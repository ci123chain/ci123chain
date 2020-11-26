package keeper

import (
	"encoding/hex"
	"fmt"
	"github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"path"
	"testing"
)

func TestCheckWasmFile(t *testing.T) {
	//
	fpath := "../simple.txt"
	fext := path.Ext(fpath)
	if fext != ".wasm" {
		fmt.Println("it is not wasm file")
		return
	}
	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	if !IsWasm(file) {
		fmt.Println("it is not wasm file")
		return
	}else {
		_, err := wasmer.Compile(file)
		if err != nil {
			fmt.Println("invalid wasm file")
			return
		}
		fmt.Println("valid wasm file")
		return
	}
}

func TestPanic(t *testing.T) {
	res := testPanic()
	fmt.Printf("ressss:%s\n",res)
}

func testPanic() (res string){
	res = "ee"
	ch := make(chan string)
	defer func() {
		if err := recover(); err != nil {
			ch <- "e"
		}
	}()

	go func() {
		select {
		case data := <- ch:
			fmt.Println(data)
			res = data
			break
		}
	}()
	panik()
	return "eee"
}

func panik() {
	panic("b")
}

func TestAscii(t *testing.T) {
	var c rune = 'a'
	i1 := int(c)
	fmt.Println(c)
	fmt.Println(i1)
	fmt.Println(string(c+1))
}

func TestHex(t *testing.T) {
	x, _ := hex.DecodeString("CA090A1438FB49531B7B0BFD130BD12F742211374E582895128002000000010000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000200000000000000000008000000000000000000000000000000000400000000000000000005000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000001AC5010A1438FB49531B7B0BFD130BD12F742211374E5828951220342827C97908E5E2F71151C08502A66D44B6F758E3AC2F1DE95F02EB95F0A7351220000000000000000000000000000000000000000000000000000000000000000012200000000000000000000000003F43E75AABA2C2FD6E227C10C6E7DC125A93DE3C20D6022A20E1ED7747B9DC35423B440580DB5C5090A59DB781FFE52D91CB0CAACBFF69519F3A200000000000000000000000000000000000000000000000000000000000000000400122C405608060405234801561001057600080FD5B5060043610610053576000357C010000000000000000000000000000000000000000000000000000000090048063893D20E814610058578063A6F9DAE1146100A2575B600080FD5B6100606100E6565B604051808273FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF1673FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF16815260200191505060405180910390F35B6100E4600480360360208110156100B857600080FD5B81019080803573FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF16906020019092919050505061010F565B005B60008060009054906101000A900473FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF16905090565B6000809054906101000A900473FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF1673FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF163373FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF16146101D1576040517F08C379A00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807F43616C6C6572206973206E6F74206F776E65720000000000000000000000000081525060200191505060405180910390FD5B8073FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF166000809054906101000A900473FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF1673FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF167F342827C97908E5E2F71151C08502A66D44B6F758E3AC2F1DE95F02EB95F0A73560405160405180910390A3806000806101000A81548173FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF021916908373FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF1602179055505056FEA265627A7A72315820F397F2733A89198BC7FED0764083694C5B828791F39EBCBC9E414BCCEF14B48064736F6C634300051000322A20E1ED7747B9DC35423B440580DB5C5090A59DB781FFE52D91CB0CAACBFF69519F")
	fmt.Println(string(x))
}

func TestEndKey(t *testing.T) {
	startKey := []byte("aaaaa")
	endKey := EndKey(startKey)
	fmt.Println(hex.EncodeToString(startKey))
	fmt.Println(hex.EncodeToString(endKey))
}