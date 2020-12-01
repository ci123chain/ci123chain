package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"math/big"
	"strings"
)

const typeInt32 = "int32"
const typeUint32 = "uint32"
const typeInt64 = "int64"
const typeUint64 = "uint64"
const typeInt128 = "int128"
const typeUint128 = "uint128"
const typeString = "string"
const typeBool = "bool"
const space = " "
const validtypesTip = "; valid types: " + typeInt32 + space + typeUint32 + space + typeInt64 + space + typeUint64 + space  + typeInt128 + space  + typeUint128 + space  + typeString + space  + typeBool

var WasmIdent = []byte("\x00\x61\x73\x6D")

func ValidType(signature *utils.Signature) error {
	for _, v := range signature.Args {
		if err := validKey(v); err != nil {
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
		typeName == typeBool ||
		typeName == "" {
		return nil
	} else {
		return errors.New("paramtype invalid : " + typeName + validtypesTip)
	}
}

func ArgsToInput(args utils.CallData) (res []byte, err error){
	sink := NewSink(res)
	sig, err := utils.ParseSignature(strings.Replace(args.Method, " ", "", -1))
	if err != nil {
		return nil, err
	}
	err = ValidType(sig)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(args.Args); i++ {
		switch sig.Args[i] {
		case typeString:
			var str string
			err := json.Unmarshal(args.Args[i], &str)
			if err != nil{
				return nil, errors.New(fmt.Sprintf("parse %d arg to string failed", i+1))
			}
			sink.WriteString(str)
			continue
		case typeInt32:
			var i32 int32
			err := json.Unmarshal(args.Args[i], &i32)
			if err != nil{
				return nil, errors.New(fmt.Sprintf("parse %d arg to i32 failed", i+1))
			}
			sink.WriteI32(i32)
			continue
		case typeInt64:
			var i64 int64
			err := json.Unmarshal(args.Args[i], &i64)
			if err != nil{
				return nil, errors.New(fmt.Sprintf("parse %d arg to i64 failed", i+1))
			}
			sink.WriteI64(i64)
			continue
		case typeInt128:
			x := new(big.Int)
			err := x.UnmarshalText(args.Args[i])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parse %d arg to i128 failed", i+1))

			}
			sink.WriteI128(NewRustI128(x))
			continue
		case typeUint32:
			var u32 uint32
			err := json.Unmarshal(args.Args[i], &u32)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parse %d arg to u32 failed", i+1))
			}
			sink.WriteU32(u32)
			continue
		case typeUint64:
			var u64 uint64
			err := json.Unmarshal(args.Args[i], &u64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parse %d arg to u64 failed", i+1))
			}
			sink.WriteU64(u64)
			continue
		case typeUint128:
			x := new(big.Int)
			err := x.UnmarshalText(args.Args[i])
			if err != nil {
				return nil, errors.New(fmt.Sprintf("parse %d arg to u128 failed", i+1))
			}
			sink.WriteU128(NewRustU128(x))
			continue
		case typeBool:
			var b bool
			err := json.Unmarshal(args.Args[i], &b)
			if err != nil{
				return nil, errors.New(fmt.Sprintf("parse %d arg to bool failed", i+1))
			}
			sink.WriteBool(b)
			continue
		}
	}

	input := sink.Bytes()
	return input, nil
}