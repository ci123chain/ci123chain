package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"unsafe"
)

type Store struct {
	parent types.KVStore
	prefix []byte
}

func NewStore(parent types.KVStore, prefix []byte) Store {
	return Store{
		parent: parent,
		prefix: prefix,
	}
}

// ImplemenmiddleIns.fun["allocate"] = allocatets KVStore
func (s Store) Set(key, value []byte) {
	AssertValidKey(key)
	AssertValidValue(value)
	s.parent.Set(s.key(key), value)
	return
}

// Implements KVStore
func (s Store) Get(key []byte) []byte {
	res := s.parent.Get(s.key(key))
	return res
}

// Implements KVStore
func (s Store) Delete(key []byte) {
	s.parent.Delete(s.key(key))
}

func (s Store) key(key []byte) (res []byte) {
	if key == nil {
		panic("nil key on Store")
	}
	res = cloneAppend(s.prefix, key)
	return
}


func cloneAppend(bz []byte, tail []byte) (res []byte) {
	res = make([]byte, len(bz)+len(tail))
	copy(res, bz)
	copy(res[len(bz):], tail)
	return
}

//export read_db
func readDB(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize, offset int32) int32 {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()

	realKey := memory[keyPtr: keyPtr + keySize]

	//fmt.Printf("read key [%s]\n", string(realKey))

	var size int;

	v := runtimeCfg.Store.Get(realKey)
	if v == nil {
		/*
		valueStr = ""
		size = 0;
		*/
		return -1
	} else {
		size = len(v)
	}
	if offset >= int32(size) {
		return 0
	}

	index := offset + valueSize
	if index > int32(size) {
		index = int32(size)
	}

	copiedData := v[offset: index]
	copy(memory[valuePtr:valuePtr+valueSize], copiedData)
	return int32(size)
}

//export write_db
func writeDB(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()

	realKey := memory[keyPtr: keyPtr + keySize]
	realValue := memory[valuePtr: valuePtr + valueSize]

	//fmt.Printf("key [%s], value [%v]\n", string(realKey), realValue)

	var Value = make([]byte, len(realValue))
	copy(Value[:], realValue[:])

	runtimeCfg.Store.Set(realKey, Value)
}

//export delete_db
func deleteDB(context unsafe.Pointer, keyPtr, keySize int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()

	realKey := memory[keyPtr: keyPtr + keySize]

	fmt.Printf("delete key [%s]\n", string(realKey))

	runtimeCfg.Store.Delete(realKey)
}