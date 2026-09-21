package wav

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/ngyewch/go-pcm"
)

// Writer is a WAV file writer.
type Writer struct {
	f               *os.File
	channels        uint16
	samplingRate    uint32
	encoding        pcm.Encoding
	sourceEncoding  pcm.Encoding
	bitsPerSample   uint16
	bytesPerSample  uint16
	blockAlign      uint16
	appendByteOrder binary.AppendByteOrder
}

// NewWriter creates a new Writer.
func NewWriter(path string, channels uint16, samplingRate uint32, encoding pcm.Encoding, sourceEncoding pcm.Encoding) (*Writer, error) {
	bitsPerSample := uint16(encoding.BitsPerSample())
	bytesPerSample := (bitsPerSample + 7) / 8
	blockAlign := bytesPerSample * channels

	var containerFormat string
	var appendByteOrder binary.AppendByteOrder
	switch encoding.ByteOrder() {
	case pcm.LittleEndian:
		containerFormat = "RIFF"
		appendByteOrder = binary.LittleEndian
	case pcm.BigEndian:
		containerFormat = "RIFX"
		appendByteOrder = binary.BigEndian
	default:
		return nil, fmt.Errorf("unsupported encoding byte order")
	}

	var audioFormat uint16
	switch encoding.BaseEncoding() {
	case pcm.BaseEncodingS16, pcm.BaseEncodingS24, pcm.BaseEncodingS32:
		audioFormat = 1
	case pcm.BaseEncodingF32:
		audioFormat = 3
	default:
		return nil, fmt.Errorf("unsupported PCM format")
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	err = func() error {
		_, err = f.Write([]byte(containerFormat))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint32(nil, 44-8))
		if err != nil {
			return err
		}

		_, err = f.Write([]byte("WAVE"))
		if err != nil {
			return err
		}

		_, err = f.Write([]byte("fmt "))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint32(nil, 16))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint16(nil, audioFormat))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint16(nil, channels))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint32(nil, samplingRate))
		if err != nil {
			return err
		}

		byteRate := uint32(blockAlign) * samplingRate
		_, err = f.Write(appendByteOrder.AppendUint32(nil, byteRate))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint16(nil, blockAlign))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint16(nil, bitsPerSample))
		if err != nil {
			return err
		}

		_, err = f.Write([]byte("data"))
		if err != nil {
			return err
		}

		_, err = f.Write(appendByteOrder.AppendUint32(nil, 0))
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return &Writer{
		f:               f,
		channels:        channels,
		samplingRate:    samplingRate,
		encoding:        encoding,
		sourceEncoding:  sourceEncoding,
		bitsPerSample:   bitsPerSample,
		bytesPerSample:  bytesPerSample,
		blockAlign:      blockAlign,
		appendByteOrder: appendByteOrder,
	}, nil
}

func (w *Writer) Write(b []byte) (int, error) {
	sourceBytesPerSample := w.sourceEncoding.BytesPerSample()
	targetBytesPerSample := w.encoding.BytesPerSample()
	sourceSampleCount := len(b) / sourceBytesPerSample
	targetBytes := make([]byte, sourceSampleCount*targetBytesPerSample)
	err := w.sourceEncoding.ConvertSlice(b, w.encoding, targetBytes)
	if err != nil {
		return 0, err
	}

	_, err = w.f.Write(targetBytes)
	if err != nil {
		return 0, err
	}

	offset, err := w.f.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	offset32 := uint32(offset)

	_, err = w.f.Seek(4, io.SeekStart)
	if err != nil {
		return 0, err
	}
	_, err = w.f.Write(w.appendByteOrder.AppendUint32(nil, offset32-8))
	if err != nil {
		return 0, err
	}

	_, err = w.f.Seek(40, io.SeekStart)
	if err != nil {
		return 0, err
	}
	_, err = w.f.Write(w.appendByteOrder.AppendUint32(nil, offset32-44))
	if err != nil {
		return 0, err
	}

	_, err = w.f.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return len(b), nil
}

// Close closes the writer.
func (w *Writer) Close() error {
	err := w.f.Close()
	if err != nil {
		return err
	}
	return nil
}
