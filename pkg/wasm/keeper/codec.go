package keeper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unicode/utf8"
)

type Sink struct {
	buf *bytes.Buffer
}

func NewSink(raw []byte) Sink {
	return Sink{
		bytes.NewBuffer(raw),
	}
}

func (sink Sink) WriteString(s string) {
	sink.WriteU32(uint32(len(s)))
	sink.buf.WriteString(s)
}

func (sink Sink) WriteU32(i uint32) {
	_ = binary.Write(sink.buf, binary.LittleEndian, i)
}

func (sink Sink) Bytes() []byte {
	return sink.buf.Bytes()
}

func (sink Sink) ReadBool() (bool, error) {
	b, err := sink.buf.ReadByte()
	if err != nil {
		return false, err
	}
	if b == 0 {
		return false, nil
	}
	return true, nil
}

func (sink Sink) ReadByte() (byte, error) {
	return sink.buf.ReadByte()
}

func (sink Sink) ReadU32() (result uint32, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) ReadI64() (result int64, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) ReadBytes(size int) ([]byte, int, error) {
	result := make([]byte, size)
	n, err := sink.buf.Read(result)
	return result, n, err
}

func (sink Sink) ReadString() (string, error) {
	size, err := sink.ReadU32()
	if err != nil {
		return "", err
	}

	b, _, err := sink.ReadBytes(int(size))
	if err == nil && !utf8.Valid(b) {
		return "", errors.New("invalid utf8 string")
	}

	return string(b), err
}
