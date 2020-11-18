package rest

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"gotest.tools/assert"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type signature struct {
	Method 	string
	Args 	[]string
	Retargs []string
}

func TestCallData(t *testing.T) {
	a := hex.EncodeToString(methodID("baz", []string{"uint32", "bool"}))
	b := hex.EncodeToString(rawEncode([]string{"uint32", "bool"}, []interface{}{"69", true}))
	assert.Equal(t, "cdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001", a+b)

	a = hex.EncodeToString(methodID("sam", []string{"bytes", "bool", "uint256[]"}))
	b = hex.EncodeToString(rawEncode([]string{"bytes", "bool", "uint256[]"}, []interface{}{[]byte("dave"), true, []int64{1, 2, 3}}))
	assert.Equal(t, "a5643bf20000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000464617665000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003", a+b)

	a = hex.EncodeToString(methodID("f", []string{"uint", "uint32[]", "bytes10", "bytes"}))
	b = hex.EncodeToString(rawEncode([]string{"uint", "uint32[]", "bytes10", "bytes"}, []interface{}{"0x123", []string{"0x456", "0x789"}, []byte("1234567890"), []byte("Hello, world!")}))
	assert.Equal(t, "8be6524600000000000000000000000000000000000000000000000000000000000001230000000000000000000000000000000000000000000000000000000000000080313233343536373839300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000004560000000000000000000000000000000000000000000000000000000000000789000000000000000000000000000000000000000000000000000000000000000d48656c6c6f2c20776f726c642100000000000000000000000000000000000000", a+b)
}

func TestMethodSig(t *testing.T) {
	a := hex.EncodeToString(methodID("test", nil))
	assert.Equal(t, "f8a8fd6d", a)

	a = hex.EncodeToString(methodID("test", []string{"uint"}))
	assert.Equal(t, "29e99f07", a)

	a = hex.EncodeToString(methodID("test", []string{"uint256"}))
	assert.Equal(t, "29e99f07", a)

	a = hex.EncodeToString(methodID("test", []string{"uint", "uint"}))
	assert.Equal(t, "eb8ac921", a)
}

func TestRawEncode(t *testing.T) {
	//encoding negative int32
	a := hex.EncodeToString(rawEncode([]string{"int32"}, []interface{}{-2}))
	assert.Equal(t, "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", a)

	//encoding negative int256
	bigNum := new(big.Int)
	bigNum.SetString("-19999999999999999999999999999999999999999999999999999999999999", 10)
	a = hex.EncodeToString(rawEncode([]string{"int256"}, []interface{}{bigNum}))
	assert.Equal(t, "fffffffffffff38dd0f10627f5529bdb2c52d4846810af0ac000000000000001", a)

	//encoding string >32bytes
	a = hex.EncodeToString(rawEncode([]string{"string"}, []interface{}{" hello world hello world hello world hello world  hello world hello world hello world hello world  hello world hello world hello world hello world hello world hello world hello world hello world"}))
	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c22068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64202068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64202068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64000000000000000000000000000000000000000000000000000000000000", a)

	//encoding uint32 response
	a = hex.EncodeToString(rawEncode([]string{"uint32"}, []interface{}{42}))
	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000002a", a)

	//encoding string response (unsupported)
	a = hex.EncodeToString(rawEncode([]string{"string"}, []interface{}{"a response string (unsupported)"}))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001f6120726573706f6e736520737472696e672028756e737570706f727465642900", a)
}

//[ 'uint32', 'bool' ], [ 69, 1 ]
func rawEncode(paramTypes []string, values []interface{}) []byte {
	var output, data []byte
	headLength := 0
	for _, v := range paramTypes{
		if isArray(v){
			var size = parseTypeArray(v)
			if size != "dynamic" {
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
		param := elementaryName(paramTypes[i])
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

func simpleEncode(method string, params string) string {
	return ""
}

//'baz', [ 'uint32', 'bool' ]
func methodID(method string, paramTypes []string) []byte {
	return eventID(method, paramTypes)[:4]
}

//'baz', [ 'uint32', 'bool' ]
func eventID(method string, paramTypes []string) []byte {
	var sig = method + "("
	for k, v := range paramTypes {
		if k != 0 {
			sig += ","
		}
		sig += elementaryName(v)
	}
	sig += ")"
	return crypto.Keccak256([]byte(sig))
}

// Convert from short to canonical names
func elementaryName (name string) string {
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
	switch number.(type) {
		case string:
		 	return ensureParse(number.(string))
		case int:
			return big.NewInt(int64(number.(int)))
		case int64:
			return big.NewInt(number.(int64))
		case float64:
			return big.NewInt(number.(int64))
		case *big.Int:
			return number.(*big.Int)
		default:
		panic("error parse number type")
	}
	return nil
}

// someMethod(bytes,uint)
// someMethod(bytes,uint):(boolean)
func parseSignature(sig string) signature {
	re := regexp.MustCompile(`^(\w+)\((.*)\)$`)
	matchedSig := re.FindAllStringSubmatch(sig, -1)
	if len(matchedSig) != 0 {
		if len(matchedSig[0]) != 4 {
			panic("Invalid method signature")
		}
	} else {
		panic("Invalid method signature")
	}

	re = regexp.MustCompile(`^(.+)\)$`)
	matchedRet := re.FindAllStringSubmatch(matchedSig[0][3], -1)
	if len(matchedRet) != 0 { //ret
		return signature{
			Method: matchedSig[0][1],
			Args: strings.Split(matchedRet[0][1], ","),
			Retargs: strings.Split(matchedRet[0][2], ","),
		}
	} else { // no ret
		return signature {
			Method: matchedSig[0][1],
			Args:   strings.Split(matchedSig[0][2], ","),
		}
	}
}

// Encodes a single item (can be dynamic array)
// @returns: Buffer
func encodeSingle(paramType string, args interface{}) []byte {
	switch paramType {
	case "address":
		numStr, ok := args.(string)
		if !ok {
			panic("number invalid")
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
					panic("int")
				}
				if reflect.TypeOf(args).Len() > s {
					panic("length")
				}
			}
			var ret []byte
			var length int
			realType := paramType[:strings.LastIndex(paramType, "[")]
			//todo fix?
			switch args.(type) {
				case []string:
					length = len(args.([]string))
					for _, v := range args.([]string) {
						ret = append(ret, encodeSingle(realType, v)...)
					}
					break;
				case []int:
					length = len(args.([]int))
					for _, v := range args.([]int) {
						ret = append(ret, encodeSingle(realType, v)...)
					}
					break;
				case []int64:
					length = len(args.([]int64))
					for _, v := range args.([]int64) {
						ret = append(ret, encodeSingle(realType, v)...)
					}
					break;
				case []float64:
					length = len(args.([]int64))
					for _, v := range args.([]int64) {
						ret = append(ret, encodeSingle(realType, v)...)
					}
					break;
				case []*big.Int:
					length = len(args.([]*big.Int))
					for _, v := range args.([]*big.Int) {
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
			msg = append(msg, buf[:length - msgLen]...)
			return msg
		} else {
			return msg[:length]
		}
	} else {
		if msgLen < length {
			msg = append(buf[:length - msgLen], msg...)
			return msg
		} else {
			return msg[msgLen - length:]
		}
	}
}