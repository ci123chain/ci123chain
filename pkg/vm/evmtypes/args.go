package evmtypes

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"math/big"
	"strconv"
	"strings"
)

const typeInt = "int"
const typeUint = "uint"
const typeBytes = "bytes"
const typeString = "string"
const typeBool = "bool"
const typeAddress = "address"
const bigInt = "64"

func ValidType(signature *utils.Signature) error {
	for _, v := range signature.Args {
		if err := validKey(v); err != nil {
			return err
		}
	}
	return nil
}

func validKey(typeName string) error {
	if strings.HasPrefix(typeName, typeInt) {
		if len(typeName) > 3 {
			size, err := strconv.Atoi(typeName[3:])
			if err != nil {
				return errors.New("paramtype is valid:" + typeName)
			}
			if size % 8 != 0 || size < 8 || size > 256 {
				return errors.New("int size is valid")
			}
		}
	} else if strings.HasPrefix(typeName, typeUint) {
		if len(typeName) > 4 {
			size, err := strconv.Atoi(typeName[4:])
			if err != nil {
				return errors.New("paramtype is valid:" + typeName)
			}
			if size % 8 != 0 || size < 8 || size > 256 {
				return errors.New("uint size is valid")
			}
		}
	} else if strings.HasPrefix(typeName, typeBytes) {
		if len(typeName) > 4 {
			size, err := strconv.Atoi(typeName[4:])
			if err != nil {
				return errors.New("paramtype is valid:" + typeName)
			}
			if size < 1 || size > 32 {
				return errors.New("uint size is valid")
			}
		}
	} else if typeName == typeAddress || typeName == typeString || typeName == typeBool || typeName == ""{
		return nil
	} else {
		return errors.New("paramtype invalid : " + typeName )
	}
	return nil
}

func EVMEncode(args utils.CallData) ([]byte, error) {
	sig, err := utils.ParseSignature(strings.Replace(args.Method, " ", "", -1))
	if err != nil {
		return nil, err
	}
	err = ValidType(sig)
	if err != nil {
		return nil, err
	}
	var rawArg []interface{}
	for i, v := range args.Args {
		switch utils.ElementaryName(sig.Args[i]) {
		case typeString:
			var str string
			if err := json.Unmarshal(v, &str); err != nil {
				return nil , errors.New("parse " + sig.Args[i] + " to string failed" )
			}
			rawArg = append(rawArg, str)
			continue
		case typeBytes:
			rawArg = append(rawArg, []byte(v))
			continue
		case typeBool:
			var boolean bool
			if err := json.Unmarshal(v, &boolean); err != nil {
				return nil , errors.New("parse " + sig.Args[i] + " to boolean failed" )
			}
			rawArg = append(rawArg, boolean)
			continue
		case typeAddress:
			//todo 0x format
			var addr string
			x  := new(big.Int)
			if err := json.Unmarshal(v, &addr); err != nil {
				if err := x.UnmarshalText(v); err != nil {
					return nil , errors.New("parse " + sig.Args[i] + " to address failed" )
				}
				rawArg = append(rawArg, x)
			} else {
				rawArg = append(rawArg, addr)
			}
			continue
		default:
			if isArray(sig.Args[i]) { //array
				realType := sig.Args[i][:strings.LastIndex(sig.Args[i], "[")]
				switch realType {
				case typeInt:
					var intArr []json.RawMessage
					var bigIntArr []*big.Int
					if err := json.Unmarshal(v, &intArr); err != nil {
						return nil, errors.New("parse " + realType +  "array error")
					}
					for _, k := range intArr {
						x := new(big.Int)
						if err := x.UnmarshalText(k); err != nil {
							return nil , errors.New("parse " + sig.Args[i] + " to bigInt failed" )
						}
						bigIntArr = append(bigIntArr, x)
					}
					rawArg = append(rawArg, bigIntArr)
					continue
				case typeUint:
					var uintArr []json.RawMessage
					var bigIntArr []*big.Int
					if err := json.Unmarshal(v, &uintArr); err != nil {
						return nil, errors.New("parse " + realType +  "array error")
					}
					for _, v := range uintArr {
						x := new(big.Int)
						if err := x.UnmarshalText(v); err != nil {
							return nil , errors.New("parse " + sig.Args[i] + " to bigInt failed" )
						}
						bigIntArr = append(bigIntArr, x)
					}
					rawArg = append(rawArg, bigIntArr)
					continue
				case typeBytes:
					var bys [][]byte
					if err := json.Unmarshal(v, &bys); err != nil {
						return nil, errors.New("parse " + realType +  "array error")
					}
					rawArg = append(rawArg, bys)
					continue
				default:
					return nil, errors.New("parse arrayType: " + realType +  " error")
				}
			} else { //not array
				if strings.HasPrefix(sig.Args[i], typeInt) || strings.HasPrefix(sig.Args[i], typeUint){
					x := new(big.Int)
					if err := x.UnmarshalText(v); err != nil {
						return nil , errors.New("parse " + sig.Args[i] + " to bigInt failed" )
					}
					rawArg = append(rawArg, x)
				} else if strings.HasPrefix(sig.Args[i], typeBytes) {
					rawArg = append(rawArg, []byte(v))
				}
			}
		}
	}
	return append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, rawArg)...), nil
}

func isArray (name string) bool {
	return strings.LastIndex(name, "]") == len(name) - 1
}