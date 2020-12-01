package moduletypes

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"gotest.tools/assert"
	"math/big"
	"testing"
)

func TestEncodeData(t *testing.T) {
	s, err := base64.StdEncoding.DecodeString("AAAAAAAAAAAAAAAAYUWSO06FhN6z0gIhTWhoLX2GuXw=")
	if err != nil {
		panic(err)
	}
	a := hex.EncodeToString(s)
	fmt.Println(a)
}

func TestCallData(t *testing.T) {
	a := hex.EncodeToString(utils.MethodID("baz", []string{"uint32", "bool"}))
	b := hex.EncodeToString(utils.RawEncode([]string{"uint32", "bool"}, []interface{}{uint32(69), true}))
	assert.Equal(t, "cdcd77c000000000000000000000000000000000000000000000000000000000000000450000000000000000000000000000000000000000000000000000000000000001", a+b)

	a = hex.EncodeToString(utils.MethodID("balanceOf", []string{"address"}))
	b = hex.EncodeToString(utils.RawEncode([]string{"address"}, []interface{}{"0x3f43e75aaba2c2fd6e227c10c6e7dc125a93de3c"}))
	fmt.Println(a+b)
	assert.Equal(t, "70a082310000000000000000000000003f43e75aaba2c2fd6e227c10c6e7dc125a93de3c", a+b)

	a = hex.EncodeToString(utils.MethodID("createPair", []string{"address", "address"}))
	b = hex.EncodeToString(utils.RawEncode([]string{"address", "address"}, []interface{}{"0xCAdc375A06bcBb0BAC4207fdbA2D413cAb1A6265", "0xA3c0daf03Df1527f23DD9Cfd74dcd6047707a81b"}))
	fmt.Println(a+b)
	assert.Equal(t, "c9c65396000000000000000000000000cadc375a06bcbb0bac4207fdba2d413cab1a6265000000000000000000000000a3c0daf03df1527f23dd9cfd74dcd6047707a81b", a+b)

	a = hex.EncodeToString(utils.MethodID("getPair", []string{"address", "address"}))
	b = hex.EncodeToString(utils.RawEncode([]string{"address", "address"}, []interface{}{"0xCAdc375A06bcBb0BAC4207fdbA2D413cAb1A6265", "0xA3c0daf03Df1527f23DD9Cfd74dcd6047707a81b"}))
	fmt.Println(a+b)
	assert.Equal(t, "e6a43905000000000000000000000000cadc375a06bcbb0bac4207fdba2d413cab1a6265000000000000000000000000a3c0daf03df1527f23dd9cfd74dcd6047707a81b", a+b)

	a = hex.EncodeToString(utils.MethodID("sam", []string{"bytes", "bool", "uint256[]"}))
	b = hex.EncodeToString(utils.RawEncode([]string{"bytes", "bool", "uint256[]"}, []interface{}{[]byte("dave"), true, []int64{1, 2, 3}}))
	assert.Equal(t, "a5643bf20000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000464617665000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003", a+b)

	a = hex.EncodeToString(utils.MethodID("f", []string{"uint", "uint32[]", "bytes10", "bytes"}))
	b = hex.EncodeToString(utils.RawEncode([]string{"uint", "uint32[]", "bytes10", "bytes"}, []interface{}{"0x123", []string{"0x456", "0x789"}, []byte("1234567890"), []byte("Hello, world!")}))
	assert.Equal(t, "8be6524600000000000000000000000000000000000000000000000000000000000001230000000000000000000000000000000000000000000000000000000000000080313233343536373839300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000004560000000000000000000000000000000000000000000000000000000000000789000000000000000000000000000000000000000000000000000000000000000d48656c6c6f2c20776f726c642100000000000000000000000000000000000000", a+b)
}

func TestMethodSig(t *testing.T) {
	a := hex.EncodeToString(utils.MethodID("test", nil))
	assert.Equal(t, "f8a8fd6d", a)

	a = hex.EncodeToString(utils.MethodID("test", []string{"uint"}))
	assert.Equal(t, "29e99f07", a)

	a = hex.EncodeToString(utils.MethodID("test", []string{"uint256"}))
	assert.Equal(t, "29e99f07", a)

	a = hex.EncodeToString(utils.MethodID("test", []string{"uint", "uint"}))
	assert.Equal(t, "eb8ac921", a)

	a = hex.EncodeToString(utils.MethodID("allPairs", nil))
	fmt.Println(a)
}

func TestRawEncode(t *testing.T) {
	//encoding negative int32
	a := hex.EncodeToString(utils.RawEncode([]string{"int32"}, []interface{}{-2}))
	assert.Equal(t, "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", a)

	//encoding negative int256
	bigNum := new(big.Int)
	bigNum.SetString("-19999999999999999999999999999999999999999999999999999999999999", 10)
	a = hex.EncodeToString(utils.RawEncode([]string{"int256"}, []interface{}{bigNum}))
	assert.Equal(t, "fffffffffffff38dd0f10627f5529bdb2c52d4846810af0ac000000000000001", a)

	//encoding string >32bytes
	a = hex.EncodeToString(utils.RawEncode([]string{"string"}, []interface{}{" hello world hello world hello world hello world  hello world hello world hello world hello world  hello world hello world hello world hello world hello world hello world hello world hello world"}))
	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c22068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64202068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64202068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c642068656c6c6f20776f726c64000000000000000000000000000000000000000000000000000000000000", a)

	//encoding uint32 response
	a = hex.EncodeToString(utils.RawEncode([]string{"uint32"}, []interface{}{42}))
	assert.Equal(t, "000000000000000000000000000000000000000000000000000000000000002a", a)

	//encoding string response (unsupported)
	a = hex.EncodeToString(utils.RawEncode([]string{"string"}, []interface{}{"a response string (unsupported)"}))
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000001f6120726573706f6e736520737472696e672028756e737570706f727465642900", a)

}