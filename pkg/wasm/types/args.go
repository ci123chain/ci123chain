package types

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/big"
	"unicode/utf8"
)

var WasmIdent = []byte("\x00\x61\x73\x6D")

type CallContractParam struct {
	Method string   `json:"method"`
	ArgsType []string
	Args   []json.RawMessage `json:"args"`
}

// 验证
func ValidCallArgs(args map[string]interface{}) error {
	for k, _ := range args {
		if err := validKey(k); err != nil {
			return err
		}
	}
	return nil
}

func validKey(typeName string) error {
	if  typeName == typeInt32 ||
		typeName == typeUint32 ||
		typeName == typeUint64 ||
		typeName == typeInt64 ||
		typeName == typeUint128 ||
		typeName == typeInt128 ||
		typeName == typeString ||
		typeName == typeBool {
		return nil
	} else {
		return errors.New("paramtype invalid : " + typeName + validtypesTip)
	}
}


func Serialize(args map[string]json.RawMessage) (res []byte) {
	sink := NewSink(res)
	for types, v := range args {

		switch types {
		case typeInt32:
			b := new(big.Int)
			if err := b.UnmarshalText(v); err == nil && len(b.Bytes()) <= 4 {
				by := b.Bytes()
				by = fillBytes(by, 4)
				bytesBuffer := bytes.NewBuffer(by)
				var x int32
				if err := binary.Read(bytesBuffer, binary.LittleEndian, &x); err == nil {
					sink.WriteI32(x)
				} else {
					panic("invalid value of int32: " + b.String())
				}
			} else {
				panic("invalid value of int32: ")
			}
		case typeUint32:
			b := new(big.Int)
			if err := b.UnmarshalText(v); err == nil && len(b.Bytes()) <= 4 {
				by := b.Bytes()
				by = fillBytes(by, 4)
				bytesBuffer := bytes.NewBuffer(by)
				var x uint32
				if err := binary.Read(bytesBuffer, binary.LittleEndian, &x); err == nil {
					sink.WriteU32(x)
				} else {
					panic("invalid value of uint32: " + b.String())
				}
			} else {
				panic("invalid value of int32: ")
			}
		case typeInt64:
			b := new(big.Int)
			lth := 8
			if err := b.UnmarshalText(v); err == nil && len(b.Bytes()) <= lth {
				by := b.Bytes()
				by = fillBytes(by, lth)
				bytesBuffer := bytes.NewBuffer(by)
				var x int64
				if err := binary.Read(bytesBuffer, binary.LittleEndian, &x); err == nil {
					sink.WriteI64(x)
				} else {
					panic("invalid value of int64: " + b.String())
				}
			} else {
				panic("invalid value of int64: ")
			}
		case typeUint64:
			b := new(big.Int)
			lth := 8
			if err := b.UnmarshalText(v); err == nil && len(b.Bytes()) <= lth {
				by := b.Bytes()
				by = fillBytes(by, lth)
				bytesBuffer := bytes.NewBuffer(by)
				var x uint64
				if err := binary.Read(bytesBuffer, binary.LittleEndian, &x); err == nil {
					sink.WriteU64(x)
				} else {
					panic("invalid value of uint64: " + b.String())
				}
			} else {
				panic("invalid value of uint64: ")
			}
		case typeInt128:
			b := new(big.Int)
			if err := b.UnmarshalText(v); err == nil && len(v) <= 128 {
				sink.WriteI128(NewRustI128(b))
			} else {
				panic("invalid value of int128: ")
			}
		case typeUint128:
			b := new(big.Int)
			if err := b.UnmarshalText(v); err == nil && len(v) <= 128 {
				sink.WriteU128(NewRustU128(b))
			} else {
				panic("invalid value of uint128: ")
			}
		case typeString:
			var value string
			if err := json.Unmarshal(v, &value); err !=nil {
				panic("invalid value of string")
			}
			//字符串必须是合法的utf8字符串
			if utf8.ValidString(value) {
				sink.WriteString(value)
			} else {
			panic("invalid utf8 string")
			}
		case typeBool:
			var value bool
			if err := json.Unmarshal(v, &value); err !=nil {
				panic("invalid value of bool")
			}
			sink.WriteBool(value)
		default:
			panic("unexpected type")
		}
	}

	return sink.Bytes()
}

func fillBytes(bz []byte, length int) []byte {
	for i := 0; i < len(bz)/2; i++ {
		bz[i], bz[len(bz)-1-i] = bz[len(bz)-1-i], bz[i]
	}
	for i := len(bz); i < length; i++ {
		bz = append(bz, 0)
	}
	return bz
}

const typeInt32 = "int32"
const typeUint32 = "uint32"
const typeInt64 = "int64"
const typeUint64 = "uint64"
const typeInt128 = "int128"
const typeUint128 = "uint128"
const typeString = "string"
const typeBool = "bool"

const validtypesTip = "; valid types: " + typeInt32 + typeUint32 + typeInt64 + typeUint64 + typeInt128 + typeUint128 + typeString + typeBool
