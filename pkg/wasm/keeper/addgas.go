package keeper

import (
	"bytes"
	"errors"
	"github.com/ci123chain/wasm-util/disasm"
	"github.com/ci123chain/wasm-util/wasm"
)

const ADDGAS = "addgas"

func tryAddgas(raw []byte) ([]byte, error) {
	r := bytes.NewReader(raw)
	m, err := wasm.DecodeModule(r)
	if err != nil {
		return nil, err
	}
	if m.Code == nil {
		return nil, errors.New("decode module fail")
	}
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