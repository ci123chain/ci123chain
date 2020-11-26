package utils

import (
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"encoding/json"
	"errors"
)

type CallData struct {
	Method string `json:"method"`
	Args []json.RawMessage `json:"args,omitempty"`
}

type Signature struct {
	Method 	string
	Args 	[]string
	Retargs []string
}

//someMethod(bytes,uint)
func ParseSignature(sig string) (*Signature, error){
	re := regexp.MustCompile(`^(\w+)\((.*)\)$`)
	matchedSig := re.FindAllStringSubmatch(sig, -1)
	if len(matchedSig) != 0 {
		if len(matchedSig[0]) < 3 {
			return nil, errors.New("Invalid method signature")
		}
	} else {
		return nil, errors.New("Invalid method signature")
	}

	//re = regexp.MustCompile(`^(.+)\)$`)
	//matchedRet := re.FindAllStringSubmatch(matchedSig[0][3], -1)
	//if len(matchedRet) != 0 { //ret
	//	return signature{
	//		Method: matchedSig[0][1],
	//		Args: strings.Split(matchedRet[0][1], ","),
	//		Retargs: strings.Split(matchedRet[0][2], ","),
	//	}
	//} else { // no ret
	return &Signature {
		Method: matchedSig[0][1],
		Args:   strings.Split(matchedSig[0][2], ","),
	}, nil
	//}
}

type doc interface {
	DocType() string
}

type evmABI struct {
	Anonymous bool 			`json:"anonymous,omitempty"`
	Input []struct{
		Indexed bool 		`json:"indexed,omitempty"`
		InternalType string `json:"internalType"`
		Name string 		`json:"name"`
		Type string 		`json:"type"`
	} 						`json:"input,omitempty"`
	Output []struct{
		Indexed bool 		`json:"indexed,omitempty"`
		InternalType string `json:"internalType"`
		Name string 		`json:"name"`
		Type string 		`json:"type"`
	} 						`json:"output,omitempty"`
	Name string 			`json:"name,omitempty"`
	Type string 			`json:"type"`
	Constant bool 			`json:"constant,omitempty"`
	Payable bool 			`json:"payable,omitempty"`
	StateMutability string 	`json:"state_mutability,omitempty"`
}

func(doc evmABI) DocType() string {return "evmABI"}

type wasmDoc struct {
	InvokeName string `json:"invoke_name"`
	ExportName string `json:"export_name"`
	InputType []string `json:"input_type"`
	OutputType string
}

func(doc wasmDoc) DocType() string {return "wasmDoc"}

//[ 'uint32', 'bool' ], [ 69, 1 ]
func RawEncode(paramTypes []string, values []interface{}) []byte {
	if len(paramTypes) == 1 && paramTypes[0] == "" {
		return nil
	}
	var output, data []byte
	headLength := 0
	for _, v := range paramTypes{
		if isArray(v){
			var size = parseTypeArray(v)
			if size != "dynamic" && size != ""{
				length, err := strconv.Atoi(size)
				if err != nil {
					panic("err parseTypeArray length")
				}
				headLength += 32 * length
			} else {
				headLength += 32
			}
		} else {
			headLength += 32
		}
	}

	for i := 0; i < len(paramTypes); i++ {
		param := ElementaryName(paramTypes[i])
		value := values[i]
		cur := encodeSingle(param, value)

		if isDynamic(param) {
			output = append(output, encodeSingle("uint256", headLength)...)
			data = append(data, cur...)
			headLength += len(cur)
		} else {
			output = append(output, cur...)
		}
	}

	output = append(output, data...)
	return output
}

//'baz', [ 'uint32', 'bool' ]
func MethodID(method string, paramTypes []string) []byte {
	return EventID(method, paramTypes)[:4]
}

//'baz', [ 'uint32', 'bool' ]
func EventID(method string, paramTypes []string) []byte {
	var sig = method + "("
	for k, v := range paramTypes {
		if k != 0 {
			sig += ","
		}
		sig += ElementaryName(v)
	}
	sig += ")"
	return crypto.Keccak256([]byte(sig))
}

// Convert from short to canonical names
func ElementaryName (name string) string {
	if (strings.HasPrefix(name,"int[")) {
		return "int256" + name[3:]
	} else if (name == "int") {
		return "int256"
	} else if (strings.HasPrefix(name,"uint[")) {
		return "uint256" + name[4:]
	} else if (name == "uint") {
		return "uint256"
	} else if (strings.HasPrefix(name,"fixed[")) {
		return "fixed128x128" + name[5:]
	} else if (name == "fixed") {
		return "fixed128x128"
	} else if (strings.HasPrefix(name,"ufixed[")) {
		return "ufixed128x128" + name[6:]
	} else if (name == "ufixed") {
		return "ufixed128x128"
	}

	return name
}

// 'uint32'
func parseTypeN(name string) int {
	re := regexp.MustCompile(`^\D+(\d+)$`)
	matched := re.FindAllStringSubmatch(name, -1)
	N, _ := strconv.Atoi(matched[0][1])
	return N
}

// 'ufixed128x128'
func parseTypeNxM(name string) []int {
	re := regexp.MustCompile(`^\D+(\d+)x(\d+)$`)
	matched := re.FindAllStringSubmatch(name, -1)
	N, _ := strconv.Atoi(matched[0][1])
	M, _ := strconv.Atoi(matched[0][2])
	return []int{N, M}
}

// 'uint32[2]'
func parseTypeArray(name string) string {
	re := regexp.MustCompile(`(.*)\[(.*?)\]$`)
	matched := re.FindAllStringSubmatch(name, -1)
	if len(matched) != 0 {
		if matched[0][2] != "" {
			return matched[0][2]
		} else {
			return "dynamic"
		}
	}
	return ""
}

// "123456 0x123"
func parseNumber(number interface{}) *big.Int {
	switch x := number.(type) {
	case string:
		return ensureParse(x)
	case int:
		return big.NewInt(int64(x))
	case int32:
		return big.NewInt(int64(x))
	case int64:
		return big.NewInt(x)
	case uint32:
		return big.NewInt(int64(x))
	case uint64:
		return big.NewInt(int64(x))
	case *big.Int:
		return number.(*big.Int)
	default:
		panic("error parse number type")
	}
	return nil
}



// Encodes a single item (can be dynamic array)
// @returns: Buffer
func encodeSingle(paramType string, args interface{}) []byte {
	switch paramType {
	case "address":
		numStr, ok := args.(string)
		if !ok {
			numBig, ok := args.(*big.Int)
			if !ok {
				panic("parse address failed")
			}
			return encodeSingle("uint160", parseNumber(numBig))
		}
		return encodeSingle("uint160", parseNumber(numStr))
	case "bool":
		boolValue := 0
		bol, ok := args.(bool)
		if !ok {
			panic("bool invalid")
		}
		if bol {
			boolValue = 1
		}
		return encodeSingle("uint8", boolValue)
	case "string":
		str, ok := args.(string)
		if !ok {
			panic("string invalid")
		}
		return encodeSingle("bytes", []byte(str))
	case "bytes":
		arg, ok := args.([]byte)
		if !ok {
			panic("bytes invalid")
		}
		ret := append(encodeSingle("uint256", len(arg)), arg...)
		if (len(arg) % 32) != 0 {
			ret = append(ret, zeros(32-(len(arg)%32))...)
		}
		return ret
	default:
		if isArray(paramType){
			size := parseTypeArray(paramType)
			if size != "dynamic" && size != "" {
				s, err := strconv.Atoi(size)
				if err != nil {
					panic("to int error")
				}
				if reflect.TypeOf(args).Len() > s {
					panic("length error")
				}
			}
			var ret []byte
			var length int
			realType := paramType[:strings.LastIndex(paramType, "[")]
			switch arg := args.(type) {
			case []int:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []int32:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []int64:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []uint:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []uint32:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []uint64:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []string:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case [][]byte:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;
			case []*big.Int:
				length = len(arg)
				for _, v := range arg {
					ret = append(ret, encodeSingle(realType, v)...)
				}
				break;

			default:
				panic("error array args type")
			}

			if size == "dynamic" {
				var length = encodeSingle("uint256", length)
				ret = append(length, ret...)
			}
			return ret
		} else if strings.HasPrefix(paramType, "bytes") {
			size := parseTypeN(paramType)
			if size < 1 || size > 32 {
				panic("Invalid bytes<N>")
			}
			arg, ok := args.([]byte)
			if !ok {
				panic("Invalid args bytes")
			}
			return setLength(arg, 32, true)
		} else if strings.HasPrefix(paramType, "uint") {
			size := parseTypeN(paramType)
			if size % 8 != 0 || size < 8 || size > 256 {
				panic("Invalid uint<N>"  )
			}
			num := parseNumber(args)
			if num.BitLen() > size {
				panic("Supplied uint exceeds width")
			}
			if num.Cmp(big.NewInt(0)) < 0 {
				panic("Supplied uint is negative")
			}
			ret := setLength(num.Bytes(), 32, false)
			return ret
		} else if strings.HasPrefix(paramType, "int") {
			size := parseTypeN(paramType)
			if size % 8 != 0 || size < 8 || size > 256 {
				panic("Invalid int<N>"  )
			}
			num := parseNumber(args)
			if num.BitLen() > size {
				panic("Supplied uint exceeds width")
			}
			if num.Cmp(big.NewInt(0)) < 0 {
				base := big.NewInt(1)
				base.Lsh(base, 256)
				base.Sub(base, num.Abs(num))
				return base.Bytes()
			} else {
				ret := setLength(num.Bytes(), 32, false)
				return ret
			}
		} else if strings.HasPrefix(paramType, "ufixed") {
			size := parseTypeNxM(paramType)
			arg, ok := args.(string)
			if !ok {
				panic("Invalid args string")
			}
			num := parseNumber(arg)
			if num.Cmp(big.NewInt(0)) < 0 {
				panic("Supplied ufixed is negative")
			}
			base := big.NewInt(2)
			base.Exp(base, big.NewInt(int64(size[1])), nil)
			return encodeSingle("uint256", num.Mul(num, base))
		} else if strings.HasPrefix(paramType, "fixed") {
			size := parseTypeNxM(paramType)
			arg, ok := args.(string)
			if !ok {
				panic("Invalid args string")
			}
			num := parseNumber(arg)
			base := big.NewInt(2)
			base.Exp(base, big.NewInt(int64(size[1])), nil)
			return encodeSingle("int256", num.Mul(num, base))
		} else {
			panic("type error")
		}
	}
}

func isDynamic (name string) bool {
	return name == "string" || name == "bytes" || parseTypeArray(name) == "dynamic"
}

func isArray (name string) bool {
	return strings.LastIndex(name, "]") == len(name) - 1
}

func ensureParse(arg string) *big.Int {
	num := new(big.Int)
	if len(arg) < 2 {
		num.SetString(arg, 10)
	} else if arg[:2] == "0x" {
		num.SetString(arg[2:], 16)
	} else {
		num.SetString(arg, 10)
	}
	return num
}

func zeros(number int) []byte {
	return make([]byte, number)
}

func setLength(msg []byte, length int, right bool) []byte {
	buf := zeros(length)
	msgLen := len(msg)
	if right {
		if msgLen < length {
			msg = append(msg, buf[:length-msgLen]...)
			return msg
		} else {
			return msg[:length]
		}
	} else {
		if msgLen < length {
			msg = append(buf[:length-msgLen], msg...)
			return msg
		} else {
			return msg[msgLen-length:]
		}
	}
}