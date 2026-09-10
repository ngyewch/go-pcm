package pcm

import "io"

// Source is a PCM data source.
type Source interface {
	io.Reader

	// Encoding returns the encoding.
	Encoding() Encoding

	// NumChannels returns the number of channels.
	NumChannels() uint16

	// SampleRate returns the sample rate (Hz).
	SampleRate() uint32
}
