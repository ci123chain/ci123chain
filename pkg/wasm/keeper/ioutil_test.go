package keeper

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/wasm/types"
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
	if !types.IsWasm(file) {
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
