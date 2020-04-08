package types

import (
	"bytes"
	"errors"
	"github.com/wasmerio/go-ext-wasm/wasmer"
)


var wasmIdent = []byte("\x00\x61\x73\x6D")

type CallContractParam struct {
	Method string   `json:"method"`
	Args   []string `json:"args"`
}

func NewCallContractParams(method string, args []string) CallContractParam {
	param := CallContractParam{
		Method: method,
		Args:   args,
	}
	return param
}



// IsWasm checks if the file contents are of wasm binary
func IsWasm(input []byte) bool {
	return bytes.Equal(input[:4], wasmIdent)
}

func IsValidaWasmFile(code []byte) error {
	if !IsWasm(code) {
		return errors.New("it is not a wasm file")
	}else {
		_, err := wasmer.Compile(code)
		if err != nil {
			return err
		}
	}
	return nil
}