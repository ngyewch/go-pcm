package wav

import (
	"encoding/binary"
	"io"
)

type readerHelper struct {
	r         io.Reader
	byteOrder binary.ByteOrder
}

func (helper readerHelper) ReadString(n int) (string, error) {
	buffer := make([]byte, n)
	_, err := io.ReadFull(helper.r, buffer)
	if err != nil {
		return "", err
	}
	return string(buffer), nil
}

func (helper readerHelper) ReadUint16() (uint16, error) {
	var buffer [2]byte
	_, err := io.ReadFull(helper.r, buffer[:])
	if err != nil {
		return 0, err
	}
	return helper.byteOrder.Uint16(buffer[:]), nil
}

func (helper readerHelper) ReadUint32() (uint32, error) {
	var buffer [4]byte
	_, err := io.ReadFull(helper.r, buffer[:])
	if err != nil {
		return 0, err
	}
	return helper.byteOrder.Uint32(buffer[:]), nil
}

func (helper readerHelper) ReadBytes(n int) ([]byte, error) {
	buffer := make([]byte, n)
	_, err := io.ReadFull(helper.r, buffer)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
