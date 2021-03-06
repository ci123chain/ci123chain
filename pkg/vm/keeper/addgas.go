package keeper

import (
	"bytes"
	"github.com/ci123chain/wasm-util/disasm"
	"github.com/ci123chain/wasm-util/wasm"
)

func tryAddgas(raw []byte) ([]byte, error) {
	r := bytes.NewReader(raw)
	m, pos, err := wasm.DecodeModuleAddGas(r)
	if err != nil {
		return nil, err
	}
	if m.Code == nil {
		return nil, err
	}

	for i := 0; i < len(m.Code.Bodies); i++ {
		d, err := disasm.DisassembleAddGas(m.Code.Bodies[i].Code, pos)
		if err != nil {
			return nil, err
		}
		code, err := disasm.Assemble(d)
		if err != nil {
			return nil, err
		}
		m.Code.Bodies[i].Code = code
	}

	buf := new(bytes.Buffer)
	err = wasm.EncodeModule(buf, m)
	return buf.Bytes(), nil
}