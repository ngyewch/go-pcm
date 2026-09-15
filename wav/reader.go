package wav

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ngyewch/go-pcm"
)

// Reader is a WAV file reader.
type Reader struct {
	f        *os.File
	header   *fmtSubChunk
	dataLen  uint32
	encoding pcm.Encoding
}

// NewReader creates a new Reader.
func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var byteOrder binary.ByteOrder
	var header *fmtSubChunk
	var dataLen uint32

	err = func() error {
		helper := readerHelper{
			r: f,
		}

		chunkID, err := helper.ReadString(4)
		if err != nil {
			return err
		}

		switch chunkID {
		case "RIFF":
			byteOrder = binary.LittleEndian
		case "RIFX":
			byteOrder = binary.BigEndian
		default:
			return fmt.Errorf("unknown chunk ID: %s", chunkID)
		}

		helper = readerHelper{
			r:         f,
			byteOrder: byteOrder,
		}

		_, err = helper.ReadUint32()
		if err != nil {
			return err
		}

		format, err := helper.ReadString(4)
		if err != nil {
			return err
		}
		if format != "WAVE" {
			return fmt.Errorf("format not supported: %s", format)
		}

		fmtSubChunkProcessed := false

		for {
			subChunkID, err := helper.ReadString(4)
			if err != nil {
				return err
			}
			subChunkSize, err := helper.ReadUint32()
			if err != nil {
				return err
			}

			switch subChunkID {
			case "fmt ":
				fmtSubChunk1, err := readFmtSubChunk(f, int(subChunkSize), byteOrder)
				if err != nil {
					return err
				}
				header = fmtSubChunk1
				fmtSubChunkProcessed = true

			case "data":
				if !fmtSubChunkProcessed {
					return fmt.Errorf("fmt sub-chunk not found")
				}
				dataLen = subChunkSize
				return nil

			default:
				_, err = f.Seek(int64(subChunkSize), io.SeekCurrent)
				if err != nil {
					return err
				}
			}
		}
	}()
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	encoding, err := header.Encoding(byteOrder)
	if err != nil {
		_ = f.Close()
		return nil, err
	}

	return &Reader{
		f:        f,
		header:   header,
		dataLen:  dataLen,
		encoding: encoding,
	}, nil
}

// Close closes the reader.
func (r *Reader) Close() error {
	_ = r.f.Close()
	return nil
}

// Encoding returns the encoding.
func (r *Reader) Encoding() pcm.Encoding {
	return r.encoding
}

// NumChannels returns the number of channels.
func (r *Reader) NumChannels() uint16 {
	return r.header.NumChannels
}

// SampleRate returns the sample rate (Hz).
func (r *Reader) SampleRate() uint32 {
	return r.header.SampleRate
}

// Duration return the duration.
func (r *Reader) Duration() time.Duration {
	return time.Duration((float64(r.dataLen) / float64(r.header.ByteRate)) * float64(time.Second))
}

func (r *Reader) Read(b []byte) (n int, err error) {
	blocks := len(b) / int(r.header.BlockAlign)
	adjustedBufferSize := blocks * int(r.header.BlockAlign)
	return r.f.Read(b[0:adjustedBufferSize])
}
