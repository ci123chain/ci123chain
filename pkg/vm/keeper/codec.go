package keeper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
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

func (sink Sink) WriteU32(i uint32) {
	sink.writeLittleEndian(i)
}

func (sink Sink) WriteU64(i uint64) {
	sink.writeLittleEndian(i)
}

func (sink Sink) WriteU128(u *types.RustU128) {
	sink.writeRawBytes(u[:])
}

func (sink Sink) WriteI32(i int32) {
	sink.writeLittleEndian(i)
}

//func (sink Sink) WriteI64(i int64) {
//	sink.writeLittleEndian(i)
//}
//
//func (sink Sink) WriteString(s string) {
//	sink.WriteU32(uint32(len(s)))
//	sink.buf.WriteString(s)
//}
//
//func (sink Sink) WriteBytes(b []byte) {
//	sink.WriteU32(uint32(len(b)))
//	sink.writeRawBytes(b)
//}

//func (sink Sink) WriteAddress(addr Address) {
//	sink.writeRawBytes(addr[:])
//}

func (sink Sink) writeRawBytes(b []byte) {
	sink.buf.Write(b)
}

func (sink Sink) writeLittleEndian(i interface{}) {
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

func (sink Sink) ReadBytes() ([]byte, int, error) {
	size, err := sink.ReadU32()
	if err != nil {
		return nil, 0, err
	}
	return sink.nextBytes(int(size))
}

func (sink Sink) ReadString() (string, error) {
	b, _, err := sink.ReadBytes()
	if err == nil && !utf8.Valid(b) {
		return "", errors.New("invalid utf8 string")
	}

	return string(b), err
}

func (sink Sink) nextBytes(size int) ([]byte, int, error) {
	result := make([]byte, size)
	n, err := sink.buf.Read(result)
	return result, n, err
}
