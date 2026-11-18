package xaudio

// Description: Complete package for generating and storing 16-bit WAV audio files,
// including naive and PolyBLEP oscillators, phase-continuous FM synthesis, a multi-voice
// unison engine, callback-interpolated ADSR envelopes (with per-note overrides),
// polyphonic sequencer, and RIFF INFO metadata.

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

// WavHeader represents the standard 44-byte RIFF/WAV PCM header.
type WavHeader struct {
	ChunkID       [4]byte // "RIFF"
	ChunkSize     uint32  // 36 + Subchunk2Size
	Format        [4]byte // "WAVE"
	Subchunk1ID   [4]byte // "fmt "
	Subchunk1Size uint32  // 16 for PCM
	AudioFormat   uint16  // 1 for PCM
	NumChannels   uint16  // 1 = Mono, 2 = Stereo
	SampleRate    uint32  // e.g., 48000
	ByteRate      uint32  // SampleRate * NumChannels * BitsPerSample / 8
	BlockAlign    uint16  // NumChannels * BitsPerSample / 8
	BitsPerSample uint16  // 16
	Subchunk2ID   [4]byte // "data"
	Subchunk2Size uint32  // NumSamples * NumChannels * BitsPerSample / 8
}

// WriteWavHeader populates and writes the 44-byte WAV header to an io.Writer.
func WriteWavHeader(w io.Writer, sampleRate, numChannels, numSamples int) error {
	const bitsPerSample = 16
	bytesPerSample := bitsPerSample / 8
	dataSize := uint32(numSamples * bytesPerSample)

	header := WavHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + dataSize,
		Format:        [4]byte{'W', 'A', 'V', 'E'},
		Subchunk1ID:   [4]byte{'f', 'm', 't', ' '},
		Subchunk1Size: 16,
		AudioFormat:   1,
		NumChannels:   uint16(numChannels),
		SampleRate:    uint32(sampleRate),
		ByteRate:      uint32(sampleRate * numChannels * bytesPerSample),
		BlockAlign:    uint16(numChannels * bytesPerSample),
		BitsPerSample: bitsPerSample,
		Subchunk2ID:   [4]byte{'d', 'a', 't', 'a'},
		Subchunk2Size: dataSize,
	}

	return binary.Write(w, binary.LittleEndian, &header)
}

// WriteWav writes 16-bit PCM samples to a WAV file.
func WriteWav(filename string, sampleRate, numChannels int, samples []int16) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := WriteWavHeader(f, sampleRate, numChannels, len(samples)); err != nil {
		return err
	}

	return binary.Write(f, binary.LittleEndian, samples)
}

func DumpVariants(dirname string) error {
	return DumpVariantsExt(dirname, 10, 110.0)
}

// DumpVariantsExt renders all configured unison and supersaw variants into dirname
// as 48 kHz 16-bit stereo WAV files with embedded RIFF INFO metadata.
func DumpVariantsExt(dirname string, durationSec int, freq float64) error {
	// 1. Fifth-harmonic shimmer: outer voices detune wider with a steep 5th-power curve.
	customFifthPower := PresetHypersaw
	customFifthPower.DetuneCurve = func(pos float64) float64 {
		return math.Pow(pos, 5)
	}

	// 2. Asymmetric curve: sharp voices detune 1.5x while flat voices detune 0.5x.
	asymmetric := PresetMultiSaw
	asymmetric.DetuneCurve = func(pos float64) float64 {
		if pos > 0 {
			return pos * 1.5
		}
		return pos * 0.5
	}

	// 3. Naive (non-PolyBLEP) JP-8000 sawtooth variant for aliasing comparison.
	naiveSupersaw := PresetJP8000
	naiveSupersaw.BaseWave = Sawtooth

	// 4. Naive (non-PolyBLEP) supersquare variant for aliasing comparison.
	naiveSupersquare := PresetSupersquare
	naiveSupersquare.BaseWave = Square

	const (
		copyright = "CC0 1.0 Public Domain Dedication (https://creativecommons.org/publicdomain/zero/1.0/)"
		software  = "xaudio Go Synth"
		date      = "2026-11-17"
	)

	variants := []struct {
		filename string
		cfg      UnisonConfig
		title    string
		desc     string
	}{
		// Classic PolyBLEP hardware flavors
		{
			filename: "qsupersaw.wav",
			cfg:      PresetJP8000,
			title:    "JP-8000 Supersaw (PolyBLEP)",
			desc:     fmt.Sprintf("7-voice mono asymmetric cubic PolyBLEP supersaw at %.2f Hz", freq),
		},
		{
			filename: "qhypersaw.wav",
			cfg:      PresetHypersaw,
			title:    "Virus TI Hypersaw (PolyBLEP)",
			desc:     fmt.Sprintf("9-voice wide stereo PolyBLEP supersaw with -1 sub-octave at %.2f Hz", freq),
		},
		{
			filename: "qmultisaw.wav",
			cfg:      PresetMultiSaw,
			title:    "Multi-Saw Linear Unison (PolyBLEP)",
			desc:     fmt.Sprintf("8-voice smooth linearly detuned stereo PolyBLEP sawtooth stack at %.2f Hz", freq),
		},

		// Naive (non-band-limited) comparison variants
		{
			filename: "qsupersaw_naive.wav",
			cfg:      naiveSupersaw,
			title:    "JP-8000 Supersaw (Naive)",
			desc:     fmt.Sprintf("7-voice mono asymmetric cubic naive sawtooth stack at %.2f Hz", freq),
		},
		{
			filename: "qsupersquare_naive.wav",
			cfg:      naiveSupersquare,
			title:    "Supersquare 35% Pulse (Naive)",
			desc:     fmt.Sprintf("7-voice stereo naive pulse wave (0.35 duty cycle) at %.2f Hz", freq),
		},

		// Waveform spin-offs
		{
			filename: "qsupersquare.wav",
			cfg:      PresetSupersquare,
			title:    "Supersquare 35% Pulse (PolyBLEP)",
			desc:     fmt.Sprintf("7-voice stereo detuned PolyBLEP pulse wave (0.35 duty cycle) at %.2f Hz", freq),
		},
		{
			filename: "qsupertriangle.wav",
			cfg:      PresetSupertriangle,
			title:    "Supertriangle Sub Stack",
			desc:     fmt.Sprintf("7-voice stereo triangle unison with -1 sub-octave at %.2f Hz", freq),
		},
		{
			filename: "qsupersine.wav",
			cfg:      PresetSupersine,
			title:    "Supersine Organ Unison",
			desc:     fmt.Sprintf("5-voice warm detuned sine wave stack at %.2f Hz", freq),
		},
		{
			filename: "qfmsaw.wav",
			cfg:      PresetFM,
			title:    "FM Unison Bell Stack",
			desc:     fmt.Sprintf("5-voice stereo 2-operator FM unison (ratio 3.5, index 1.8) at %.2f Hz", freq),
		},

		// Structural / synthesis variants
		{
			filename: "qadrift.wav",
			cfg:      PresetAnalogDrift,
			title:    "Analog Thermal Drift Saw",
			desc:     fmt.Sprintf("6-voice stereo PolyBLEP sawtooth with independent LFO pitch drift at %.2f Hz", freq),
		},
		{
			filename: "qhardsync.wav",
			cfg:      PresetHardSync,
			title:    "Hard-Sync Supersaw Lead",
			desc:     fmt.Sprintf("7-voice hard-synced PolyBLEP sawtooth stack locked to master clock at %.2f Hz", freq),
		},

		// Custom detune curves
		{
			filename: "qhyper5saw.wav",
			cfg:      customFifthPower,
			title:    "Fifth-Power Hypersaw",
			desc:     fmt.Sprintf("9-voice stereo PolyBLEP hypersaw using x^5 outer-spread detune curve at %.2f Hz", freq),
		},
		{
			filename: "qasymultisaw.wav",
			cfg:      asymmetric,
			title:    "Asymmetric Multi-Saw",
			desc:     fmt.Sprintf("8-voice stereo PolyBLEP sawtooth with 1.5x sharp / 0.5x flat detune curve at %.2f Hz", freq),
		},
	}

	// Iterate over all variants and write each WAV file with its dedicated RIFF INFO metadata.
	for _, v := range variants {
		meta := WavMetadata{
			Title:       v.title,
			Copyright:   copyright,
			Description: v.desc,
			Software:    software,
			Date:        date,
		}
		outPath := filepath.Join(dirname, v.filename)
		if err := WriteUnisonWavWithMetadata(outPath, v.cfg, durationSec, freq, meta); err != nil {
			return fmt.Errorf("failed writing %s: %w", v.filename, err)
		}
	}

	return nil
}
