package types

var WasmIdent = []byte("\x00\x61\x73\x6D")

type CallContractParam struct {
	Method string   `json:"method"`
	Args   []string `json:"args"`
}