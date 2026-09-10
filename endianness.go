package pcm

// Endianness represents byte ordering.
type Endianness string

const (
	// EndiannessUnknown represents an unknown byte ordering.
	EndiannessUnknown Endianness = "unknown"
	// EndiannessLittleEndian represents little-endian byte ordering.
	EndiannessLittleEndian Endianness = "le"
	// EndiannessBigEndian represents big-endian byte ordering.
	EndiannessBigEndian Endianness = "be"
)

func (endianness Endianness) String() string {
	return string(endianness)
}

// ByteOrder returns the byte order.
func (endianness Endianness) ByteOrder() ByteOrder {
	switch endianness {
	case EndiannessLittleEndian:
		return LittleEndian
	case EndiannessBigEndian:
		return BigEndian
	default:
		return nil
	}
}

// ValueGetter returns the value getter.
func (endianness Endianness) ValueGetter() ValueGetter {
	switch endianness {
	case EndiannessLittleEndian:
		return LittleEndianValueGetter
	case EndiannessBigEndian:
		return BigEndianValueGetter
	default:
		return nil
	}
}

// ValuePutter returns the value putter.
func (endianness Endianness) ValuePutter() ValuePutter {
	switch endianness {
	case EndiannessLittleEndian:
		return LittleEndianValuePutter
	case EndiannessBigEndian:
		return BigEndianValuePutter
	default:
		return nil
	}
}

// IsReverseOf returns true if the byte slice representation is the reverse order of the other.
func (endianness Endianness) IsReverseOf(other Endianness) bool {
	if endianness == other {
		return false
	}
	switch endianness {
	case EndiannessLittleEndian:
		return other == EndiannessBigEndian
	case EndiannessBigEndian:
		return other == EndiannessLittleEndian
	default:
		return false
	}
}
