package pcm

import (
	"encoding/binary"
	"math"
)

// ByteOrder specifies how to convert byte slices into various types and vice versa.
type ByteOrder interface {
	binary.ByteOrder

	Int16([]byte) int16
	Int32([]byte) int32
	Int64([]byte) int64
	Float32([]byte) float32
	PutInt16(b []byte, v int16)
	PutInt32(b []byte, v int32)
	PutInt64(b []byte, v int64)
	PutFloat32(b []byte, v float32)
	Int24(b []byte) int32
	Uint24(b []byte) uint32
	PutUint24(b []byte, v uint32)
	PutInt24(b []byte, v int32)
	ValueGetter() ValueGetter
	ValuePutter() ValuePutter
	IsReverseOf(otherByteOrder ByteOrder) bool
}

type pcmByteOrder struct {
	binary.ByteOrder
}

func newPCMByteOrder(byteOrder binary.ByteOrder) *pcmByteOrder {
	return &pcmByteOrder{
		ByteOrder: byteOrder,
	}
}

var (
	// BigEndian is the big-endian implementation of ByteOrder
	BigEndian = newPCMByteOrder(binary.BigEndian)
	// LittleEndian is the little-endian implementation of ByteOrder
	LittleEndian = newPCMByteOrder(binary.LittleEndian)
)

func (byteOrder pcmByteOrder) Int16(b []byte) int16 {
	return int16(byteOrder.Uint16(b))
}

func (byteOrder pcmByteOrder) Int32(b []byte) int32 {
	return int32(byteOrder.Uint32(b))
}

func (byteOrder pcmByteOrder) Int64(b []byte) int64 {
	return int64(byteOrder.Uint64(b))
}

func (byteOrder pcmByteOrder) Float32(b []byte) float32 {
	return math.Float32frombits(byteOrder.Uint32(b))
}

func (byteOrder pcmByteOrder) PutInt16(b []byte, v int16) {
	byteOrder.PutUint16(b, uint16(v))
}

func (byteOrder pcmByteOrder) PutInt32(b []byte, v int32) {
	byteOrder.PutUint32(b, uint32(v))
}

func (byteOrder pcmByteOrder) PutInt64(b []byte, v int64) {
	byteOrder.PutUint64(b, uint64(v))
}

func (byteOrder pcmByteOrder) PutFloat32(b []byte, v float32) {
	byteOrder.PutUint32(b, math.Float32bits(v))
}

func (byteOrder pcmByteOrder) Uint24(b []byte) uint32 {
	var tmp []byte
	switch byteOrder.ByteOrder {
	case binary.LittleEndian:
		tmp = []byte{b[0], b[1], b[2], 0}
	case binary.BigEndian:
		tmp = []byte{0, b[0], b[1], b[2]}
	default:
		return 0
	}
	return byteOrder.Uint32(tmp)
}

func (byteOrder pcmByteOrder) Int24(b []byte) int32 {
	return int32(byteOrder.Uint24(b))
}

func (byteOrder pcmByteOrder) PutUint24(b []byte, v uint32) {
	var tmp [4]byte
	byteOrder.PutUint32(tmp[:], v)
	switch byteOrder.ByteOrder {
	case binary.LittleEndian:
		copy(b, tmp[0:3])
	case binary.BigEndian:
		copy(b, tmp[1:4])
	}
}

func (byteOrder pcmByteOrder) PutInt24(b []byte, v int32) {
	byteOrder.PutUint32(b, uint32(v))
}

func (byteOrder pcmByteOrder) ValueGetter() ValueGetter {
	return valueGetter{
		ByteOrder: byteOrder,
	}
}

func (byteOrder pcmByteOrder) ValuePutter() ValuePutter {
	return valuePutter{
		ByteOrder: byteOrder,
	}
}

func (byteOrder pcmByteOrder) IsReverseOf(otherByteOrder ByteOrder) bool {
	if byteOrder == otherByteOrder {
		return false
	}
	otherPCMByteOrder, ok := otherByteOrder.(*pcmByteOrder)
	if !ok {
		return false
	}
	switch byteOrder.ByteOrder {
	case binary.LittleEndian:
		return otherPCMByteOrder.ByteOrder == binary.BigEndian
	case binary.BigEndian:
		return otherPCMByteOrder.ByteOrder == binary.LittleEndian
	default:
		return false
	}
}
