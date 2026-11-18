package xaudio

import (
	"encoding/binary"
	"os"
)

// WavMetadata holds standard RIFF INFO metadata fields.
type WavMetadata struct {
	Title       string // INAM
	Artist      string // IART
	Copyright   string // ICOP
	Description string // ICMT (Comments / Notes)
	Software    string // ISFT
	Date        string // ICRD
}

// buildInfoChunk serializes non-empty metadata fields into a valid RIFF "LIST" ("INFO") chunk.
// Returns nil if no metadata fields are populated.
func buildInfoChunk(meta WavMetadata) []byte {
	entries := []struct {
		id  [4]byte
		val string
	}{
		{[4]byte{'I', 'N', 'A', 'M'}, meta.Title},
		{[4]byte{'I', 'A', 'R', 'T'}, meta.Artist},
		{[4]byte{'I', 'C', 'O', 'P'}, meta.Copyright},
		{[4]byte{'I', 'C', 'M', 'T'}, meta.Description},
		{[4]byte{'I', 'S', 'F', 'T'}, meta.Software},
		{[4]byte{'I', 'C', 'R', 'D'}, meta.Date},
	}

	var subchunks []byte
	for _, e := range entries {
		if e.val == "" {
			continue
		}

		// RIFF INFO strings are null-terminated
		strBytes := append([]byte(e.val), 0x00)
		subSize := uint32(len(strBytes))

		// Sub-chunk header: 4-byte ID + 4-byte little-endian size
		var hdr [8]byte
		copy(hdr[0:4], e.id[:])
		binary.LittleEndian.PutUint32(hdr[4:8], subSize)

		subchunks = append(subchunks, hdr[:]...)
		subchunks = append(subchunks, strBytes...)

		// RIFF specification requires every sub-chunk to be word-aligned (even byte count).
		// Note: padding byte is NOT included in subSize, but IS included in the outer LIST size.
		if len(strBytes)%2 != 0 {
			subchunks = append(subchunks, 0x00)
		}
	}

	if len(subchunks) == 0 {
		return nil
	}

	// Outer LIST chunk: "LIST" (4) + Size (4) + "INFO" (4) + subchunks
	listSize := uint32(4 + len(subchunks))
	out := make([]byte, 12+len(subchunks))
	copy(out[0:4], "LIST")
	binary.LittleEndian.PutUint32(out[4:8], listSize)
	copy(out[8:12], "INFO")
	copy(out[12:], subchunks)

	return out
}

// WriteWavWithMetadata writes 16-bit PCM samples along with RIFF INFO metadata tags.
func WriteWavWithMetadata(filename string, sampleRate, numChannels int, samples []int16, meta WavMetadata) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	infoChunk := buildInfoChunk(meta)

	const bitsPerSample = 16
	bytesPerSample := bitsPerSample / 8
	dataSize := uint32(len(samples) * bytesPerSample)

	header := WavHeader{
		ChunkID:       [4]byte{'R', 'I', 'F', 'F'},
		ChunkSize:     36 + dataSize + uint32(len(infoChunk)), // Must include LIST/INFO bytes
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

	// 1. Write 44-byte RIFF/fmt/data header
	if err := binary.Write(f, binary.LittleEndian, &header); err != nil {
		return err
	}

	// 2. Write raw PCM audio samples
	if err := binary.Write(f, binary.LittleEndian, samples); err != nil {
		return err
	}

	// 3. Append LIST/INFO metadata chunk at the end
	if len(infoChunk) > 0 {
		if _, err := f.Write(infoChunk); err != nil {
			return err
		}
	}

	return nil
}
