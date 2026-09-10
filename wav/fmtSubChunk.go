package wav

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ngyewch/go-pcm"
)

type fmtSubChunk struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16
	BitsPerSample uint16
}

func readFmtSubChunk(r io.Reader, size int, byteOrder binary.ByteOrder) (*fmtSubChunk, error) {
	subChunkBytes := make([]byte, size)
	_, err := io.ReadFull(r, subChunkBytes)
	if err != nil {
		return nil, err
	}
	helper := readerHelper{
		r:         bytes.NewReader(subChunkBytes),
		byteOrder: byteOrder,
	}
	audioFormat, err := helper.ReadUint16()
	if err != nil {
		return nil, err
	}
	numChannels, err := helper.ReadUint16()
	if err != nil {
		return nil, err
	}
	sampleRate, err := helper.ReadUint32()
	if err != nil {
		return nil, err
	}
	byteRate, err := helper.ReadUint32()
	if err != nil {
		return nil, err
	}
	blockAlign, err := helper.ReadUint16()
	if err != nil {
		return nil, err
	}
	bitsPerSample, err := helper.ReadUint16()
	if err != nil {
		return nil, err
	}
	return &fmtSubChunk{
		AudioFormat:   audioFormat,
		NumChannels:   numChannels,
		SampleRate:    sampleRate,
		ByteRate:      byteRate,
		BlockAlign:    blockAlign,
		BitsPerSample: bitsPerSample,
	}, nil
}

func (chunk fmtSubChunk) Encoding(byteOrder binary.ByteOrder) (pcm.Encoding, error) {
	switch chunk.AudioFormat {
	case 1:
		switch chunk.BitsPerSample {
		case 16:
			switch byteOrder {
			case binary.LittleEndian:
				return pcm.EncodingS16LE, nil
			case binary.BigEndian:
				return pcm.EncodingS16BE, nil
			}
		case 24:
			switch byteOrder {
			case binary.LittleEndian:
				return pcm.EncodingS24LE, nil
			case binary.BigEndian:
				return pcm.EncodingS24BE, nil
			}
		case 32:
			switch byteOrder {
			case binary.LittleEndian:
				return pcm.EncodingS24LE, nil
			case binary.BigEndian:
				return pcm.EncodingS24BE, nil
			}
		}
	case 3:
		switch chunk.BitsPerSample {
		case 32:
			switch byteOrder {
			case binary.LittleEndian:
				return pcm.EncodingF32LE, nil
			case binary.BigEndian:
				return pcm.EncodingF32BE, nil
			}
		}
	}
	return pcm.EncodingUnknown, fmt.Errorf("unsupported format: audioFormat=%d, bitsPerSample=%d", chunk.AudioFormat, chunk.BitsPerSample)
}
