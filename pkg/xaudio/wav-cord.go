package xaudio

import (
	"fmt"
	"math"
)

// ============================================================================
// Chord & Sequencer Helpers
// ============================================================================

// EnvelopeCurve maps normalized stage progress [0.0, 1.0] to an interpolation factor [0.0, 1.0].
type EnvelopeCurve func(progress float64) float64

// NewBicubicWobbleCurve constructs a parametric EnvelopeCurve combining a cubic Hermite
// base transition (3t^2 - 2t^3) with a boundary-damped sinusoidal wobble controlled by
// the number of oscillation cycles and modulation depth.
func NewBicubicWobbleCurve(cycles, depth float64) EnvelopeCurve {
	return func(progress float64) float64 {
		// Cubic Hermite smoothstep base curve guaranteeing 0.0 at t=0 and 1.0 at t=1.
		base := progress * progress * (3.0 - 2.0*progress)
		// Boundary-damped wobble envelope: sin(pi * t) vanishes at t=0 and t=1.
		window := math.Sin(math.Pi * progress)
		wobble := depth * window * math.Sin(2.0*math.Pi*cycles*progress)
		val := base + wobble
		if val < 0.0 {
			return 0.0
		}
		if val > 1.0 {
			return 1.0
		}
		return val
	}
}

// Built-in EnvelopeCurve interpolation callbacks.
var (
	// EnvLinear interpolates linearly between envelope stage endpoints.
	EnvLinear EnvelopeCurve = func(progress float64) float64 {
		return progress
	}

	// EnvExponential creates a concave curve that changes slowly at first and accelerates near the end.
	EnvExponential EnvelopeCurve = func(progress float64) float64 {
		return progress * progress
	}

	// EnvLogarithmic creates a convex curve that changes rapidly at first and tapers off smoothly.
	EnvLogarithmic EnvelopeCurve = func(progress float64) float64 {
		return 1.0 - (1.0-progress)*(1.0-progress)
	}

	// EnvHermite uses cubic Hermite smoothstep interpolation (3t^2 - 2t^3) with zero derivatives at endpoints.
	EnvHermite EnvelopeCurve = func(progress float64) float64 {
		return progress * progress * (3.0 - 2.0*progress)
	}

	// EnvRadian uses half-cycle cosine radian interpolation (0.5 * (1 - cos(pi * t))) across [0, pi].
	EnvRadian EnvelopeCurve = func(progress float64) float64 {
		return 0.5 * (1.0 - math.Cos(math.Pi*progress))
	}

	// EnvBicubicWobble provides a default parametric bicubic wobble curve (3 cycles, 0.15 depth).
	EnvBicubicWobble EnvelopeCurve = NewBicubicWobbleCurve(3.0, 0.15)
)

// ADSR defines a time-varying amplitude envelope in seconds with pluggable stage interpolation curves.
type ADSR struct {
	Attack       float64       // Rise time after note-on
	Decay        float64       // Fall time from peak to Sustain level
	Sustain      float64       // Held amplitude multiplier (0.0 to 1.0)
	Release      float64       // Fade-out tail after note-off
	AttackCurve  EnvelopeCurve // Interpolation callback for Attack stage (defaults to EnvLinear if nil)
	DecayCurve   EnvelopeCurve // Interpolation callback for Decay stage (defaults to EnvLinear if nil)
	ReleaseCurve EnvelopeCurve // Interpolation callback for Release stage (defaults to EnvLinear if nil)
}

// DefaultPadADSR provides a smooth envelope for supersaw chords and leads.
var DefaultPadADSR = ADSR{
	Attack:       0.03,
	Decay:        0.15,
	Sustain:      0.80,
	Release:      0.25,
	AttackCurve:  EnvLinear,
	DecayCurve:   EnvExponential,
	ReleaseCurve: EnvExponential,
}

// NoteEvent represents a single note event on a timeline.
type NoteEvent struct {
	StartSec    float64       // Start time in seconds
	DurationSec float64       // Held duration in seconds
	Freq        float64       // Pitch frequency in Hz
	Velocity    float64       // Note velocity (0.0 to 1.0)
	Config      *UnisonConfig // Optional per-note configuration override
	Envelope    *ADSR         // Optional per-note ADSR envelope override (uses sequence default if nil)
}

// NoteToFreq converts a MIDI note number (e.g., 69 = A4) to Hz.
func NoteToFreq(midiNote int) float64 {
	return 440.0 * math.Pow(2.0, float64(midiNote-69)/12.0)
}

// Predefined chord interval patterns in semitones.
var (
	ChordMajor    = []int{0, 4, 7}
	ChordMinor    = []int{0, 3, 7}
	ChordMaj7     = []int{0, 4, 7, 11}
	ChordMin7     = []int{0, 3, 7, 10}
	ChordSus4     = []int{0, 5, 7}
	ChordAdd9     = []int{0, 4, 7, 14}
	ChordPower    = []int{0, 7, 12}
	ChordSupersaw = []int{-12, 0, 7, 12, 16, 19}
)

// NewChordEvents generates simultaneous NoteEvents for a chord rooted at rootMIDI.
func NewChordEvents(rootMIDI int, intervals []int, startSec, durationSec, velocity float64) []NoteEvent {
	events := make([]NoteEvent, len(intervals))
	noteVel := velocity / math.Sqrt(float64(len(intervals)))
	for i, semitone := range intervals {
		events[i] = NoteEvent{
			StartSec:    startSec,
			DurationSec: durationSec,
			Freq:        NoteToFreq(rootMIDI + semitone),
			Velocity:    noteVel,
		}
	}
	return events
}

// evalStageCurve evaluates a stage curve callback safely, defaulting to linear interpolation if nil.
func evalStageCurve(curve EnvelopeCurve, progress float64) float64 {
	if progress <= 0.0 {
		return 0.0
	}
	if progress >= 1.0 {
		return 1.0
	}
	if curve == nil {
		return progress
	}
	return curve(progress)
}

// envelopeGain computes the ADSR gain multiplier at elapsed time t for hold duration holdSec.
func (env ADSR) envelopeGain(t, holdSec float64) float64 {
	if t < 0 {
		return 0.0
	}

	// Determine the envelope amplitude at the exact moment the note is released.
	levelAtRelease := env.Sustain
	if holdSec < env.Attack && env.Attack > 0 {
		attFactor := evalStageCurve(env.AttackCurve, holdSec/env.Attack)
		levelAtRelease = attFactor
	} else if holdSec < env.Attack+env.Decay && env.Decay > 0 {
		decFactor := evalStageCurve(env.DecayCurve, (holdSec-env.Attack)/env.Decay)
		levelAtRelease = 1.0 - decFactor*(1.0-env.Sustain)
	}

	// Release stage (t >= holdSec).
	if t >= holdSec {
		if env.Release <= 0 {
			return 0.0
		}
		relProgress := (t - holdSec) / env.Release
		if relProgress >= 1.0 {
			return 0.0
		}
		relFactor := evalStageCurve(env.ReleaseCurve, relProgress)
		return levelAtRelease * (1.0 - relFactor)
	}

	// Attack stage (0 <= t < Attack).
	if t < env.Attack && env.Attack > 0 {
		return evalStageCurve(env.AttackCurve, t/env.Attack)
	}

	// Decay stage (Attack <= t < Attack + Decay).
	if t < env.Attack+env.Decay && env.Decay > 0 {
		decFactor := evalStageCurve(env.DecayCurve, (t-env.Attack)/env.Decay)
		return 1.0 - decFactor*(1.0-env.Sustain)
	}

	// Sustain stage.
	return env.Sustain
}

// RenderSequence mixes polyphonic NoteEvents into a 16-bit stereo buffer with peak normalization.
func RenderSequence(defaultCfg UnisonConfig, defaultEnv ADSR, sampleRate int, events []NoteEvent) ([]int16, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("events list is empty")
	}

	// Determine total timeline duration using each note's resolved ADSR release tail.
	var totalSec float64
	for _, ev := range events {
		env := defaultEnv
		if ev.Envelope != nil {
			env = *ev.Envelope
		}
		end := ev.StartSec + ev.DurationSec + env.Release
		if end > totalSec {
			totalSec = end
		}
	}

	totalFrames := int(math.Ceil(totalSec * float64(sampleRate)))
	mixBuffer := make([]float64, totalFrames*2)

	for idx, ev := range events {
		cfg := defaultCfg
		if ev.Config != nil {
			cfg = *ev.Config
		}

		env := defaultEnv
		if ev.Envelope != nil {
			env = *ev.Envelope
		}

		noteTotalSec := ev.DurationSec + env.Release
		noteSamples, err := generateUnisonFloat(cfg, sampleRate, noteTotalSec, ev.Freq, ev.Velocity, int64(100+idx))
		if err != nil {
			return nil, fmt.Errorf("event %d: %w", idx, err)
		}

		startFrame := int(math.Round(ev.StartSec * float64(sampleRate)))
		noteFrames := len(noteSamples) / 2

		for f := 0; f < noteFrames; f++ {
			dstFrame := startFrame + f
			if dstFrame < 0 || dstFrame >= totalFrames {
				continue
			}

			t := float64(f) / float64(sampleRate)
			gain := env.envelopeGain(t, ev.DurationSec)

			mixBuffer[dstFrame*2] += noteSamples[f*2] * gain
			mixBuffer[dstFrame*2+1] += noteSamples[f*2+1] * gain
		}
	}

	var peak float64
	for _, s := range mixBuffer {
		abs := math.Abs(s)
		if abs > peak {
			peak = abs
		}
	}

	masterScale := 0.95 * math.MaxInt16
	if peak > 1.0 {
		masterScale /= peak
	}

	out := make([]int16, len(mixBuffer))
	for i, s := range mixBuffer {
		out[i] = clampToInt16(s * masterScale)
	}

	return out, nil
}

// WriteSequenceWav renders a polyphonic sequence and writes it to a 48 kHz stereo WAV file.
func WriteSequenceWav(filename string, defaultCfg UnisonConfig, defaultEnv ADSR, events []NoteEvent) error {
	const sampleRate = 48000
	samples, err := RenderSequence(defaultCfg, defaultEnv, sampleRate, events)
	if err != nil {
		return err
	}
	return WriteWav(filename, sampleRate, 2, samples)
}

// RenderNote renders a single NoteEvent (including its ADSR envelope and release tail)
// into a 16-bit stereo PCM sample buffer.
func RenderNote(defaultCfg UnisonConfig, defaultEnv ADSR, sampleRate int, note NoteEvent) ([]int16, error) {
	if note.DurationSec < 0 {
		return nil, fmt.Errorf("note duration must be >= 0, got %f", note.DurationSec)
	}
	// Normalize StartSec to 0.0 so a single note renders immediately without leading silence.
	single := note
	single.StartSec = 0.0
	return RenderSequence(defaultCfg, defaultEnv, sampleRate, []NoteEvent{single})
}

// WriteNoteWav renders a single NoteEvent and writes it as a 48 kHz stereo WAV file.
func WriteNoteWav(filename string, defaultCfg UnisonConfig, defaultEnv ADSR, note NoteEvent) error {
	return WriteNoteWavWithMetadata(filename, defaultCfg, defaultEnv, note, WavMetadata{})
}

// WriteNoteWavWithMetadata renders a single NoteEvent and writes it as a 48 kHz stereo WAV file
// with embedded RIFF INFO metadata.
func WriteNoteWavWithMetadata(filename string, defaultCfg UnisonConfig, defaultEnv ADSR, note NoteEvent, meta WavMetadata) error {
	const sampleRate = 48000
	samples, err := RenderNote(defaultCfg, defaultEnv, sampleRate, note)
	if err != nil {
		return err
	}
	return WriteWavWithMetadata(filename, sampleRate, 2, samples, meta)
}

func DumpTranceAnthem(filename string) error {
	var song []NoteEvent

	// 1. Build a 4-bar supersaw chord progression (Am -> F -> C -> G)
	// Using wide trance voicings (root bass + open triad)
	song = append(song, NewChordEvents(57, ChordMinor, 0.0, 1.8, 0.85)...) // A3 minor
	song = append(song, NewChordEvents(53, ChordMajor, 2.0, 1.8, 0.85)...) // F3 major
	song = append(song, NewChordEvents(60, ChordMajor, 4.0, 1.8, 0.85)...) // C4 major
	song = append(song, NewChordEvents(55, ChordMajor, 6.0, 1.8, 0.85)...) // G3 major

	// 2. Layer a higher-pitched Hypersaw lead melody over the chords
	leadCfg := PresetHypersaw
	melody := []struct {
		midi  int
		start float64
		dur   float64
	}{
		{72, 0.0, 0.45}, // C5
		{76, 0.5, 0.45}, // E5
		{79, 1.0, 0.90}, // G5
		{77, 2.0, 0.90}, // F5
		{76, 3.0, 0.90}, // E5
		{72, 4.0, 1.80}, // C5
		{74, 6.0, 1.80}, // D5
	}

	for _, m := range melody {
		song = append(song, NoteEvent{
			StartSec:    m.start,
			DurationSec: m.dur,
			Freq:        NoteToFreq(m.midi),
			Velocity:    0.75,
			Config:      &leadCfg, // Override default patch for the lead line
		})
	}

	// 4. Render everything into a single 48kHz stereo WAV file
	meta := WavMetadata{
		Title:       "Trance Anthem",
		Artist:      "Your Name",
		Copyright:   "Copyright 2026 Your Name. All rights reserved.",
		Description: "Rendered with 9-voice PolyBLEP Hypersaw and MultiSaw chords at 48kHz.",
		Software:    "xaudio Go Synth",
		Date:        "2026-11-16",
	}

	// 4. Render everything into a single 48kHz stereo WAV file
	samples, _ := RenderSequence(PresetMultiSaw, DefaultPadADSR, 48000, song)
	return WriteWavWithMetadata(filename, 48000, 2, samples, meta)
}
