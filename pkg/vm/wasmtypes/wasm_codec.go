package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"unicode/utf8"
)

type number128 [16]byte

func newNumber128(i *big.Int) *number128 {
	ib := i.Bytes() // big-endian

	size := len(ib)
	if size > 16 {
		panic("u128最大16字节") //链上处理
	}

	// little-endian
	for i := 0; i < size/2; i++ {
		ib[i], ib[size-1-i] = ib[size-1-i], ib[i]
	}

	// 补全
	for i := 0; i+size < 16; i++ {
		ib = append(ib, 0)
	}

	var n number128
	copy(n[:], ib)
	return &n
}

type RustU128 number128

func NewRustU128(i *big.Int) *RustU128 {
	n128 := newNumber128(i)
	var u128 RustU128
	copy(u128[:], n128[:])
	return &u128
}

func (u128 *RustU128) Bytes() []byte {
	return u128[:]
}

type RustI128 number128

func NewRustI128(i *big.Int) *RustI128 {
	n128 := newNumber128(i)
	var i128 RustI128
	copy(i128[:], n128[:])
	return &i128
}

func (i128 *RustI128) Bytes() []byte {
	return i128[:]
}

type Sink struct {
	buf *bytes.Buffer
}

func NewSink(raw []byte) Sink {
	return Sink{
		bytes.NewBuffer(raw),
	}
}

// Write ------------------------------------------------------------------------------------

func (sink Sink) WriteBool(b bool) {
	// Boolean values encode as one byte: 1 for true, and 0 for false.
	sink.writeLittleEndian(b)
}

func (sink Sink) WriteBytes(b []byte) {
	// length
	sink.writeLittleEndian(uint32(len(b)))
	sink.writeRawBytes(b)
}

func (sink Sink) WriteU32(i uint32) {
	sink.writeLittleEndian(i)
}

func (sink Sink) WriteI32(i int32) {
	sink.writeLittleEndian(i)
}

func (sink Sink) WriteU64(i uint64) {
	sink.writeLittleEndian(i)
}

func (sink Sink) WriteI64(i int64) {
	sink.writeLittleEndian(i)
}

func (sink Sink) WriteU128(u128 *RustU128) {
	sink.writeRawBytes(u128.Bytes())
}

func (sink Sink) WriteI128(i128 *RustI128) {
	sink.writeRawBytes(i128.Bytes())
}

func (sink Sink) WriteString(s string) {
	// length
	sink.writeLittleEndian(uint32(len(s)))
	sink.buf.WriteString(s)
}

func (sink Sink) writeRawBytes(b []byte) {
	sink.buf.Write(b)
}

func (sink Sink) writeLittleEndian(i interface{}) {
	_ = binary.Write(sink.buf, binary.LittleEndian, i)
}

// Read -------------------------------------------------------------------------------------

func (sink Sink) Bytes() []byte {
	return sink.buf.Bytes()
}

func (sink Sink) ReadByte() (byte, error) {
	return sink.buf.ReadByte()
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

func (sink Sink) ReadBytes() ([]byte, int, error) {
	size, err := sink.ReadUSize()
	if err != nil {
		return nil, 0, err
	}
	return sink.nextBytes(int(size))
}

func (sink Sink) ReadU32() (result uint32, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) ReadI32() (result int32, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) ReadU64() (result uint64, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) ReadI64() (result int64, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) ReadString() (string, error) {
	size, err := sink.ReadUSize()
	if err != nil {
		return "", err
	}
	b, _, err := sink.nextBytes(int(size))
	if err == nil && !utf8.Valid(b) {
		return "", errors.New("invalid utf8 string")
	}

	return string(b), err
}

func (sink Sink) ReadUSize() (result uint32, err error) {
	err = binary.Read(sink.buf, binary.LittleEndian, &result)
	return
}

func (sink Sink) nextBytes(size int) ([]byte, int, error) {
	result := make([]byte, size)
	n, err := sink.buf.Read(result)
	return result, n, err
}
