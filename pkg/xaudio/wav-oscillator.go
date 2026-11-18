package xaudio

import (
	"fmt"
	"math"
	"math/rand"
)

// Waveform defines the supported built-in oscillator shapes.
type Waveform int

const (
	// Sine generates a pure sine wave.
	Sine Waveform = iota
	// Sawtooth generates a naive (non-band-limited) rising sawtooth wave.
	Sawtooth
	// SawtoothPoly generates a PolyBLEP band-limited rising sawtooth wave.
	SawtoothPoly
	// Square generates a naive (non-band-limited) pulse/square wave.
	Square
	// SquarePoly generates a PolyBLEP band-limited pulse/square wave.
	SquarePoly
	// Triangle generates a continuous bipolar triangle wave.
	Triangle
	// FM generates a 2-operator frequency modulation oscillator (carrier + modulator).
	FM
)

// OscillatorFunc computes a bipolar sample [-1.0, +1.0] for normalized phase [0.0, 1.0)
// and phase increment dt (frequency / sampleRate).
type OscillatorFunc func(phase, dt float64) float64

// FMOperatorFunc computes a bipolar FM sample [-1.0, +1.0] using independent carrier and
// modulator phases in [0.0, 1.0) along with the carrier phase increment dt.
type FMOperatorFunc func(carrierPhase, modPhase, dt float64) float64

// DetuneCurve maps a normalized voice position pos in [-1.0, +1.0]
// to a detune offset multiplier.
type DetuneCurve func(pos float64) float64

// Built-in DetuneCurve callbacks.
var (
	// CurveLinear distributes voices evenly across the detune range.
	CurveLinear DetuneCurve = func(pos float64) float64 {
		return pos
	}

	// CurveJP8000 uses Roland's non-linear cubic spacing.
	CurveJP8000 DetuneCurve = func(pos float64) float64 {
		return 0.75*(pos*pos*pos) + 0.25*pos
	}

	// CurveExponential pushes voices further toward the outer edges.
	CurveExponential DetuneCurve = func(pos float64) float64 {
		if pos < 0 {
			return -math.Sqrt(-pos)
		}
		return math.Sqrt(pos)
	}
)

// UnisonConfig configures the multi-voice unison and supersaw synthesis engine.
type UnisonConfig struct {
	// BaseWave selects the built-in waveform shape (used when EvalOscillator and EvalFM are nil).
	BaseWave Waveform
	// EvalOscillator allows supplying a custom single-phase oscillator callback; overrides BaseWave if non-nil.
	EvalOscillator OscillatorFunc
	// EvalFM allows supplying a custom two-phase FM oscillator callback; overrides BaseWave and EvalOscillator if non-nil.
	EvalFM FMOperatorFunc
	// Voices determines the number of stacked oscillator voices.
	Voices int
	// Detune controls pitch spread from 0.0 to 1.0.
	Detune float64
	// DetuneCurve determines the detune distribution across voices.
	DetuneCurve DetuneCurve
	// Mix controls the balance between side voices and center voice (0.0 to 1.0).
	Mix float64
	// StereoWidth sets stereo spread from 0.0 (Mono) to 1.0 (full L/R pan).
	StereoWidth float64
	// SubOctave drops the outermost voice pair down by -1 octave.
	SubOctave bool
	// AnalogDrift adds slow, independent LFO pitch wander per voice (0.0 to 1.0).
	AnalogDrift float64
	// HardSync resets all carrier and modulator phases whenever the master oscillator completes a cycle.
	HardSync bool
	// PulseWidth controls duty cycle for Square and SquarePoly waveforms (0.05 to 0.95).
	PulseWidth float64
	// ModRatio sets the modulator-to-carrier frequency ratio when BaseWave == FM.
	ModRatio float64
	// ModIndex sets the modulation depth (phase deviation) when BaseWave == FM.
	ModIndex float64
}

// Preset configurations for standard supersaw and unison variants.
var (
	// PresetJP8000 models the classic 7-voice mono Roland JP-8000 supersaw.
	PresetJP8000 = UnisonConfig{
		BaseWave: SawtoothPoly, Voices: 7, Detune: 0.35, DetuneCurve: CurveJP8000,
		Mix: 0.75, StereoWidth: 0.0,
	}

	// PresetHypersaw models a 9-voice wide stereo supersaw with sub-octave reinforcement.
	PresetHypersaw = UnisonConfig{
		BaseWave: SawtoothPoly, Voices: 9, Detune: 0.45, DetuneCurve: CurveJP8000,
		Mix: 0.85, StereoWidth: 0.9, SubOctave: true,
	}

	// PresetMultiSaw creates a smoothly detuned 8-voice stereo sawtooth stack.
	PresetMultiSaw = UnisonConfig{
		BaseWave: SawtoothPoly, Voices: 8, Detune: 0.30, DetuneCurve: CurveLinear,
		Mix: 0.80, StereoWidth: 0.8,
	}

	// PresetSupersquare creates a detuned 7-voice pulse wave stack.
	PresetSupersquare = UnisonConfig{
		BaseWave: SquarePoly, Voices: 7, Detune: 0.40, DetuneCurve: CurveJP8000,
		Mix: 0.80, StereoWidth: 0.85, PulseWidth: 0.35,
	}

	// PresetSupertriangle creates a warm 7-voice triangle stack with sub-octave.
	PresetSupertriangle = UnisonConfig{
		BaseWave: Triangle, Voices: 7, Detune: 0.25, DetuneCurve: CurveLinear,
		Mix: 0.85, StereoWidth: 0.75, SubOctave: true,
	}

	// PresetSupersine creates a 5-voice detuned sine unison.
	PresetSupersine = UnisonConfig{
		BaseWave: Sine, Voices: 5, Detune: 0.20, DetuneCurve: CurveLinear,
		Mix: 0.80, StereoWidth: 0.60,
	}

	// PresetAnalogDrift simulates vintage thermal pitch drift across 6 voices.
	PresetAnalogDrift = UnisonConfig{
		BaseWave: SawtoothPoly, Voices: 6, Detune: 0.25, DetuneCurve: CurveLinear,
		Mix: 0.85, StereoWidth: 0.70, AnalogDrift: 0.65,
	}

	// PresetHardSync creates an aggressive 7-voice hard-synced lead.
	PresetHardSync = UnisonConfig{
		BaseWave: SawtoothPoly, Voices: 7, Detune: 0.85, DetuneCurve: CurveJP8000,
		Mix: 0.90, StereoWidth: 0.75, HardSync: true,
	}

	// PresetFM creates a 5-voice bell-like FM unison patch with non-integer ratio 3.5.
	PresetFM = UnisonConfig{
		BaseWave: FM, Voices: 5, Detune: 0.20, DetuneCurve: CurveLinear,
		Mix: 0.80, StereoWidth: 0.80, ModRatio: 3.5, ModIndex: 1.8,
	}
)

// GenerateWave produces mono 16-bit PCM samples for the chosen waveform.
// gain should be between 0.0 and 1.0 (e.g., 0.8 to leave headroom).
func GenerateWave(wave Waveform, sampleRate, durationSec int, freq, gain float64) ([]int16, error) {
	if gain < 0.0 || gain > 1.0 {
		return nil, fmt.Errorf("gain must be between 0.0 and 1.0, got %f", gain)
	}

	numSamples := sampleRate * durationSec
	samples := make([]int16, numSamples)
	amplitude := gain * math.MaxInt16

	phaseIncrement := freq / float64(sampleRate)
	phase := 0.0

	for i := 0; i < numSamples; i++ {
		var sample float64

		switch wave {
		case Sine:
			sample = math.Sin(2.0 * math.Pi * phase)
		case SawtoothPoly:
			fallthrough
		case Sawtooth:
			sample = 2.0*phase - 1.0
		case SquarePoly:
			fallthrough
		case Square:
			if phase < 0.5 {
				sample = 1.0
			} else {
				sample = -1.0
			}
		case Triangle:
			sample = 4.0*math.Abs(phase-0.5) - 1.0
		default:
			return nil, fmt.Errorf("unsupported waveform: %v", wave)
		}

		samples[i] = int16(amplitude * sample)

		phase += phaseIncrement
		if phase >= 1.0 {
			phase -= 1.0
		}
	}

	return samples, nil
}

// GenerateSupersaw renders a 7-voice detuned sawtooth oscillator.
//   - detune (0.0 to 1.0): pitch spread across the 6 side voices (~0 to 50 cents)
//   - mix (0.0 to 1.0): volume balance between the 6 detuned voices and the center voice
func GenerateSupersaw(sampleRate, durationSec int, freq, gain, detune, mix float64) []int16 {
	// JP-8000 style asymmetric detune ratios around the fundamental (0.0)
	detuneOffsets := []float64{
		-0.11002313,
		-0.06288439,
		-0.01952356,
		0.0, // Center fundamental voice
		0.01991221,
		0.06216538,
		0.10745242,
	}

	return GenerateSupersawExt(sampleRate, durationSec, freq, gain, detune, mix, detuneOffsets)
}

func GenerateSupersawExt(sampleRate, durationSec int, freq, gain, detune, mix float64, detuneOffsets []float64) []int16 {
	var phases [7]float64
	var increments [7]float64
	var weights [7]float64

	dLen := len(detuneOffsets)

	// Max detune spread of ~3% frequency (~50 cents)
	detuneScale := detune * 0.03
	sideWeight := mix
	centerWeight := 1.0 - (mix * 0.5)
	totalWeight := centerWeight + float64(dLen-1)*sideWeight

	rng := rand.New(rand.NewSource(42)) // Deterministic random start phases to prevent initial click/spike
	for v := 0; v < dLen; v++ {
		voiceFreq := freq * (1.0 + detuneOffsets[v]*detuneScale)
		increments[v] = voiceFreq / float64(sampleRate)
		phases[v] = rng.Float64()

		if v == int(dLen/2) {
			weights[v] = centerWeight / totalWeight
		} else {
			weights[v] = sideWeight / totalWeight
		}
	}

	numSamples := sampleRate * durationSec
	samples := make([]int16, numSamples)
	amplitude := gain * math.MaxInt16

	for i := 0; i < numSamples; i++ {
		var sum float64
		for v := 0; v < 7; v++ {
			saw := 2.0*phases[v] - 1.0
			sum += saw * weights[v]

			phases[v] += increments[v]
			if phases[v] >= 1.0 {
				phases[v] -= 1.0
			}
		}

		// Clamp to [-1.0, 1.0] for safety before scaling to int16
		if sum > 1.0 {
			sum = 1.0
		} else if sum < -1.0 {
			sum = -1.0
		}

		samples[i] = int16(amplitude * sum)
	}

	return samples
}

// WriteWaveWav is a convenience helper that generates a waveform and writes it to disk.
func WriteWaveWav(filename string, wave Waveform, durationSec int, freq float64) error {
	const sampleRate = 48000
	samples, err := GenerateWave(wave, sampleRate, durationSec, freq, 0.9)
	if err != nil {
		return err
	}
	return WriteWav(filename, sampleRate, 1, samples)
}

// PolyBLEP computes the polynomial band-limited step correction at phase
// discontinuities (t = 0 and t = 1).
func PolyBLEP(t, dt float64) float64 {
	if dt <= 0.0 {
		return 0.0
	}
	if t < dt {
		t /= dt
		return t + t - t*t - 1.0
	} else if t > 1.0-dt {
		t = (t - 1.0) / dt
		return t*t + t + t + 1.0
	}
	return 0.0
}

// defaultOscillator returns an OscillatorFunc for the built-in Waveform shapes.
func defaultOscillator(wave Waveform, pulseWidth float64) OscillatorFunc {
	switch wave {
	case Sine:
		return func(phase, dt float64) float64 {
			return math.Sin(2.0 * math.Pi * phase)
		}
	case Sawtooth:
		return func(phase, dt float64) float64 {
			val := 2.0*phase - 1.0
			return val
		}
	case SawtoothPoly:
		return func(phase, dt float64) float64 {
			val := 2.0*phase - 1.0
			val -= PolyBLEP(phase, dt)
			return val
		}
	case Square:
		return func(phase, dt float64) float64 {
			if phase < pulseWidth {
				return 1.0
			}
			return -1.0
		}
	case SquarePoly:
		return func(phase, dt float64) float64 {
			var val float64
			if phase < pulseWidth {
				val = 1.0
			} else {
				val = -1.0
			}
			val += PolyBLEP(phase, dt)
			val -= PolyBLEP(math.Mod(phase+1.0-pulseWidth, 1.0), dt)
			return val
		}
	case Triangle:
		return func(phase, dt float64) float64 {
			return 4.0*math.Abs(phase-0.5) - 1.0
		}
	default:
		return func(phase, dt float64) float64 {
			return 0.0
		}
	}
}

// defaultFMOscillator returns a phase-continuous 2-operator FM callback using independent
// carrierPhase and modPhase accumulators.
func defaultFMOscillator(modIndex float64) FMOperatorFunc {
	return func(carrierPhase, modPhase, dt float64) float64 {
		modulator := math.Sin(2.0 * math.Pi * modPhase)
		return math.Sin(2.0*math.Pi*carrierPhase + modIndex*modulator)
	}
}

// generateUnisonFloat renders interleaved stereo samples in [-1.0, +1.0] for a single note.
func generateUnisonFloat(cfg UnisonConfig, sampleRate int, durationSec, freq, gain float64, seed int64) ([]float64, error) {
	if cfg.Voices < 1 {
		return nil, fmt.Errorf("voices must be >= 1, got %d", cfg.Voices)
	}
	if gain < 0.0 || gain > 1.0 {
		return nil, fmt.Errorf("gain must be between 0.0 and 1.0, got %f", gain)
	}
	if cfg.PulseWidth <= 0.0 || cfg.PulseWidth >= 1.0 {
		cfg.PulseWidth = 0.5
	}
	if cfg.ModRatio == 0.0 {
		cfg.ModRatio = 2.0
	}

	// Resolve FM vs standard single-phase oscillator callback once before the sample loop.
	evalFM := cfg.EvalFM
	if evalFM == nil && cfg.EvalOscillator == nil && cfg.BaseWave == FM {
		evalFM = defaultFMOscillator(cfg.ModIndex)
	}

	evalOsc := cfg.EvalOscillator
	if evalFM == nil && evalOsc == nil {
		evalOsc = defaultOscillator(cfg.BaseWave, cfg.PulseWidth)
	}

	curve := cfg.DetuneCurve
	if curve == nil {
		curve = CurveLinear
	}

	rng := rand.New(rand.NewSource(seed))

	phases := make([]float64, cfg.Voices)
	modPhases := make([]float64, cfg.Voices)
	baseIncrements := make([]float64, cfg.Voices)
	weights := make([]float64, cfg.Voices)
	panL := make([]float64, cfg.Voices)
	panR := make([]float64, cfg.Voices)
	driftPhases := make([]float64, cfg.Voices)
	driftRates := make([]float64, cfg.Voices)

	detuneScale := cfg.Detune * 0.035
	centerIdx := cfg.Voices / 2
	hasCenter := cfg.Voices%2 == 1

	sideWeight := cfg.Mix
	centerWeight := 1.0 - (cfg.Mix * 0.5)
	var totalWeight float64

	for v := 0; v < cfg.Voices; v++ {
		var pos float64
		if cfg.Voices > 1 {
			pos = (2.0*float64(v)/float64(cfg.Voices-1) - 1.0)
		}

		offset := curve(pos)
		voiceFreq := freq * (1.0 + offset*detuneScale)

		if cfg.SubOctave && cfg.Voices >= 3 && (v == 0 || v == cfg.Voices-1) {
			voiceFreq *= 0.5
		}

		baseIncrements[v] = voiceFreq / float64(sampleRate)
		initPhase := rng.Float64()
		phases[v] = initPhase
		// Initialize modulator phase coherently with initial carrier phase.
		_, modFrac := math.Modf(initPhase * cfg.ModRatio)
		if modFrac < 0.0 {
			modFrac += 1.0
		}
		modPhases[v] = modFrac

		if hasCenter && v == centerIdx {
			weights[v] = centerWeight
		} else {
			weights[v] = sideWeight
		}
		totalWeight += weights[v]

		pan := pos * cfg.StereoWidth
		angle := (pan + 1.0) * (math.Pi / 4.0)
		panL[v] = math.Cos(angle)
		panR[v] = math.Sin(angle)

		driftPhases[v] = rng.Float64() * 2.0 * math.Pi
		driftRates[v] = (0.15 + rng.Float64()*0.30) * 2.0 * math.Pi / float64(sampleRate)
	}

	for v := 0; v < cfg.Voices; v++ {
		weights[v] /= totalWeight
	}

	numFrames := int(math.Round(float64(sampleRate) * durationSec))
	samples := make([]float64, numFrames*2)

	masterPhase := 0.0
	masterInc := freq / float64(sampleRate)

	for i := 0; i < numFrames; i++ {
		syncReset := false
		if cfg.HardSync {
			masterPhase += masterInc
			if masterPhase >= 1.0 {
				masterPhase -= 1.0
				syncReset = true
			}
		}

		var sumL, sumR float64

		for v := 0; v < cfg.Voices; v++ {
			if syncReset {
				phases[v] = 0.0
				modPhases[v] = 0.0
			}

			inc := baseIncrements[v]
			if cfg.AnalogDrift > 0.0 {
				drift := math.Sin(driftPhases[v]) * cfg.AnalogDrift * 0.004
				inc *= (1.0 + drift)
				driftPhases[v] += driftRates[v]
			}

			var rawSample float64
			if evalFM != nil {
				rawSample = evalFM(phases[v], modPhases[v], inc)
			} else {
				rawSample = evalOsc(phases[v], inc)
			}

			val := rawSample * weights[v]
			sumL += val * panL[v]
			sumR += val * panR[v]

			// Advance carrier phase accumulator and wrap cleanly into [0.0, 1.0).
			phases[v] += inc
			if phases[v] >= 1.0 {
				phases[v] -= math.Floor(phases[v])
			}

			// Advance independent modulator phase accumulator and wrap cleanly into [0.0, 1.0).
			modPhases[v] += inc * cfg.ModRatio
			if modPhases[v] >= 1.0 || modPhases[v] < 0.0 {
				modPhases[v] -= math.Floor(modPhases[v])
			}
		}

		samples[i*2] = sumL * gain
		samples[i*2+1] = sumR * gain
	}

	return samples, nil
}

// GenerateUnison renders interleaved 16-bit stereo PCM samples for any unison variant.
func GenerateUnison(cfg UnisonConfig, sampleRate, durationSec int, freq, gain float64) ([]int16, error) {
	raw, err := generateUnisonFloat(cfg, sampleRate, float64(durationSec), freq, gain, 42)
	if err != nil {
		return nil, err
	}
	out := make([]int16, len(raw))
	for i, v := range raw {
		out[i] = clampToInt16(v * math.MaxInt16)
	}
	return out, nil
}

func clampToInt16(v float64) int16 {
	if v > math.MaxInt16 {
		return math.MaxInt16
	}
	if v < math.MinInt16 {
		return math.MinInt16
	}
	return int16(v)
}

// WriteUnisonWav generates any supersaw variant and writes it as a 48kHz stereo WAV file.
func WriteUnisonWav(filename string, cfg UnisonConfig, durationSec int, freq float64) error {
	const sampleRate = 48000
	samples, err := GenerateUnison(cfg, sampleRate, durationSec, freq, 0.5)
	if err != nil {
		return err
	}
	return WriteWav(filename, sampleRate, 2, samples)
}

// WriteUnisonWavWithMetadata generates any supersaw variant and writes it as a 48kHz stereo WAV file with RIFF INFO metadata.
func WriteUnisonWavWithMetadata(filename string, cfg UnisonConfig, durationSec int, freq float64, meta WavMetadata) error {
	const sampleRate = 48000
	samples, err := GenerateUnison(cfg, sampleRate, durationSec, freq, 0.5)
	if err != nil {
		return err
	}
	return WriteWavWithMetadata(filename, sampleRate, 2, samples, meta)
}
