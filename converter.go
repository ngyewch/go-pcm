package pcm

import (
	"io"
)

// Converter is a io.Writer that converts from one encoding to another.
type Converter struct {
	w              io.Writer
	sourceEncoding Encoding
	targetEncoding Encoding
}

// NewConverter creates a new Converter.
func NewConverter(w io.Writer, sourceEncoding Encoding, targetEncoding Encoding) *Converter {
	return &Converter{
		w:              w,
		sourceEncoding: sourceEncoding,
		targetEncoding: targetEncoding,
	}
}

func (converter *Converter) Write(b []byte) (int, error) {
	targetBytes, err := converter.sourceEncoding.Convert(b, converter.targetEncoding)
	if err != nil {
		return 0, err
	}
	_, err = converter.w.Write(targetBytes)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}
