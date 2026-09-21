package pcm

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

// Encoding is an encoding.
type Encoding string

const (
	// EncodingUnknown represents an unknown encoding.
	EncodingUnknown Encoding = "unknown"
	// EncodingS16LE represents a signed 16-bit integer encoding (little-endian).
	EncodingS16LE Encoding = "s16le"
	// EncodingS24LE represents a signed 24-bit integer encoding (little-endian).
	EncodingS24LE Encoding = "s24le"
	// EncodingS32LE represents a signed 32-bit integer encoding (little-endian).
	EncodingS32LE Encoding = "s32le"
	// EncodingF32LE represents a 32-bit floating point encoding (little-endian).
	EncodingF32LE Encoding = "f32le"
	// EncodingS16BE represents a signed 16-bit integer encoding (big-endian).
	EncodingS16BE Encoding = "s16be"
	// EncodingS24BE represents a signed 24-bit integer encoding (big-endian).
	EncodingS24BE Encoding = "s24be"
	// EncodingS32BE represents a signed 32-bit integer encoding (big-endian).
	EncodingS32BE Encoding = "s32be"
	// EncodingF32BE represents a 32-bit floating point encoding (big-endian).
	EncodingF32BE Encoding = "f32be"
)

func (encoding Encoding) String() string {
	return string(encoding)
}

// BaseEncoding returns the base encoding.
func (encoding Encoding) BaseEncoding() BaseEncoding {
	switch encoding {
	case EncodingS16LE, EncodingS16BE:
		return BaseEncodingS16
	case EncodingS24LE, EncodingS24BE:
		return BaseEncodingS24
	case EncodingS32LE, EncodingS32BE:
		return BaseEncodingS32
	case EncodingF32LE, EncodingF32BE:
		return BaseEncodingF32
	default:
		return BaseEncodingUnknown
	}
}

// ByteOrder returns the byte order.
func (encoding Encoding) ByteOrder() ByteOrder {
	switch encoding {
	case EncodingS16LE, EncodingS24LE, EncodingS32LE, EncodingF32LE:
		return LittleEndian
	case EncodingS16BE, EncodingS24BE, EncodingS32BE, EncodingF32BE:
		return BigEndian
	default:
		return nil
	}
}

// BitsPerSample returns the number of bits used to represent a sample.
func (encoding Encoding) BitsPerSample() int {
	return encoding.BaseEncoding().BitsPerSample()
}

// BytesPerSample returns the number of bytes used to represent a sample.
func (encoding Encoding) BytesPerSample() int {
	return encoding.BaseEncoding().BytesPerSample()
}

// Float64Func returns a function for converting a byte slice representation of the sample to a float64 value.
func (encoding Encoding) Float64Func() func(b []byte) float64 {
	f := encoding.BaseEncoding().Float64Func(encoding.ByteOrder().ValueGetter())
	if f == nil {
		return func(_ []byte) float64 {
			return math.NaN()
		}
	}
	return func(b []byte) float64 {
		return f(b)
	}
}

// PutFloat64Func returns a function for converting a float64 value to its byte slice representation.
func (encoding Encoding) PutFloat64Func() func(b []byte, v float64) {
	f := encoding.BaseEncoding().PutFloat64Func(encoding.ByteOrder().ValuePutter())
	if f == nil {
		return func(_ []byte, _ float64) {
			// do nothing
		}
	}
	return func(b []byte, v float64) {
		f(b, v)
	}
}

// Float64 converts a byte slice representation of the sample to a float64 value.
func (encoding Encoding) Float64(b []byte) float64 {
	f := encoding.Float64Func()
	return f(b)
}

// PutFloat64 converts a float64 value to its byte slice representation.
func (encoding Encoding) PutFloat64(b []byte, v float64) {
	f := encoding.PutFloat64Func()
	f(b, v)
}

// Convert converts a sample from this encoding to the target encoding.
func (encoding Encoding) Convert(sourceBytes []byte, targetEncoding Encoding, targetBytes []byte) error {
	sourceBytesPerSample := encoding.BytesPerSample()
	targetBytesPerSample := targetEncoding.BytesPerSample()

	if len(sourceBytes) != sourceBytesPerSample {
		return fmt.Errorf("incorrect source bytes length")
	}
	if len(targetBytes) != targetBytesPerSample {
		return fmt.Errorf("incorrect target bytes length")
	}

	if encoding == targetEncoding {
		copy(targetBytes, sourceBytes)
		return nil
	}

	if (sourceBytesPerSample == targetBytesPerSample) &&
		(encoding.BaseEncoding() == targetEncoding.BaseEncoding()) &&
		encoding.ByteOrder().IsReverseOf(targetEncoding.ByteOrder()) {
		copy(targetBytes, sourceBytes)
		slices.Reverse(targetBytes)
		return nil
	}

	float64Func := encoding.Float64Func()
	putFloat64Func := targetEncoding.PutFloat64Func()
	v := float64Func(sourceBytes)
	putFloat64Func(targetBytes, v)
	return nil
}

// ConvertSlice converts a sample slice from this encoding to the target encoding.
func (encoding Encoding) ConvertSlice(sourceBytes []byte, targetEncoding Encoding, targetBytes []byte) error {
	sourceBytesPerSample := encoding.BytesPerSample()
	if (len(sourceBytes) % sourceBytesPerSample) != 0 {
		return fmt.Errorf("source bytes length is not a multiple of source bytes per sample")
	}
	targetBytesPerSample := targetEncoding.BytesPerSample()
	if (len(targetBytes) % targetBytesPerSample) != 0 {
		return fmt.Errorf("target bytes length is not a multiple of target bytes per sample")
	}
	sourceSampleCount := len(sourceBytes) / sourceBytesPerSample
	targetSampleCount := len(targetBytes) / targetBytesPerSample
	if sourceSampleCount != targetSampleCount {
		return fmt.Errorf("source sample count is not the same as target sample count")
	}

	if encoding == targetEncoding {
		copy(targetBytes, sourceBytes)
		return nil
	}

	if (sourceBytesPerSample == targetBytesPerSample) &&
		(encoding.BaseEncoding() == targetEncoding.BaseEncoding()) &&
		encoding.ByteOrder().IsReverseOf(targetEncoding.ByteOrder()) {
		copy(targetBytes, sourceBytes)
		for i := 0; i < len(targetBytes); i += targetBytesPerSample {
			b := targetBytes[i : i+targetBytesPerSample]
			slices.Reverse(b)
		}
		return nil
	}

	float64Func := encoding.Float64Func()
	putFloat64Func := targetEncoding.PutFloat64Func()
	sourceOffset := 0
	targetOffset := 0
	for range sourceSampleCount {
		sourceSampleBytes := sourceBytes[sourceOffset : sourceOffset+sourceBytesPerSample]
		targetSampleBytes := targetBytes[targetOffset : targetOffset+targetBytesPerSample]
		v := float64Func(sourceSampleBytes)
		putFloat64Func(targetSampleBytes, v)
		sourceOffset += sourceBytesPerSample
		targetOffset += targetBytesPerSample
	}
	return nil
}

// EncodingFromString returns the Encoding from its ID.
func EncodingFromString(s string) (Encoding, error) {
	switch strings.ToLower(s) {
	case "s16le":
		return EncodingS16LE, nil
	case "s24le":
		return EncodingS24LE, nil
	case "s32le":
		return EncodingS32LE, nil
	case "f32le":
		return EncodingF32LE, nil
	case "s16be":
		return EncodingS16BE, nil
	case "s24be":
		return EncodingS24BE, nil
	case "s32be":
		return EncodingS32BE, nil
	case "f32be":
		return EncodingF32BE, nil
	default:
		return EncodingS16LE, fmt.Errorf("unknown encoding: %s", s)
	}
}
