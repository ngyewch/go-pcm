package pcm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoding(t *testing.T) {
	testEncodingFromString(t, "S16LE", EncodingS16LE)
	testEncodingFromString(t, "S24LE", EncodingS24LE)
	testEncodingFromString(t, "S32LE", EncodingS32LE)
	testEncodingFromString(t, "F32LE", EncodingF32LE)
	testEncodingFromString(t, "S16BE", EncodingS16BE)
	testEncodingFromString(t, "S24BE", EncodingS24BE)
	testEncodingFromString(t, "S32BE", EncodingS32BE)
	testEncodingFromString(t, "F32BE", EncodingF32BE)
	testEncodingFromString(t, "s16le", EncodingS16LE)
	testEncodingFromString(t, "s24le", EncodingS24LE)
	testEncodingFromString(t, "s32le", EncodingS32LE)
	testEncodingFromString(t, "f32le", EncodingF32LE)
	testEncodingFromString(t, "s16be", EncodingS16BE)
	testEncodingFromString(t, "s24be", EncodingS24BE)
	testEncodingFromString(t, "s32be", EncodingS32BE)
	testEncodingFromString(t, "f32be", EncodingF32BE)

	testEncodingToString(t, EncodingS16LE, "s16le")
	testEncodingToString(t, EncodingS24LE, "s24le")
	testEncodingToString(t, EncodingS32LE, "s32le")
	testEncodingToString(t, EncodingF32LE, "f32le")
	testEncodingToString(t, EncodingS16BE, "s16be")
	testEncodingToString(t, EncodingS24BE, "s24be")
	testEncodingToString(t, EncodingS32BE, "s32be")
	testEncodingToString(t, EncodingF32BE, "f32be")

	testEncodingByteOrder(t, EncodingS16LE, LittleEndian)
	testEncodingByteOrder(t, EncodingS24LE, LittleEndian)
	testEncodingByteOrder(t, EncodingS32LE, LittleEndian)
	testEncodingByteOrder(t, EncodingF32LE, LittleEndian)
	testEncodingByteOrder(t, EncodingS16BE, BigEndian)
	testEncodingByteOrder(t, EncodingS24BE, BigEndian)
	testEncodingByteOrder(t, EncodingS32BE, BigEndian)
	testEncodingByteOrder(t, EncodingF32BE, BigEndian)
}

func testEncodingFromString(t *testing.T, s string, expected Encoding) {
	encoding, err := EncodingFromString(s)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, encoding)
	}
}

func testEncodingToString(t *testing.T, encoding Encoding, expected string) {
	assert.Equal(t, expected, encoding.String())
}

func testEncodingByteOrder(t *testing.T, encoding Encoding, expected ByteOrder) {
	assert.Equal(t, expected, encoding.ByteOrder())
}
