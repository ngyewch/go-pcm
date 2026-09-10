package pcm

import (
	"math"
)

const (
	maxInt24 = (1 << (24 - 1)) - 1
)

var (
	// LittleEndianValueGetter is the little-endian implementation of ValueGetter.
	LittleEndianValueGetter = valueGetter{
		ByteOrder: LittleEndian,
	}
	// BigEndianValueGetter is the big-endian implementation of ValueGetter.
	BigEndianValueGetter = valueGetter{
		ByteOrder: BigEndian,
	}
	// LittleEndianValuePutter is the little-endian implementation of ValuePutter.
	LittleEndianValuePutter = valuePutter{
		ByteOrder: LittleEndian,
	}
	// BigEndianValuePutter is the big-endian implementation of ValuePutter.
	BigEndianValuePutter = valuePutter{
		ByteOrder: BigEndian,
	}
)

// ValueGetter converts byte slices to a float64 value.
type ValueGetter interface {
	S16(b []byte) float64
	S24(b []byte) float64
	S32(b []byte) float64
	F32(b []byte) float64
}

// ValuePutter converts a float64 value its byte slice representation.
type ValuePutter interface {
	PutS16(b []byte, v float64)
	PutS24(b []byte, v float64)
	PutS32(b []byte, v float64)
	PutF32(b []byte, v float64)
}

type valueGetter struct {
	ByteOrder ByteOrder
}

func (getter valueGetter) S16(b []byte) float64 {
	v := getter.ByteOrder.Int16(b)
	return float64(v) / float64(math.MaxInt16)
}

func (getter valueGetter) S24(b []byte) float64 {
	v := getter.ByteOrder.Int24(b)
	return float64(v) / float64(maxInt24)
}

func (getter valueGetter) S32(b []byte) float64 {
	v := getter.ByteOrder.Int32(b)
	return float64(v) / float64(math.MaxInt32)
}

func (getter valueGetter) F32(b []byte) float64 {
	v := getter.ByteOrder.Float32(b)
	return float64(v)
}

type valuePutter struct {
	ByteOrder ByteOrder
}

func (putter valuePutter) PutS16(b []byte, v float64) {
	putter.ByteOrder.PutInt16(b, quantize[int16](v, math.MaxInt16))
}

func (putter valuePutter) PutS24(b []byte, v float64) {
	putter.ByteOrder.PutInt24(b, quantize[int32](v, maxInt24))
}

func (putter valuePutter) PutS32(b []byte, v float64) {
	putter.ByteOrder.PutInt32(b, quantize[int32](v, math.MaxInt32))
}

func (putter valuePutter) PutF32(b []byte, v float64) {
	putter.ByteOrder.PutFloat32(b, float32(v))
}

func quantize[T int16 | int32 | int64](v float64, maxValue T) T {
	minValue := -maxValue
	n := math.Round(v * float64(maxValue))
	if n > float64(maxValue) {
		return maxValue
	}
	if n < float64(minValue) {
		return minValue
	}
	return T(n)
}
