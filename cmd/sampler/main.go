// cmd/sampler/main.go
// Fully commented implementation of a Sample Composer and Unison synthesizer GUI in Fyne,
package main

import (
	"bufio"
	"fmt"
	"image/color"
	"math"
	"math/cmplx"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/terefang/gommons/pkg/xaudio"
)

// hardcodedPresetNames lists the built-in presets that are always available.
var hardcodedPresetNames = []string{
	"Analog Trance Synth",
	"Hypersaw",
	"MultiSaw",
}

// isHardcodedPreset returns true if name matches one of the built-in presets.
func isHardcodedPreset(name string) bool {
	trimmed := strings.TrimSpace(name)
	for _, h := range hardcodedPresetNames {
		if h == trimmed {
			return true
		}
	}
	return false
}

// customSynthTheme implements fyne.Theme to allow instant runtime switching
// between dark and light appearance modes while delegating base metrics/colors to fyne.DefaultTheme().
type customSynthTheme struct {
	dark bool
}

func (t *customSynthTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	v := theme.VariantLight
	if t.dark {
		v = theme.VariantDark
	}
	return theme.DefaultTheme().Color(name, v)
}

func (t *customSynthTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customSynthTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customSynthTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// autoFitAccordion wraps widget.Accordion so opening or closing an item
// immediately recalculates and fits the window size without scrollbars.
type autoFitAccordion struct {
	widget.Accordion
	onToggle func()
}

func newAutoFitAccordion(onToggle func(), title string, content fyne.CanvasObject) *autoFitAccordion {
	acc := &autoFitAccordion{onToggle: onToggle}
	item := widget.NewAccordionItem(title, content)
	acc.Items = []*widget.AccordionItem{item}
	acc.ExtendBaseWidget(acc)
	return acc
}

func (a *autoFitAccordion) Open(index int) {
	a.Accordion.Open(index)
	if a.onToggle != nil {
		fyne.Do(a.onToggle)
	}
}

func (a *autoFitAccordion) Close(index int) {
	a.Accordion.Close(index)
	if a.onToggle != nil {
		fyne.Do(a.onToggle)
	}
}

// stretchVizLayout resizes the visualization background and triggers redraws whenever
// the visualization container changes size.
type stretchVizLayout struct {
	ui *synthUI
}

func (l *stretchVizLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(300, 160)
}

func (l *stretchVizLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if size.Width <= 1 || size.Height <= 1 {
		return
	}
	prevW := l.ui.vizLastSize.Width
	prevH := l.ui.vizLastSize.Height
	l.ui.vizBg.Move(fyne.NewPos(0, 0))
	l.ui.vizBg.Resize(size)

	// Only re-render waveforms/bars when container size meaningfully changes to avoid layout loops.
	if math.Abs(float64(size.Width-prevW)) > 1.0 || math.Abs(float64(size.Height-prevH)) > 1.0 {
		l.ui.vizLastSize = size
		fyne.Do(func() {
			l.ui.redrawVisualization()
		})
	}
}

// SequenceStep defines a single editable step on the timeline that maps to xaudio.NoteEvent.
type SequenceStep struct {
	StartSec    float64 // Step start time in seconds along the timeline
	Freq        float64 // Pitch frequency in Hz
	DurationSec float64 // Held duration in seconds
	Velocity    float64 // Note velocity multiplier (0.0 to 1.0)
}

// SynthParams holds every parametric value as an explicit configuration field.
type SynthParams struct {
	SampleRate   int                  // Sample rate in Hz (8000, 16000, 32000, 48000)
	LengthMs     float64              // Total single-note duration in milliseconds (1000 to 30000)
	Tune         float64              // Base oscillator frequency in Hz
	Tone         float64              // Tone filter / modulation index balance (0.0 to 1.0)
	Snappy       float64              // Noise transient mix and attack scaling (0.0 to 1.0)
	BaseWave     xaudio.Waveform      // Oscillator waveform shape (Sine, Sawtooth, SawtoothPoly, Square, SquarePoly, Triangle, FM)
	Voices       int                  // Unison voice count (>= 1)
	Detune       float64              // Pitch spread across voices (0.0 to 1.0)
	DetuneCurve  xaudio.DetuneCurve   // Voice detune distribution curve (CurveLinear, CurveJP8000, CurveExponential)
	Mix          float64              // Balance between side voices and center voice (0.0 to 1.0)
	StereoWidth  float64              // Stereo pan spread (0.0 to 1.0)
	SubOctave    bool                 // Enable -1 octave sub-bass reinforcement
	AnalogDrift  float64              // LFO pitch wander per voice (0.0 to 1.0)
	HardSync     bool                 // Lock oscillator phases to master cycle
	PulseWidth   float64              // Square wave duty cycle (0.05 to 0.95)
	ModRatio     float64              // FM modulator-to-carrier frequency ratio
	ModIndex     float64              // FM modulation depth
	Attack       float64              // Attack time in seconds
	Decay        float64              // Decay time in seconds
	Sustain      float64              // Sustain amplitude level (0.0 to 1.0)
	Release      float64              // Release tail in seconds
	AttackCurve  xaudio.EnvelopeCurve // Attack stage interpolation curve
	DecayCurve   xaudio.EnvelopeCurve // Decay stage interpolation curve
	ReleaseCurve xaudio.EnvelopeCurve // Release stage interpolation curve
	WobbleCycles float64              // Modulation cycles for NewBicubicWobbleCurve
	WobbleDepth  float64              // Modulation depth for NewBicubicWobbleCurve
}

// defaultSynthParams returns the initial default configuration for all synthesizer parameters.
func defaultSynthParams() SynthParams {
	return SynthParams{
		SampleRate:   48000,
		LengthMs:     1500.0,
		Tune:         65.41,
		Tone:         0.65,
		Snappy:       0.50,
		BaseWave:     xaudio.SawtoothPoly,
		Voices:       7,
		Detune:       0.35,
		DetuneCurve:  xaudio.CurveJP8000,
		Mix:          0.75,
		StereoWidth:  0.80,
		SubOctave:    true,
		AnalogDrift:  0.15,
		HardSync:     false,
		PulseWidth:   0.50,
		ModRatio:     2.0,
		ModIndex:     1.0,
		Attack:       0.01,
		Decay:        0.35,
		Sustain:      0.60,
		Release:      0.25,
		AttackCurve:  xaudio.EnvHermite,
		DecayCurve:   xaudio.EnvExponential,
		ReleaseCurve: xaudio.EnvLogarithmic,
		WobbleCycles: 3.0,
		WobbleDepth:  0.15,
	}
}

// defaultSequenceSteps returns a fresh copy of the initial default timeline steps.
func defaultSequenceSteps() []SequenceStep {
	return []SequenceStep{
		{StartSec: 0.0, Freq: 65.41, DurationSec: 0.45, Velocity: 0.90},
		{StartSec: 0.5, Freq: 77.78, DurationSec: 0.45, Velocity: 0.85},
		{StartSec: 1.0, Freq: 98.00, DurationSec: 0.45, Velocity: 0.85},
		{StartSec: 1.5, Freq: 130.81, DurationSec: 0.60, Velocity: 0.95},
	}
}

// resolvePresetsIniPath returns <XDG user config dir>/gommons-xaudio/presets.ini.
func resolvePresetsIniPath() string {
	cfgDir, err := os.UserConfigDir()
	if err != nil || strings.TrimSpace(cfgDir) == "" {
		return "presets.ini"
	}
	return filepath.Join(cfgDir, "SamplesComposer", "presets.ini")
}

// formatWaveformName converts an xaudio.Waveform value to its UI/INI label.
func formatWaveformName(w xaudio.Waveform) string {
	switch w {
	case xaudio.Sine:
		return "Sine"
	case xaudio.Sawtooth:
		return "Sawtooth"
	case xaudio.Square:
		return "Square"
	case xaudio.SquarePoly:
		return "SquarePoly"
	case xaudio.Triangle:
		return "Triangle"
	case xaudio.FM:
		return "FM"
	default:
		return "SawtoothPoly"
	}
}

// formatDetuneCurveName converts an xaudio.DetuneCurve value to its UI/INI label.
func formatDetuneCurveName(dc *xaudio.DetuneCurve) string {
	switch dc {
	case &xaudio.CurveLinear:
		return "Linear"
	case &xaudio.CurveExponential:
		return "Exponential"
	case &xaudio.CurveJP8000:
		return "Analog Trance Synth"
	default:
		return "Custom"
	}
}

// formatEnvelopeCurveName converts an xaudio.EnvelopeCurve value to a curve name when possible.
func formatEnvelopeCurveName(ec xaudio.EnvelopeCurve, fallback string) string {
	switch {
	case fmt.Sprintf("%p", ec) == fmt.Sprintf("%p", xaudio.EnvLinear):
		return "Linear"
	case fmt.Sprintf("%p", ec) == fmt.Sprintf("%p", xaudio.EnvExponential):
		return "Exponential"
	case fmt.Sprintf("%p", ec) == fmt.Sprintf("%p", xaudio.EnvLogarithmic):
		return "Logarithmic"
	case fmt.Sprintf("%p", ec) == fmt.Sprintf("%p", xaudio.EnvHermite):
		return "Hermite"
	case fmt.Sprintf("%p", ec) == fmt.Sprintf("%p", xaudio.EnvRadian):
		return "Radian"
	default:
		return fallback
	}
}

// loadPresetsFromINI reads INI sections from presets.ini in the XDG user config directory.
func loadPresetsFromINI(path string) (map[string]SynthParams, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	presets := make(map[string]SynthParams)
	var currentSection string
	currentParams := defaultSynthParams()

	flushSection := func() {
		if strings.TrimSpace(currentSection) != "" {
			presets[currentSection] = currentParams
		}
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			flushSection()
			currentSection = strings.TrimSpace(line[1 : len(line)-1])
			currentParams = defaultSynthParams()
			continue
		}
		if currentSection == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch key {
		case "sample_rate":
			if v, err := strconv.Atoi(val); err == nil {
				currentParams.SampleRate = v
			}
		case "length_ms":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.LengthMs = v
			}
		case "tune":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Tune = v
			}
		case "tone":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Tone = v
			}
		case "snappy":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Snappy = v
			}
		case "base_wave":
			currentParams.BaseWave = resolveWaveform(val)
		case "voices":
			if v, err := strconv.Atoi(val); err == nil && v >= 1 {
				currentParams.Voices = v
			}
		case "detune":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Detune = v
			}
		case "detune_curve":
			currentParams.DetuneCurve = resolveDetuneCurve(val)
		case "mix":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Mix = v
			}
		case "stereo_width":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.StereoWidth = v
			}
		case "sub_octave":
			if v, err := strconv.ParseBool(val); err == nil {
				currentParams.SubOctave = v
			}
		case "analog_drift":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.AnalogDrift = v
			}
		case "hard_sync":
			if v, err := strconv.ParseBool(val); err == nil {
				currentParams.HardSync = v
			}
		case "pulse_width":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.PulseWidth = v
			}
		case "mod_ratio":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.ModRatio = v
			}
		case "mod_index":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.ModIndex = v
			}
		case "attack":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Attack = v
			}
		case "decay":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Decay = v
			}
		case "sustain":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Sustain = v
			}
		case "release":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.Release = v
			}
		case "wobble_cycles":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.WobbleCycles = v
			}
		case "wobble_depth":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				currentParams.WobbleDepth = v
			}
		case "attack_curve":
			currentParams.AttackCurve = resolveEnvelopeCurve(val, currentParams.WobbleCycles, currentParams.WobbleDepth)
		case "decay_curve":
			currentParams.DecayCurve = resolveEnvelopeCurve(val, currentParams.WobbleCycles, currentParams.WobbleDepth)
		case "release_curve":
			currentParams.ReleaseCurve = resolveEnvelopeCurve(val, currentParams.WobbleCycles, currentParams.WobbleDepth)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	flushSection()
	return presets, nil
}

// writePresetsINI serializes the provided map of INI presets to presets.ini.
func writePresetsINI(path string, presets map[string]SynthParams) error {
	if len(presets) == 0 {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	names := make([]string, 0, len(presets))
	for k := range presets {
		names = append(names, k)
	}
	sort.Strings(names)

	var sb strings.Builder
	sb.WriteString("# SamplesComposer presets\n")
	for _, k := range names {
		pr := presets[k]
		sb.WriteString(fmt.Sprintf("\n[%s]\n", k))
		sb.WriteString(fmt.Sprintf("sample_rate=%d\n", pr.SampleRate))
		sb.WriteString(fmt.Sprintf("length_ms=%.2f\n", pr.LengthMs))
		sb.WriteString(fmt.Sprintf("tune=%.2f\n", pr.Tune))
		sb.WriteString(fmt.Sprintf("tone=%.4f\n", pr.Tone))
		sb.WriteString(fmt.Sprintf("snappy=%.4f\n", pr.Snappy))
		sb.WriteString(fmt.Sprintf("base_wave=%s\n", formatWaveformName(pr.BaseWave)))
		sb.WriteString(fmt.Sprintf("voices=%d\n", pr.Voices))
		sb.WriteString(fmt.Sprintf("detune=%.4f\n", pr.Detune))
		sb.WriteString(fmt.Sprintf("detune_curve=%s\n", formatDetuneCurveName(&pr.DetuneCurve)))
		sb.WriteString(fmt.Sprintf("mix=%.4f\n", pr.Mix))
		sb.WriteString(fmt.Sprintf("stereo_width=%.4f\n", pr.StereoWidth))
		sb.WriteString(fmt.Sprintf("sub_octave=%t\n", pr.SubOctave))
		sb.WriteString(fmt.Sprintf("analog_drift=%.4f\n", pr.AnalogDrift))
		sb.WriteString(fmt.Sprintf("hard_sync=%t\n", pr.HardSync))
		sb.WriteString(fmt.Sprintf("pulse_width=%.4f\n", pr.PulseWidth))
		sb.WriteString(fmt.Sprintf("mod_ratio=%.4f\n", pr.ModRatio))
		sb.WriteString(fmt.Sprintf("mod_index=%.4f\n", pr.ModIndex))
		sb.WriteString(fmt.Sprintf("attack=%.4f\n", pr.Attack))
		sb.WriteString(fmt.Sprintf("decay=%.4f\n", pr.Decay))
		sb.WriteString(fmt.Sprintf("sustain=%.4f\n", pr.Sustain))
		sb.WriteString(fmt.Sprintf("release=%.4f\n", pr.Release))
		sb.WriteString(fmt.Sprintf("attack_curve=%s\n", formatEnvelopeCurveName(pr.AttackCurve, "Hermite")))
		sb.WriteString(fmt.Sprintf("decay_curve=%s\n", formatEnvelopeCurveName(pr.DecayCurve, "Exponential")))
		sb.WriteString(fmt.Sprintf("release_curve=%s\n", formatEnvelopeCurveName(pr.ReleaseCurve, "Logarithmic")))
		sb.WriteString(fmt.Sprintf("wobble_cycles=%.2f\n", pr.WobbleCycles))
		sb.WriteString(fmt.Sprintf("wobble_depth=%.4f\n", pr.WobbleDepth))
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

// savePresetToINI writes or updates a named preset section inside presets.ini in the XDG user config dir.
func savePresetToINI(path, name string, p SynthParams) error {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return fmt.Errorf("preset name cannot be empty")
	}

	existing, err := loadPresetsFromINI(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if existing == nil {
		existing = make(map[string]SynthParams)
	}
	existing[trimmedName] = p
	return writePresetsINI(path, existing)
}

// removePresetFromINI removes a named INI preset section from presets.ini in the XDG user config dir.
func removePresetFromINI(path, name string) (bool, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return false, fmt.Errorf("preset name cannot be empty")
	}

	existing, err := loadPresetsFromINI(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if _, exists := existing[trimmedName]; !exists {
		return false, nil
	}
	delete(existing, trimmedName)
	if err := writePresetsINI(path, existing); err != nil {
		return false, err
	}
	return true, nil
}

// BuildUnisonConfig constructs an xaudio.UnisonConfig from current synth parameters.
func BuildUnisonConfig(p SynthParams) xaudio.UnisonConfig {
	effectiveModIndex := p.ModIndex + (p.Tone * 2.5)
	effectiveMix := math.Min(1.0, math.Max(0.0, p.Mix*(0.5+0.5*p.Tone)))
	return xaudio.UnisonConfig{
		BaseWave:    p.BaseWave,
		Voices:      p.Voices,
		Detune:      p.Detune,
		DetuneCurve: p.DetuneCurve,
		Mix:         effectiveMix,
		StereoWidth: p.StereoWidth,
		SubOctave:   p.SubOctave,
		AnalogDrift: p.AnalogDrift,
		HardSync:    p.HardSync,
		PulseWidth:  p.PulseWidth,
		ModRatio:    p.ModRatio,
		ModIndex:    effectiveModIndex,
	}
}

// BuildADSR constructs an xaudio.ADSR envelope from current synth parameters.
func BuildADSR(p SynthParams) xaudio.ADSR {
	effectiveAttack := math.Max(0.001, p.Attack*(1.0-0.7*p.Snappy))
	return xaudio.ADSR{
		Attack:       effectiveAttack,
		Decay:        p.Decay,
		Sustain:      p.Sustain,
		Release:      p.Release,
		AttackCurve:  p.AttackCurve,
		DecayCurve:   p.DecayCurve,
		ReleaseCurve: p.ReleaseCurve,
	}
}

// SynthesizeParams renders a single NoteEvent using xaudio.RenderNote.
func SynthesizeParams(p SynthParams) ([]int16, error) {
	cfg := BuildUnisonConfig(p)
	env := BuildADSR(p)
	note := xaudio.NoteEvent{
		StartSec:    0.0,
		DurationSec: p.LengthMs / 1000.0,
		Freq:        p.Tune,
		Velocity:    0.85 + (0.15 * p.Snappy),
		Config:      &cfg,
		Envelope:    &env,
	}
	return xaudio.RenderNote(cfg, env, p.SampleRate, note)
}

// SynthesizeSequence converts user-edited SequenceStep entries into xaudio.NoteEvent items
// and mixes them via xaudio.RenderSequence.
func SynthesizeSequence(p SynthParams, steps []SequenceStep) ([]int16, error) {
	cfg := BuildUnisonConfig(p)
	env := BuildADSR(p)

	var events []xaudio.NoteEvent
	for _, st := range steps {
		events = append(events, xaudio.NoteEvent{
			StartSec:    st.StartSec,
			DurationSec: st.DurationSec,
			Freq:        st.Freq,
			Velocity:    st.Velocity,
			Config:      &cfg,
			Envelope:    &env,
		})
	}

	return xaudio.RenderSequence(cfg, env, p.SampleRate, events)
}

// computeSpectrum performs an FFT on left-channel samples to yield normalized magnitudes [0.0, 1.0].
func computeSpectrum(stereoSamples []int16, numBins int) []float64 {
	fftSize := 1024
	x := make([]complex128, fftSize)
	for i := 0; i < fftSize && (i*2) < len(stereoSamples); i++ {
		window := 0.5 * (1.0 - math.Cos(2.0*math.Pi*float64(i)/float64(fftSize-1)))
		sample := (float64(stereoSamples[i*2]) / math.MaxInt16) * window
		x[i] = complex(sample, 0)
	}

	var fft func([]complex128) []complex128
	fft = func(a []complex128) []complex128 {
		n := len(a)
		if n <= 1 {
			return a
		}
		even := make([]complex128, n/2)
		odd := make([]complex128, n/2)
		for i := 0; i < n/2; i++ {
			even[i] = a[2*i]
			odd[i] = a[2*i+1]
		}
		even = fft(even)
		odd = fft(odd)
		out := make([]complex128, n)
		for k := 0; k < n/2; k++ {
			t := cmplx.Exp(complex(0, -2.0*math.Pi*float64(k)/float64(n))) * odd[k]
			out[k] = even[k] + t
			out[k+n/2] = even[k] - t
		}
		return out
	}

	spectrum := fft(x)
	half := fftSize / 2
	if numBins > half {
		numBins = half
	}
	magnitudes := make([]float64, numBins)
	maxMag := 1e-9
	step := half / numBins
	for b := 0; b < numBins; b++ {
		mag := cmplx.Abs(spectrum[b*step])
		magnitudes[b] = mag
		if mag > maxMag {
			maxMag = mag
		}
	}
	for b := range magnitudes {
		magnitudes[b] /= maxMag
	}
	return magnitudes
}

// playWavFile asynchronously executes the operating system's default audio player.
func playWavFile(path string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("afplay", path)
	case "windows":
		cmd = exec.Command("powershell", "-c", fmt.Sprintf("(New-Object Media.SoundPlayer '%s').PlaySync()", path))
	default:
		cmd = exec.Command("aplay", path)
	}
	_ = cmd.Start()
}

// resolveEnvelopeCurve translates curve string selections into xaudio.EnvelopeCurve callbacks.
func resolveEnvelopeCurve(name string, cycles, depth float64) xaudio.EnvelopeCurve {
	switch name {
	case "Exponential":
		return xaudio.EnvExponential
	case "Logarithmic":
		return xaudio.EnvLogarithmic
	case "Hermite":
		return xaudio.EnvHermite
	case "Radian":
		return xaudio.EnvRadian
	case "BicubicWobble":
		return xaudio.NewBicubicWobbleCurve(cycles, depth)
	default:
		return xaudio.EnvLinear
	}
}

// resolveDetuneCurve translates detune curve string selections into xaudio.DetuneCurve callbacks.
func resolveDetuneCurve(name string) xaudio.DetuneCurve {
	switch name {
	case "Analog Trance Synth":
		return xaudio.CurveJP8000
	case "Exponential":
		return xaudio.CurveExponential
	default:
		return xaudio.CurveLinear
	}
}

// resolveWaveform translates waveform string selections into xaudio.Waveform constants.
func resolveWaveform(name string) xaudio.Waveform {
	switch name {
	case "Sine":
		return xaudio.Sine
	case "Sawtooth":
		return xaudio.Sawtooth
	case "Square":
		return xaudio.Square
	case "SquarePoly":
		return xaudio.SquarePoly
	case "Triangle":
		return xaudio.Triangle
	case "FM":
		return xaudio.FM
	default:
		return xaudio.SawtoothPoly
	}
}

// synthUI consolidates all Fyne application state, widgets, and layout containers into one struct.
type synthUI struct {
	app           fyne.App
	window        fyne.Window
	params        SynthParams
	steps         []SequenceStep
	currentSample []int16
	vizMode       string
	vizLastSize   fyne.Size
	iniPath       string
	iniPresets    map[string]SynthParams
	applyingState bool

	// Status & visualization widgets
	statusLabel    *widget.Label
	vizContainer   *fyne.Container
	vizBg          *canvas.Rectangle
	vizGroupHeader *fyne.Container
	vizContent     *fyne.Container
	vizAccordion   *autoFitAccordion

	// Selectors, combobox & checkboxes
	modeSelect         *widget.Select
	presetCombo        *widget.SelectEntry
	savePresetBtn      *widget.Button
	loadPresetsBtn     *widget.Button
	deletePresetBtn    *widget.Button
	srSelect           *widget.Select
	waveSelect         *widget.Select
	detuneCurveSelect  *widget.Select
	subOctCheck        *widget.Check
	hardSyncCheck      *widget.Check
	attackCurveSelect  *widget.Select
	decayCurveSelect   *widget.Select
	releaseCurveSelect *widget.Select

	// Sliders and real-time value readouts
	lengthSlider       *widget.Slider
	lengthValLbl       *widget.Label
	tuneSlider         *widget.Slider
	tuneValLbl         *widget.Label
	toneSlider         *widget.Slider
	toneValLbl         *widget.Label
	snappySlider       *widget.Slider
	snappyValLbl       *widget.Label
	modRatioSlider     *widget.Slider
	modRatioValLbl     *widget.Label
	modIndexSlider     *widget.Slider
	modIndexValLbl     *widget.Label
	voicesSlider       *widget.Slider
	voicesValLbl       *widget.Label
	detuneSlider       *widget.Slider
	detuneValLbl       *widget.Label
	mixSlider          *widget.Slider
	mixValLbl          *widget.Label
	widthSlider        *widget.Slider
	widthValLbl        *widget.Label
	driftSlider        *widget.Slider
	driftValLbl        *widget.Label
	pulseWidthSlider   *widget.Slider
	pulseWidthValLbl   *widget.Label
	attackSlider       *widget.Slider
	attackValLbl       *widget.Label
	decaySlider        *widget.Slider
	decayValLbl        *widget.Label
	sustainSlider      *widget.Slider
	sustainValLbl      *widget.Label
	releaseSlider      *widget.Slider
	releaseValLbl      *widget.Label
	wobbleCyclesSlider *widget.Slider
	wobbleCyclesValLbl *widget.Label
	wobbleDepthSlider  *widget.Slider
	wobbleDepthValLbl  *widget.Label

	// Action buttons & metadata entries
	playBtn        *widget.Button
	playSeqBtn     *widget.Button
	saveBtn        *widget.Button
	resetBtn       *widget.Button
	sortStepsBtn   *widget.Button
	addStepBtn     *widget.Button
	saveSeqMetaBtn *widget.Button
	titleEntry     *widget.Entry
	artistEntry    *widget.Entry
	descEntry      *widget.Entry

	// Layout containers & independent accordions
	stepRowsContainer *fyne.Container
	seqContent        *fyne.Container
	seqAccordion      *autoFitAccordion
	metaContent       *fyne.Container
	metaAccordion     *autoFitAccordion
	colComposer       *fyne.Container
	colUnison         *fyne.Container
	colADSR           *fyne.Container
	topBar            *fyne.Container
	topBarLine1       *fyne.Container
	topBarLine2       *fyne.Container
}

// newSynthUI initializes defaults, loads XDG INI presets if present, and constructs all widgets.
func newSynthUI() *synthUI {
	a := app.New()
	t := &customSynthTheme{dark: true}
	a.Settings().SetTheme(t)

	w := a.NewWindow("SamplesComposer & Unison Synthesizer")
	w.SetFixedSize(false)

	iniPath := resolvePresetsIniPath()
	iniPresets, _ := loadPresetsFromINI(iniPath)
	if iniPresets == nil {
		iniPresets = make(map[string]SynthParams)
	}

	ui := &synthUI{
		app:        a,
		window:     w,
		vizMode:    "Waveform",
		params:     defaultSynthParams(),
		steps:      defaultSequenceSteps(),
		iniPath:    iniPath,
		iniPresets: iniPresets,
	}

	ui.initWidgets()
	ui.refreshPresetOptions()
	ui.buildLayout()
	ui.regenerateNote()
	ui.fitWindowToContent()
	return ui
}

// refreshPresetOptions merges hardcoded presets with any presets loaded from presets.ini.
func (ui *synthUI) refreshPresetOptions() {
	if ui.presetCombo == nil {
		return
	}
	seen := make(map[string]bool)
	var options []string
	for _, name := range hardcodedPresetNames {
		options = append(options, name)
		seen[name] = true
	}
	var iniNames []string
	for name := range ui.iniPresets {
		if !seen[name] {
			iniNames = append(iniNames, name)
			seen[name] = true
		}
	}
	sort.Strings(iniNames)
	options = append(options, iniNames...)
	ui.presetCombo.SetOptions(options)
	ui.presetCombo.Refresh()
}

// applyParamsToControls synchronizes all slider, selector, and checkbox widgets to p.
func (ui *synthUI) applyParamsToControls(p SynthParams) {
	ui.applyingState = true
	defer func() { ui.applyingState = false }()

	ui.params = p

	switch p.SampleRate {
	case 8000:
		ui.srSelect.SetSelected("8 kHz")
	case 16000:
		ui.srSelect.SetSelected("16 kHz")
	case 32000:
		ui.srSelect.SetSelected("32 kHz")
	default:
		ui.srSelect.SetSelected("48 kHz")
	}

	ui.lengthSlider.SetValue(p.LengthMs)
	ui.lengthValLbl.SetText(fmt.Sprintf("%.0f ms", p.LengthMs))

	ui.tuneSlider.SetValue(p.Tune)
	ui.tuneValLbl.SetText(fmt.Sprintf("%.2f Hz", p.Tune))

	ui.toneSlider.SetValue(p.Tone)
	ui.toneValLbl.SetText(fmt.Sprintf("%.2f", p.Tone))

	ui.snappySlider.SetValue(p.Snappy)
	ui.snappyValLbl.SetText(fmt.Sprintf("%.2f", p.Snappy))

	ui.modRatioSlider.SetValue(p.ModRatio)
	ui.modRatioValLbl.SetText(fmt.Sprintf("%.2f", p.ModRatio))

	ui.modIndexSlider.SetValue(p.ModIndex)
	ui.modIndexValLbl.SetText(fmt.Sprintf("%.2f", p.ModIndex))

	ui.waveSelect.SetSelected(formatWaveformName(p.BaseWave))
	ui.detuneCurveSelect.SetSelected(formatDetuneCurveName(&p.DetuneCurve))

	ui.voicesSlider.SetValue(float64(p.Voices))
	ui.voicesValLbl.SetText(fmt.Sprintf("%d", p.Voices))

	ui.detuneSlider.SetValue(p.Detune)
	ui.detuneValLbl.SetText(fmt.Sprintf("%.2f", p.Detune))

	ui.mixSlider.SetValue(p.Mix)
	ui.mixValLbl.SetText(fmt.Sprintf("%.2f", p.Mix))

	ui.widthSlider.SetValue(p.StereoWidth)
	ui.widthValLbl.SetText(fmt.Sprintf("%.2f", p.StereoWidth))

	ui.driftSlider.SetValue(p.AnalogDrift)
	ui.driftValLbl.SetText(fmt.Sprintf("%.2f", p.AnalogDrift))

	ui.pulseWidthSlider.SetValue(p.PulseWidth)
	ui.pulseWidthValLbl.SetText(fmt.Sprintf("%.2f", p.PulseWidth))

	ui.subOctCheck.SetChecked(p.SubOctave)
	ui.hardSyncCheck.SetChecked(p.HardSync)

	ui.attackSlider.SetValue(p.Attack)
	ui.attackValLbl.SetText(fmt.Sprintf("%.3fs", p.Attack))
	ui.attackCurveSelect.SetSelected(formatEnvelopeCurveName(p.AttackCurve, "Hermite"))

	ui.decaySlider.SetValue(p.Decay)
	ui.decayValLbl.SetText(fmt.Sprintf("%.2fs", p.Decay))
	ui.decayCurveSelect.SetSelected(formatEnvelopeCurveName(p.DecayCurve, "Exponential"))

	ui.sustainSlider.SetValue(p.Sustain)
	ui.sustainValLbl.SetText(fmt.Sprintf("%.2f", p.Sustain))

	ui.releaseSlider.SetValue(p.Release)
	ui.releaseValLbl.SetText(fmt.Sprintf("%.2fs", p.Release))
	ui.releaseCurveSelect.SetSelected(formatEnvelopeCurveName(p.ReleaseCurve, "Logarithmic"))

	ui.wobbleCyclesSlider.SetValue(p.WobbleCycles)
	ui.wobbleCyclesValLbl.SetText(fmt.Sprintf("%.1f", p.WobbleCycles))

	ui.wobbleDepthSlider.SetValue(p.WobbleDepth)
	ui.wobbleDepthValLbl.SetText(fmt.Sprintf("%.2f", p.WobbleDepth))
}

// applyPreset updates synthesizer parameters from either a hardcoded preset or an INI preset.
func (ui *synthUI) applyPreset(name string) {
	if ui.applyingState {
		return
	}
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return
	}

	p := ui.params
	var sourceLabel string

	switch trimmed {
	case "Analog Trance Synth":
		p.Voices = 7
		p.Detune = 0.35
		p.Mix = 0.75
		p.StereoWidth = 0.0
		p.SubOctave = false
		p.BaseWave = xaudio.SawtoothPoly
		p.DetuneCurve = xaudio.CurveJP8000
		sourceLabel = "Internal Preset: Analog Trance Synth"
	case "Hypersaw":
		p.Voices = 9
		p.Detune = 0.45
		p.Mix = 0.85
		p.StereoWidth = 0.90
		p.SubOctave = true
		p.BaseWave = xaudio.SawtoothPoly
		p.DetuneCurve = xaudio.CurveJP8000
		sourceLabel = "Internal Preset: Hypersaw"
	case "MultiSaw":
		p.Voices = 8
		p.Detune = 0.30
		p.Mix = 0.80
		p.StereoWidth = 0.80
		p.SubOctave = false
		p.BaseWave = xaudio.SawtoothPoly
		p.DetuneCurve = xaudio.CurveLinear
		sourceLabel = "Internal Preset: MultiSaw"
	default:
		loaded, ok := ui.iniPresets[trimmed]
		if !ok {
			if fresh, err := loadPresetsFromINI(ui.iniPath); err == nil {
				ui.iniPresets = fresh
				ui.refreshPresetOptions()
				loaded, ok = ui.iniPresets[trimmed]
			}
		}
		if !ok {
			return
		}
		p = loaded
		sourceLabel = fmt.Sprintf("INI Preset: %s", trimmed)
	}

	ui.applyParamsToControls(p)
	ui.statusLabel.SetText("Loaded " + sourceLabel)
	ui.regenerateNote()
}

// reloadPresetsFromConfig reloads presets.ini from the XDG user config directory.
func (ui *synthUI) reloadPresetsFromConfig() {
	presets, err := loadPresetsFromINI(ui.iniPath)
	if err != nil {
		if os.IsNotExist(err) {
			ui.iniPresets = make(map[string]SynthParams)
			ui.refreshPresetOptions()
			ui.statusLabel.SetText("No presets.ini yet at " + ui.iniPath + " (hardcoded presets active)")
			return
		}
		ui.statusLabel.SetText("INI load error: " + err.Error())
		return
	}
	ui.iniPresets = presets
	ui.refreshPresetOptions()
	ui.statusLabel.SetText(fmt.Sprintf("Loaded %d INI preset(s) from %s", len(presets), ui.iniPath))
}

// saveCurrentPreset saves current synthesizer parameters under the combo box text into presets.ini.
func (ui *synthUI) saveCurrentPreset() {
	name := strings.TrimSpace(ui.presetCombo.Text)
	if name == "" {
		name = "Custom Preset"
		ui.presetCombo.SetText(name)
	}
	if err := savePresetToINI(ui.iniPath, name, ui.params); err != nil {
		ui.statusLabel.SetText("Save preset error: " + err.Error())
		return
	}
	ui.iniPresets[name] = ui.params
	ui.refreshPresetOptions()
	ui.presetCombo.SetText(name)
	ui.statusLabel.SetText(fmt.Sprintf("Saved preset [%s] to %s", name, ui.iniPath))
}

// deleteSelectedPreset removes the selected/entered INI preset from presets.ini in the XDG config dir.
func (ui *synthUI) deleteSelectedPreset() {
	name := strings.TrimSpace(ui.presetCombo.Text)
	if name == "" {
		ui.statusLabel.SetText("Select or enter an INI preset name to delete")
		return
	}
	if isHardcodedPreset(name) {
		ui.statusLabel.SetText(fmt.Sprintf("Preset [%s] is hardcoded and cannot be deleted from presets.ini", name))
		return
	}

	deleted, err := removePresetFromINI(ui.iniPath, name)
	if err != nil {
		ui.statusLabel.SetText("Delete preset error: " + err.Error())
		return
	}
	if !deleted {
		ui.statusLabel.SetText(fmt.Sprintf("Preset [%s] not found in %s", name, ui.iniPath))
		return
	}

	delete(ui.iniPresets, name)
	ui.refreshPresetOptions()
	ui.presetCombo.SetText("")
	ui.statusLabel.SetText(fmt.Sprintf("Deleted INI preset [%s] from %s", name, ui.iniPath))
}

// fitWindowToContent sizes the window so all visible widgets fit without scrollbars.
func (ui *synthUI) fitWindowToContent() {
	if ui.window == nil || ui.window.Content() == nil {
		return
	}
	min := ui.window.Content().MinSize()
	cur := ui.window.Content().Size()
	targetW := cur.Width
	if min.Width > targetW {
		targetW = min.Width
	}
	ui.window.Resize(fyne.NewSize(targetW, min.Height))
}

// resetAllToDefaults restores all SynthParams, widget controls, and timeline steps to initial defaults.
func (ui *synthUI) resetAllToDefaults() {
	ui.steps = defaultSequenceSteps()
	if ui.presetCombo != nil {
		ui.presetCombo.SetText("")
	}
	ui.applyParamsToControls(defaultSynthParams())
	ui.rebuildTimelineEditor()
	ui.regenerateNote()
	ui.fitWindowToContent()
	ui.statusLabel.SetText("Reset all synthesizer parameters and timeline steps to defaults")
}

// redrawVisualization updates the waveform or FFT spectrum canvas using the current container width.
func (ui *synthUI) redrawVisualization() {
	ui.vizBg.FillColor = color.RGBA{R: 18, G: 20, B: 26, A: 255}
	ui.vizContainer.Objects = []fyne.CanvasObject{ui.vizBg}

	width := ui.vizContainer.Size().Width
	height := ui.vizContainer.Size().Height
	if width <= 1 && ui.window.Content() != nil {
		width = ui.window.Content().Size().Width - 16
	}
	if width <= 1 {
		width = 1000
	}
	if height <= 1 {
		height = 160
	}
	ui.vizBg.Resize(fyne.NewSize(width, height))

	if len(ui.currentSample) == 0 {
		ui.vizContainer.Refresh()
		return
	}

	waveColor := color.RGBA{R: 0, G: 220, B: 180, A: 255}
	specColor := color.RGBA{R: 255, G: 140, B: 40, A: 255}

	if ui.vizMode == "Waveform" {
		numPoints := int(width / 3)
		if numPoints < 60 {
			numPoints = 60
		}
		if numPoints > 600 {
			numPoints = 600
		}
		step := len(ui.currentSample) / 2 / numPoints
		if step < 1 {
			step = 1
		}
		midY := height / 2.0
		dx := width / float32(numPoints)
		var prevX, prevY float32
		for i := 0; i < numPoints && (i*step*2) < len(ui.currentSample); i++ {
			norm := float32(ui.currentSample[i*step*2]) / float32(math.MaxInt16)
			x := float32(i) * dx
			y := midY - (norm * (midY - 10))
			if i > 0 {
				line := canvas.NewLine(waveColor)
				line.StrokeWidth = 2
				line.Position1 = fyne.NewPos(prevX, prevY)
				line.Position2 = fyne.NewPos(x, y)
				ui.vizContainer.Add(line)
			}
			prevX, prevY = x, y
		}
	} else {
		bins := int(width / 16)
		if bins < 16 {
			bins = 16
		}
		if bins > 128 {
			bins = 128
		}
		mags := computeSpectrum(ui.currentSample, bins)
		barW := width / float32(bins)
		for b, mag := range mags {
			barH := float32(mag) * (height - 16)
			bar := canvas.NewRectangle(specColor)
			bar.Move(fyne.NewPos(float32(b)*barW+1, height-barH))
			bar.Resize(fyne.NewSize(barW-2, barH))
			ui.vizContainer.Add(bar)
		}
	}
	ui.vizContainer.Refresh()
}

// regenerateNote renders a single note via SynthesizeParams and redraws the visualization.
func (ui *synthUI) regenerateNote() {
	samples, err := SynthesizeParams(ui.params)
	if err != nil {
		ui.statusLabel.SetText("Synthesis error: " + err.Error())
		return
	}
	ui.currentSample = samples
	ui.statusLabel.SetText(fmt.Sprintf("Rendered %d stereo frames @ %d Hz", len(samples)/2, ui.params.SampleRate))
	ui.redrawVisualization()
}

// regenerateSequence renders the timeline steps via SynthesizeSequence and redraws the visualization.
func (ui *synthUI) regenerateSequence() {
	seqSamples, err := SynthesizeSequence(ui.params, ui.steps)
	if err != nil {
		ui.statusLabel.SetText("Sequence error: " + err.Error())
		return
	}
	ui.currentSample = seqSamples
	ui.statusLabel.SetText(fmt.Sprintf("Rendered sequence (%d steps, %d stereo frames)", len(ui.steps), len(seqSamples)/2))
	ui.redrawVisualization()
}

// buildStepRow creates a single timeline step row with full-width sliders and a Remove button.
func (ui *synthUI) buildStepRow(idx int) *fyne.Container {
	st := &ui.steps[idx]

	startLbl := widget.NewLabel(fmt.Sprintf("%.2fs", st.StartSec))
	startSlider := widget.NewSlider(0.0, 8.0)
	startSlider.Step = 0.05
	startSlider.SetValue(st.StartSec)

	freqLbl := widget.NewLabel(fmt.Sprintf("%.1f Hz", st.Freq))
	freqSlider := widget.NewSlider(30.0, 880.0)
	freqSlider.Step = 1.0
	freqSlider.SetValue(st.Freq)

	durLbl := widget.NewLabel(fmt.Sprintf("%.2fs", st.DurationSec))
	durSlider := widget.NewSlider(0.05, 3.0)
	durSlider.Step = 0.05
	durSlider.SetValue(st.DurationSec)

	velLbl := widget.NewLabel(fmt.Sprintf("%.2f", st.Velocity))
	velSlider := widget.NewSlider(0.05, 1.0)
	velSlider.Step = 0.05
	velSlider.SetValue(st.Velocity)

	startSlider.OnChanged = func(v float64) {
		st.StartSec = v
		startLbl.SetText(fmt.Sprintf("%.2fs", v))
		ui.regenerateSequence()
	}
	freqSlider.OnChanged = func(v float64) {
		st.Freq = v
		freqLbl.SetText(fmt.Sprintf("%.1f Hz", v))
		ui.regenerateSequence()
	}
	durSlider.OnChanged = func(v float64) {
		st.DurationSec = v
		durLbl.SetText(fmt.Sprintf("%.2fs", v))
		ui.regenerateSequence()
	}
	velSlider.OnChanged = func(v float64) {
		st.Velocity = v
		velLbl.SetText(fmt.Sprintf("%.2f", v))
		ui.regenerateSequence()
	}

	removeBtn := widget.NewButton("Remove", func() {
		if len(ui.steps) <= 1 {
			ui.statusLabel.SetText("Sequence must contain at least 1 step")
			return
		}
		ui.steps = append(ui.steps[:idx], ui.steps[idx+1:]...)
		ui.rebuildTimelineEditor()
		ui.regenerateSequence()
		ui.fitWindowToContent()
	})

	slidersGrid := container.NewGridWithColumns(4,
		makeParameterRow(fmt.Sprintf("Step %d Start", idx+1), startSlider, startLbl),
		makeParameterRow("Freq", freqSlider, freqLbl),
		makeParameterRow("Duration", durSlider, durLbl),
		makeParameterRow("Velocity", velSlider, velLbl),
	)

	return container.NewBorder(nil, nil, nil, removeBtn, slidersGrid)
}

// rebuildTimelineEditor refreshes all rows inside stepRowsContainer.
func (ui *synthUI) rebuildTimelineEditor() {
	ui.stepRowsContainer.Objects = nil
	for i := range ui.steps {
		ui.stepRowsContainer.Add(ui.buildStepRow(i))
	}
	ui.stepRowsContainer.Refresh()
}

// initWidgets instantiates all interactive controls on synthUI.
func (ui *synthUI) initWidgets() {
	ui.statusLabel = widget.NewLabel("Ready")

	ui.vizBg = canvas.NewRectangle(color.RGBA{R: 18, G: 20, B: 26, A: 255})
	ui.vizContainer = container.New(&stretchVizLayout{ui: ui}, ui.vizBg)

	ui.modeSelect = widget.NewSelect([]string{"Waveform", "Spectrum"}, func(s string) {
		ui.vizMode = s
		ui.redrawVisualization()
	})
	ui.modeSelect.SetSelected("Waveform")

	ui.presetCombo = widget.NewSelectEntry(hardcodedPresetNames)
	ui.presetCombo.PlaceHolder = "Select or type preset name..."
	ui.presetCombo.OnChanged = func(s string) {
		ui.applyPreset(s)
	}
	ui.presetCombo.OnSubmitted = func(s string) {
		ui.applyPreset(s)
	}

	ui.savePresetBtn = widget.NewButton("Save INI", func() {
		ui.saveCurrentPreset()
	})
	ui.loadPresetsBtn = widget.NewButton("Reload INI", func() {
		ui.reloadPresetsFromConfig()
	})
	ui.deletePresetBtn = widget.NewButton("Delete INI", func() {
		ui.deleteSelectedPreset()
	})

	ui.srSelect = widget.NewSelect([]string{"8 kHz", "16 kHz", "32 kHz", "48 kHz"}, func(s string) {
		khz, err := strconv.Atoi(strings.TrimSuffix(s, " kHz"))
		if err == nil {
			ui.params.SampleRate = khz * 1000
			if !ui.applyingState {
				ui.regenerateNote()
			}
		}
	})
	ui.srSelect.SetSelected("48 kHz")

	ui.lengthValLbl = widget.NewLabel(fmt.Sprintf("%.0f ms", ui.params.LengthMs))
	ui.lengthSlider = widget.NewSlider(1000, 30000)
	ui.lengthSlider.Step = 100
	ui.lengthSlider.SetValue(ui.params.LengthMs)
	ui.lengthSlider.OnChanged = func(v float64) {
		ui.params.LengthMs = v
		ui.lengthValLbl.SetText(fmt.Sprintf("%.0f ms", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.tuneValLbl = widget.NewLabel(fmt.Sprintf("%.2f Hz", ui.params.Tune))
	ui.tuneSlider = widget.NewSlider(30.0, 880.0)
	ui.tuneSlider.SetValue(ui.params.Tune)
	ui.tuneSlider.OnChanged = func(v float64) {
		ui.params.Tune = v
		ui.tuneValLbl.SetText(fmt.Sprintf("%.2f Hz", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.toneValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.Tone))
	ui.toneSlider = widget.NewSlider(0.0, 1.0)
	ui.toneSlider.Step = 0.01
	ui.toneSlider.SetValue(ui.params.Tone)
	ui.toneSlider.OnChanged = func(v float64) {
		ui.params.Tone = v
		ui.toneValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.snappyValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.Snappy))
	ui.snappySlider = widget.NewSlider(0.0, 1.0)
	ui.snappySlider.Step = 0.01
	ui.snappySlider.SetValue(ui.params.Snappy)
	ui.snappySlider.OnChanged = func(v float64) {
		ui.params.Snappy = v
		ui.snappyValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.waveSelect = widget.NewSelect([]string{"Sine", "Sawtooth", "SawtoothPoly", "Square", "SquarePoly", "Triangle", "FM"}, func(s string) {
		ui.params.BaseWave = resolveWaveform(s)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})
	ui.waveSelect.SetSelected("SawtoothPoly")

	ui.detuneCurveSelect = widget.NewSelect([]string{"Linear", "Analog Trance Synth", "Exponential"}, func(s string) {
		ui.params.DetuneCurve = resolveDetuneCurve(s)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})
	ui.detuneCurveSelect.SetSelected("Analog Trance Synth")

	ui.voicesValLbl = widget.NewLabel(fmt.Sprintf("%d", ui.params.Voices))
	ui.voicesSlider = widget.NewSlider(1, 16)
	ui.voicesSlider.Step = 1
	ui.voicesSlider.SetValue(float64(ui.params.Voices))
	ui.voicesSlider.OnChanged = func(v float64) {
		ui.params.Voices = int(v)
		ui.voicesValLbl.SetText(fmt.Sprintf("%d", ui.params.Voices))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.detuneValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.Detune))
	ui.detuneSlider = widget.NewSlider(0.0, 1.0)
	ui.detuneSlider.Step = 0.01
	ui.detuneSlider.SetValue(ui.params.Detune)
	ui.detuneSlider.OnChanged = func(v float64) {
		ui.params.Detune = v
		ui.detuneValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.mixValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.Mix))
	ui.mixSlider = widget.NewSlider(0.0, 1.0)
	ui.mixSlider.Step = 0.01
	ui.mixSlider.SetValue(ui.params.Mix)
	ui.mixSlider.OnChanged = func(v float64) {
		ui.params.Mix = v
		ui.mixValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.widthValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.StereoWidth))
	ui.widthSlider = widget.NewSlider(0.0, 1.0)
	ui.widthSlider.Step = 0.01
	ui.widthSlider.SetValue(ui.params.StereoWidth)
	ui.widthSlider.OnChanged = func(v float64) {
		ui.params.StereoWidth = v
		ui.widthValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.driftValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.AnalogDrift))
	ui.driftSlider = widget.NewSlider(0.0, 1.0)
	ui.driftSlider.Step = 0.01
	ui.driftSlider.SetValue(ui.params.AnalogDrift)
	ui.driftSlider.OnChanged = func(v float64) {
		ui.params.AnalogDrift = v
		ui.driftValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.pulseWidthValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.PulseWidth))
	ui.pulseWidthSlider = widget.NewSlider(0.05, 0.95)
	ui.pulseWidthSlider.Step = 0.01
	ui.pulseWidthSlider.SetValue(ui.params.PulseWidth)
	ui.pulseWidthSlider.OnChanged = func(v float64) {
		ui.params.PulseWidth = v
		ui.pulseWidthValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.modRatioValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.ModRatio))
	ui.modRatioSlider = widget.NewSlider(0.5, 10.0)
	ui.modRatioSlider.Step = 0.1
	ui.modRatioSlider.SetValue(ui.params.ModRatio)
	ui.modRatioSlider.OnChanged = func(v float64) {
		ui.params.ModRatio = v
		ui.modRatioValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.modIndexValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.ModIndex))
	ui.modIndexSlider = widget.NewSlider(0.0, 10.0)
	ui.modIndexSlider.Step = 0.1
	ui.modIndexSlider.SetValue(ui.params.ModIndex)
	ui.modIndexSlider.OnChanged = func(v float64) {
		ui.params.ModIndex = v
		ui.modIndexValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.subOctCheck = widget.NewCheck("Sub-Octave (-1 Oct)", func(b bool) {
		ui.params.SubOctave = b
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})
	ui.subOctCheck.SetChecked(ui.params.SubOctave)

	ui.hardSyncCheck = widget.NewCheck("Hard Sync", func(b bool) {
		ui.params.HardSync = b
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})

	curveOptions := []string{"Linear", "Exponential", "Logarithmic", "Hermite", "Radian", "BicubicWobble"}

	ui.attackValLbl = widget.NewLabel(fmt.Sprintf("%.3fs", ui.params.Attack))
	ui.attackSlider = widget.NewSlider(0.001, 2.0)
	ui.attackSlider.Step = 0.005
	ui.attackSlider.SetValue(ui.params.Attack)
	ui.attackSlider.OnChanged = func(v float64) {
		ui.params.Attack = v
		ui.attackValLbl.SetText(fmt.Sprintf("%.3fs", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.decayValLbl = widget.NewLabel(fmt.Sprintf("%.2fs", ui.params.Decay))
	ui.decaySlider = widget.NewSlider(0.01, 3.0)
	ui.decaySlider.Step = 0.01
	ui.decaySlider.SetValue(ui.params.Decay)
	ui.decaySlider.OnChanged = func(v float64) {
		ui.params.Decay = v
		ui.decayValLbl.SetText(fmt.Sprintf("%.2fs", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.sustainValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.Sustain))
	ui.sustainSlider = widget.NewSlider(0.0, 1.0)
	ui.sustainSlider.Step = 0.01
	ui.sustainSlider.SetValue(ui.params.Sustain)
	ui.sustainSlider.OnChanged = func(v float64) {
		ui.params.Sustain = v
		ui.sustainValLbl.SetText(fmt.Sprintf("%.2f", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.releaseValLbl = widget.NewLabel(fmt.Sprintf("%.2fs", ui.params.Release))
	ui.releaseSlider = widget.NewSlider(0.01, 3.0)
	ui.releaseSlider.Step = 0.01
	ui.releaseSlider.SetValue(ui.params.Release)
	ui.releaseSlider.OnChanged = func(v float64) {
		ui.params.Release = v
		ui.releaseValLbl.SetText(fmt.Sprintf("%.2fs", v))
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.attackCurveSelect = widget.NewSelect(curveOptions, func(s string) {
		ui.params.AttackCurve = resolveEnvelopeCurve(s, ui.params.WobbleCycles, ui.params.WobbleDepth)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})
	ui.attackCurveSelect.SetSelected("Hermite")

	ui.decayCurveSelect = widget.NewSelect(curveOptions, func(s string) {
		ui.params.DecayCurve = resolveEnvelopeCurve(s, ui.params.WobbleCycles, ui.params.WobbleDepth)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})
	ui.decayCurveSelect.SetSelected("Exponential")

	ui.releaseCurveSelect = widget.NewSelect(curveOptions, func(s string) {
		ui.params.ReleaseCurve = resolveEnvelopeCurve(s, ui.params.WobbleCycles, ui.params.WobbleDepth)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	})
	ui.releaseCurveSelect.SetSelected("Logarithmic")

	ui.wobbleCyclesValLbl = widget.NewLabel(fmt.Sprintf("%.1f", ui.params.WobbleCycles))
	ui.wobbleCyclesSlider = widget.NewSlider(0.5, 10.0)
	ui.wobbleCyclesSlider.Step = 0.5
	ui.wobbleCyclesSlider.SetValue(ui.params.WobbleCycles)
	ui.wobbleCyclesSlider.OnChanged = func(v float64) {
		ui.params.WobbleCycles = v
		ui.wobbleCyclesValLbl.SetText(fmt.Sprintf("%.1f", v))
		ui.params.AttackCurve = resolveEnvelopeCurve(ui.attackCurveSelect.Selected, ui.params.WobbleCycles, ui.params.WobbleDepth)
		ui.params.DecayCurve = resolveEnvelopeCurve(ui.decayCurveSelect.Selected, ui.params.WobbleCycles, ui.params.WobbleDepth)
		ui.params.ReleaseCurve = resolveEnvelopeCurve(ui.releaseCurveSelect.Selected, ui.params.WobbleCycles, ui.params.WobbleDepth)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.wobbleDepthValLbl = widget.NewLabel(fmt.Sprintf("%.2f", ui.params.WobbleDepth))
	ui.wobbleDepthSlider = widget.NewSlider(0.0, 0.5)
	ui.wobbleDepthSlider.Step = 0.01
	ui.wobbleDepthSlider.SetValue(ui.params.WobbleDepth)
	ui.wobbleDepthSlider.OnChanged = func(v float64) {
		ui.params.WobbleDepth = v
		ui.wobbleDepthValLbl.SetText(fmt.Sprintf("%.2f", v))
		ui.params.AttackCurve = resolveEnvelopeCurve(ui.attackCurveSelect.Selected, ui.params.WobbleCycles, ui.params.WobbleDepth)
		ui.params.DecayCurve = resolveEnvelopeCurve(ui.decayCurveSelect.Selected, ui.params.WobbleCycles, ui.params.WobbleDepth)
		ui.params.ReleaseCurve = resolveEnvelopeCurve(ui.releaseCurveSelect.Selected, ui.params.WobbleCycles, ui.params.WobbleDepth)
		if !ui.applyingState {
			ui.regenerateNote()
		}
	}

	ui.titleEntry = widget.NewEntry()
	ui.titleEntry.SetText("Sample Composer Unison Sequence")
	ui.artistEntry = widget.NewEntry()
	ui.artistEntry.SetText("xaudio Synth")
	ui.descEntry = widget.NewEntry()
	ui.descEntry.SetText("Generated with gommons xaudio RenderSequence")

	ui.stepRowsContainer = container.NewVBox()
	ui.sortStepsBtn = widget.NewButton("Sort by StartSec", func() {
		sort.SliceStable(ui.steps, func(i, j int) bool {
			return ui.steps[i].StartSec < ui.steps[j].StartSec
		})
		ui.rebuildTimelineEditor()
		ui.statusLabel.SetText("Sorted sequence steps chronologically by StartSec")
		ui.regenerateSequence()
	})

	ui.addStepBtn = widget.NewButton("+ Add Step", func() {
		nextStart := 0.0
		nextFreq := ui.params.Tune
		if len(ui.steps) > 0 {
			last := ui.steps[len(ui.steps)-1]
			nextStart = math.Min(8.0, last.StartSec+0.5)
			nextFreq = math.Min(880.0, last.Freq*1.122)
		}
		ui.steps = append(ui.steps, SequenceStep{
			StartSec:    nextStart,
			Freq:        nextFreq,
			DurationSec: 0.45,
			Velocity:    0.85,
		})
		ui.rebuildTimelineEditor()
		ui.regenerateSequence()
		ui.fitWindowToContent()
	})

	ui.saveSeqMetaBtn = widget.NewButton("Save Sequence WAV", func() {
		ui.regenerateSequence()
		meta := xaudio.WavMetadata{
			Title:       ui.titleEntry.Text,
			Artist:      ui.artistEntry.Text,
			Description: ui.descEntry.Text,
			Software:    "gommons xaudio Fyne Synth",
			Date:        "2026-11-19",
		}
		outFile := "sequence_timeline_output.wav"
		if err := xaudio.WriteWavWithMetadata(outFile, ui.params.SampleRate, 2, ui.currentSample, meta); err != nil {
			ui.statusLabel.SetText("Sequence WAV save error: " + err.Error())
			return
		}
		ui.statusLabel.SetText("Saved rendered sequence to " + outFile + " with RIFF INFO metadata")
	})

	ui.rebuildTimelineEditor()

	ui.playBtn = widget.NewButton("Play Note", func() {
		ui.regenerateNote()
		tmpPath := "preview_note.wav"
		if err := xaudio.WriteWav(tmpPath, ui.params.SampleRate, 2, ui.currentSample); err != nil {
			ui.statusLabel.SetText("Playback write error: " + err.Error())
			return
		}
		playWavFile(tmpPath)
		ui.statusLabel.SetText("Playing " + tmpPath)
	})

	ui.playSeqBtn = widget.NewButton("Play Sequence", func() {
		ui.regenerateSequence()
		tmpPath := "preview_sequence.wav"
		if err := xaudio.WriteWav(tmpPath, ui.params.SampleRate, 2, ui.currentSample); err != nil {
			ui.statusLabel.SetText("Sequence write error: " + err.Error())
			return
		}
		playWavFile(tmpPath)
		ui.statusLabel.SetText("Playing edited timeline sequence (" + tmpPath + ")")
	})

	ui.saveBtn = widget.NewButton("Save WAV with Metadata", func() {
		meta := xaudio.WavMetadata{
			Title:       ui.titleEntry.Text,
			Artist:      ui.artistEntry.Text,
			Description: ui.descEntry.Text,
			Software:    "gommons xaudio Fyne Synth",
			Date:        "2026-11-19",
		}
		outFile := "sample_composer_output.wav"
		if err := xaudio.WriteWavWithMetadata(outFile, ui.params.SampleRate, 2, ui.currentSample, meta); err != nil {
			ui.statusLabel.SetText("Save error: " + err.Error())
			return
		}
		ui.statusLabel.SetText("Saved " + outFile + " with RIFF INFO metadata")
	})

	ui.resetBtn = widget.NewButton("Reset Defaults", func() {
		ui.resetAllToDefaults()
	})

}

// makeParameterRow places left label, stretching center input widget, and right numeric label
// in a Border container so the input widget expands across all available horizontal space.
func makeParameterRow(labelText string, input fyne.CanvasObject, valLabel *widget.Label) *fyne.Container {
	var right fyne.CanvasObject = layout.NewSpacer()
	if valLabel != nil {
		right = valLabel
	}
	return container.NewBorder(nil, nil, widget.NewLabel(labelText), right, input)
}

// buildLayout arranges the toolbar, parameter columns, independently collapsible Sequence,
// Meta-Data, and Visualization accordions directly into the window without scrollbars.
func (ui *synthUI) buildLayout() {
	ui.colComposer = container.NewVBox(
		widget.NewLabelWithStyle("Sample Composer & Format", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeParameterRow("Sample Rate", ui.srSelect, nil),
		makeParameterRow("Sample Length", ui.lengthSlider, ui.lengthValLbl),
		makeParameterRow("Tune", ui.tuneSlider, ui.tuneValLbl),
		makeParameterRow("Tone", ui.toneSlider, ui.toneValLbl),
		makeParameterRow("Snappy", ui.snappySlider, ui.snappyValLbl),
		makeParameterRow("FM Mod Ratio", ui.modRatioSlider, ui.modRatioValLbl),
		makeParameterRow("FM Mod Index", ui.modIndexSlider, ui.modIndexValLbl),
	)

	ui.colUnison = container.NewVBox(
		widget.NewLabelWithStyle("xaudio Unison & DetuneCurve", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeParameterRow("Waveform", ui.waveSelect, nil),
		makeParameterRow("Detune Curve", ui.detuneCurveSelect, nil),
		makeParameterRow("Voices", ui.voicesSlider, ui.voicesValLbl),
		makeParameterRow("Detune", ui.detuneSlider, ui.detuneValLbl),
		makeParameterRow("Mix", ui.mixSlider, ui.mixValLbl),
		makeParameterRow("Stereo Width", ui.widthSlider, ui.widthValLbl),
		makeParameterRow("Analog Drift", ui.driftSlider, ui.driftValLbl),
		makeParameterRow("Pulse Width", ui.pulseWidthSlider, ui.pulseWidthValLbl),
		container.NewHBox(ui.subOctCheck, ui.hardSyncCheck),
	)

	ui.colADSR = container.NewVBox(
		widget.NewLabelWithStyle("xaudio ADSR & EnvelopeCurve", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		makeParameterRow("Attack", ui.attackSlider, ui.attackValLbl),
		makeParameterRow("Attack Curve", ui.attackCurveSelect, nil),
		makeParameterRow("Decay", ui.decaySlider, ui.decayValLbl),
		makeParameterRow("Decay Curve", ui.decayCurveSelect, nil),
		makeParameterRow("Sustain", ui.sustainSlider, ui.sustainValLbl),
		makeParameterRow("Release", ui.releaseSlider, ui.releaseValLbl),
		makeParameterRow("Release Curve", ui.releaseCurveSelect, nil),
		makeParameterRow("Wobble Cycles", ui.wobbleCyclesSlider, ui.wobbleCyclesValLbl),
		makeParameterRow("Wobble Depth", ui.wobbleDepthSlider, ui.wobbleDepthValLbl),
	)

	// Standalone Sequence accordion that auto-fits the window when opened or closed
	ui.seqContent = container.NewVBox(
		container.NewHBox(
			layout.NewSpacer(),
			ui.sortStepsBtn,
			ui.addStepBtn,
			ui.saveSeqMetaBtn,
		),
		ui.stepRowsContainer,
	)
	ui.seqAccordion = newAutoFitAccordion(
		func() { ui.fitWindowToContent() },
		"Sequence",
		ui.seqContent,
	)
	ui.seqAccordion.Open(0)

	// Standalone Meta-Data accordion that auto-fits the window when opened or closed
	ui.metaContent = container.NewVBox(
		makeParameterRow("Title (INAM)", ui.titleEntry, nil),
		makeParameterRow("Artist (IART)", ui.artistEntry, nil),
		makeParameterRow("Description (ICMT)", ui.descEntry, nil),
	)
	ui.metaAccordion = newAutoFitAccordion(
		func() { ui.fitWindowToContent() },
		"Meta-Data",
		ui.metaContent,
	)
	ui.metaAccordion.Open(0)

	// Group label row paired directly with the visualization mode selector
	ui.vizGroupHeader = container.NewHBox(
		widget.NewLabelWithStyle("Visualization", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		ui.modeSelect,
		layout.NewSpacer(),
	)

	// Collapsible Visualization accordion that auto-fits the window when opened or closed
	ui.vizContent = container.NewVBox(
		ui.vizGroupHeader,
		ui.vizContainer,
	)
	ui.vizAccordion = newAutoFitAccordion(
		func() { ui.fitWindowToContent() },
		"Visualization",
		ui.vizContent,
	)
	ui.vizAccordion.Open(0)

	ui.topBarLine1 = container.NewBorder(nil, nil,
		widget.NewLabel("Preset:"),
		container.NewHBox(
			ui.savePresetBtn,
			ui.deletePresetBtn,
			ui.loadPresetsBtn,
		),
		ui.presetCombo,
	)
	ui.topBarLine2 = container.NewBorder(nil, nil,
		nil,
		ui.resetBtn,
		container.NewHBox(
			ui.playBtn,
			ui.playSeqBtn,
			ui.saveBtn,
		),
	)
	ui.topBar = container.NewVBox(
		ui.topBarLine1,
		ui.topBarLine2,
	)

	// Vertical stack without VScroll lets the window size directly to all widgets with zero scrollbars.
	root := container.NewVBox(
		ui.topBar,
		container.NewGridWithColumns(3, ui.colComposer, ui.colUnison, ui.colADSR),
		ui.seqAccordion,
		ui.metaAccordion,
		ui.vizAccordion,
		ui.statusLabel,
	)

	ui.window.SetContent(root)
}

func main() {
	ui := newSynthUI()
	ui.window.ShowAndRun()
}
