package keeper

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/go-interpreter/wagon/disasm"
	"github.com/go-interpreter/wagon/wasm"
	"github.com/go-interpreter/wagon/wasm/leb128"
	ops "github.com/go-interpreter/wagon/wasm/operators"
	"io"
	"math"
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
	pos, err := decodeAddgas(m)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(m.Code.Bodies); i++ {
		d, err := disassembleAddGas(m.Code.Bodies[i].Code, pos)
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
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeAddgas(m *wasm.Module) (int, error) {
	//判断type是否存在
	hasType := false
	var tPos int
	if m.Types != nil {
		for k,v := range m.Types.Entries{
			if len(v.ParamTypes) == 1 && v.ParamTypes[0] == wasm.ValueTypeI32 && len(v.ReturnTypes) == 0{
				hasType = true
				tPos = k
				break
			}
		}
	} else {
		m.Types = &wasm.SectionTypes{
			RawSection: wasm.RawSection{
				Start: 0,
				End:   0,
				ID:    wasm.SectionID(uint8(wasm.SectionIDType)),
				Bytes: nil,
			},
			Entries:    []wasm.FunctionSig{},
		}
		m.Sections = append([]wasm.Section{m.Types}, m.Sections...)
	}

	if !hasType {
		entry := wasm.FunctionSig{
			Form:        wasm.TypeFunc,
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: nil,
		}
		m.Types.Entries = append(m.Types.Entries, entry)
		tPos = len(m.Types.Entries) - 1
	}

	//判断sectoinImport(02) 是否存在, 不存在就实例化
	if m.Import == nil {
		m.Import = &wasm.SectionImports{
			RawSection: wasm.RawSection{
				Start: 0,
				End:   0,
				ID:    wasm.SectionID(uint8(wasm.SectionIDImport)),
				Bytes: nil,
			},
			Entries:    []wasm.ImportEntry{},
		}
		temp := append([]wasm.Section{}, m.Sections[1:]...)
		m.Sections = append(m.Sections[:1], m.Import)
		m.Sections = append(m.Sections, temp...)
	}

	//添加 addgas function
	entry := wasm.ImportEntry{
		ModuleName: "env",
		FieldName:  ADDGAS,
		Type:       wasm.FuncImport{Type:uint32(tPos)},
	}
	m.Import.Entries = append(m.Import.Entries,entry)

	mp := make(map[string]wasm.ExportEntry)
	for k,v := range m.Export.Entries {
		if v.Kind == wasm.ExternalFunction {
			v.Index++
		}
		mp[k] = v
	}
	m.Export.Entries = mp

	if m.Elements != nil {
		for i := 0; i < len(m.Elements.Entries[0].Elems); i++ {
			m.Elements.Entries[0].Elems[i]++
		}
	}

	return len(m.Import.Entries) - 1, nil
}

func disassembleAddGas(code []byte, pos int) ([]disasm.Instr, error) {
	var out []disasm.Instr
	var gasStack Stack
	var posStack Stack
	paramInstr := newInstr(0x41,"i32.const")
	//paramInstr.Immediates = append(paramInstr.Immediates, int32(0))
	gasInstr := newInstr(0x10,"call")
	gasInstr.Immediates = append(gasInstr.Immediates, uint32(pos))
	out = append(out, paramInstr, gasInstr)
	gasStack.push(0)
	posStack.push(0)

	reader := bytes.NewReader(code)
	for {
		op, err := reader.ReadByte()
		if err == io.EOF {
			gas := gasStack.pop()
			pos := posStack.pop()
			fixIm(out, pos, gas)
			break
		} else if err != nil {
			return nil, err
		}

		opStr, err := ops.New(op)
		if err != nil {
			return nil, err
		}

		instr := disasm.Instr{
			Op: opStr,
		}

		switch op {
		case ops.Block, ops.Loop, ops.If:
			sig, err := wasm.ReadByte(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, wasm.BlockType(sig))
			out = append(out, instr)
			out = append(out, paramInstr, gasInstr)
			gasStack.push(0)
			posStack.push(int32(len(out)) - 2)
		case ops.BrIf:
			depth, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, depth)
			gas := gasStack.pop()
			pos := posStack.pop()
			fixIm(out, pos, gas)
			out = append(out, instr)
			out = append(out, paramInstr, gasInstr)
			gasStack.push(0)
			posStack.push(int32(len(out)) - 2)
		case ops.Br:
			depth, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, depth)
			out = append(out, instr)
		case ops.BrTable:
			targetCount, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, targetCount)
			for i := uint32(0); i < targetCount; i++ {
				entry, err := leb128.ReadVarUint32(reader)
				if err != nil {
					return nil, err
				}
				instr.Immediates = append(instr.Immediates, entry)
			}
			defaultTarget, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, defaultTarget)
			gasStack.topAdd()
			out = append(out, instr)
		case ops.Call, ops.CallIndirect:
			index, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			if op == ops.Call && int(index) > pos {
				index ++
			}
			instr.Immediates = append(instr.Immediates, index)
			if op == ops.CallIndirect {
				idx, err := wasm.ReadByte(reader)
				if err != nil {
					return nil, err
				}
				if idx != 0x00 {
					return nil, errors.New("disasm: table index in call_indirect must be 0")
				}
				instr.Immediates = append(instr.Immediates, uint32(idx))
			}
			gasStack.topAdd()
			out = append(out, instr)
		case ops.GetLocal, ops.SetLocal, ops.TeeLocal, ops.GetGlobal, ops.SetGlobal:
			index, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, index)
			gasStack.topAdd()
			out = append(out, instr)
		case ops.I32Const:
			i, err := leb128.ReadVarint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, i)
			gasStack.topAdd()
			out = append(out, instr)
		case ops.I64Const:
			i, err := leb128.ReadVarint64(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, i)
			gasStack.topAdd()
			out = append(out, instr)
		case ops.F32Const:
			var b [4]byte
			if _, err := io.ReadFull(reader, b[:]); err != nil {
				return nil, err
			}
			i := binary.LittleEndian.Uint32(b[:])
			instr.Immediates = append(instr.Immediates, math.Float32frombits(i))
			gasStack.topAdd()
			out = append(out, instr)
		case ops.F64Const:
			var b [8]byte
			if _, err := io.ReadFull(reader, b[:]); err != nil {
				return nil, err
			}
			i := binary.LittleEndian.Uint64(b[:])
			instr.Immediates = append(instr.Immediates, math.Float64frombits(i))
			gasStack.topAdd()
			out = append(out, instr)
		case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			// read memory_immediate
			align, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, align)

			offset, err := leb128.ReadVarUint32(reader)
			if err != nil {
				return nil, err
			}
			instr.Immediates = append(instr.Immediates, offset)
			gasStack.topAdd()
			out = append(out, instr)
		case ops.CurrentMemory, ops.GrowMemory:
			idx, err := wasm.ReadByte(reader)
			if err != nil {
				return nil, err
			}
			if idx != 0x00 {
				return nil, errors.New("disasm: memory index must be 0")
			}
			instr.Immediates = append(instr.Immediates, uint8(idx))
			gasStack.topAdd()
			out = append(out, instr)
		case ops.Return:
			gas := gasStack.pop()
			pos := posStack.pop()
			fixIm(out, pos, gas)
			out = append(out, instr)
			out = append(out, paramInstr, gasInstr)
			gasStack.push(0)
			posStack.push(int32(len(out)) - 2)
		case ops.End:
			gas := gasStack.pop()
			pos := posStack.pop()
			fixIm(out, pos, gas)
			out = append(out, instr)
		default:
			gasStack.topAdd()
			out = append(out, instr)
		}
	}
	return out, nil
}

func newInstr(code byte, name string) disasm.Instr{
	ins := disasm.Instr{
		Op:          ops.Op{
			Code: code,
			Name: name,
		},
		Immediates:  []interface{}{},
		NewStack:    nil,
		Block:       nil,
		Unreachable: false,
		IsReturn:    false,
		Branches:    nil,
	}
	return ins
}

func fixIm(out []disasm.Instr, pos, im int32) {
	if len(out[pos].Immediates) > 0 {
		panic("fix error")
	}
	out[pos].Immediates = []interface{}{im}
}

type Stack struct {
	stack []int32
}

func(s *Stack) push(n int32) {
	s.stack = append(s.stack, n)
}

//pop -> fixIm
func(s *Stack) pop() (top int32) {
	lenth := len(s.stack)
	if lenth == 0 {
		panic("pop error")
	} else {
		top = s.stack[lenth-1]
		if lenth == 1 {
			s.stack = []int32{}
		} else {
			s.stack = s.stack[:lenth-1]
		}
	}
	return
}

func(s *Stack) topAdd() {
	s.stack[len(s.stack)-1]++
}