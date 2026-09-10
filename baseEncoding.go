package pcm

import (
	"fmt"
	"math"
	"strings"
)

// BaseEncoding is a base encoding (independent of byte order).
type BaseEncoding string

const (
	// BaseEncodingUnknown represents an unknown base encoding.
	BaseEncodingUnknown BaseEncoding = "unknown"

	// BaseEncodingS16 represents a signed 16-bit integer encoding.
	BaseEncodingS16 BaseEncoding = "s16"

	// BaseEncodingS24 represents a signed 24-bit integer encoding.
	BaseEncodingS24 BaseEncoding = "s24"

	// BaseEncodingS32 represents a signed 32-bit integer encoding.
	BaseEncodingS32 BaseEncoding = "s32"

	// BaseEncodingF32 represents a 32-bit floating point encoding.
	BaseEncodingF32 BaseEncoding = "f32"
)

func (baseEncoding BaseEncoding) String() string {
	return string(baseEncoding)
}

// BitsPerSample returns the number of bits used to represent a sample.
func (baseEncoding BaseEncoding) BitsPerSample() int {
	switch baseEncoding {
	case BaseEncodingS16:
		return 16
	case BaseEncodingS24:
		return 24
	case BaseEncodingS32, BaseEncodingF32:
		return 32
	default:
		return 0
	}
}

// BytesPerSample returns the number of bytes used to represent a sample.
func (baseEncoding BaseEncoding) BytesPerSample() int {
	return (baseEncoding.BitsPerSample() + 7) / 8
}

// Float64Func returns a function for converting a byte slice representation of the sample to a float64 value.
func (baseEncoding BaseEncoding) Float64Func(valueGetter ValueGetter) func(b []byte) float64 {
	switch baseEncoding {
	case BaseEncodingS16:
		return valueGetter.S16
	case BaseEncodingS24:
		return valueGetter.S24
	case BaseEncodingS32:
		return valueGetter.S32
	case BaseEncodingF32:
		return valueGetter.F32
	default:
		return nil
	}
}

// PutFloat64Func returns a function for converting a float64 value to its byte slice representation.
func (baseEncoding BaseEncoding) PutFloat64Func(valuePutter ValuePutter) func(b []byte, v float64) {
	switch baseEncoding {
	case BaseEncodingS16:
		return valuePutter.PutS16
	case BaseEncodingS24:
		return valuePutter.PutS24
	case BaseEncodingS32:
		return valuePutter.PutS32
	case BaseEncodingF32:
		return valuePutter.PutF32
	default:
		return nil
	}
}

// Float64 converts a byte slice representation of the sample to a float64 value.
func (baseEncoding BaseEncoding) Float64(b []byte, valueGetter ValueGetter) float64 {
	f := baseEncoding.Float64Func(valueGetter)
	if f != nil {
		return f(b)
	}
	return math.NaN()
}

// PutFloat64 converts a float64 value to its byte slice representation.
func (baseEncoding BaseEncoding) PutFloat64(b []byte, v float64, valuePutter ValuePutter) {
	f := baseEncoding.PutFloat64Func(valuePutter)
	if f != nil {
		f(b, v)
	}
}

// BaseEncodingFromString returns the BaseEncoding from its ID.
func BaseEncodingFromString(s string) (BaseEncoding, error) {
	switch strings.ToLower(s) {
	case "s16":
		return BaseEncodingS16, nil
	case "s24":
		return BaseEncodingS24, nil
	case "s32":
		return BaseEncodingS32, nil
	case "f32":
		return BaseEncodingF32, nil
	default:
		return BaseEncodingUnknown, fmt.Errorf("unknown base encoding: %s", s)
	}
}
