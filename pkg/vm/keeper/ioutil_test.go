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
	x, _ := hex.DecodeString("31303030303030")
	fmt.Println(string(x))
}

func TestEndKey(t *testing.T) {
	startKey := []byte("aaaaa")
	endKey := EndKey(startKey)
	fmt.Println(hex.EncodeToString(startKey))
	fmt.Println(hex.EncodeToString(endKey))
}