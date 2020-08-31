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

var store Store
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

//set the store that be used by rust contract.
func SetStore(kvStore Store) {
	store = kvStore
}

//export read_db
func readDB(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize, offset int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	realKey := memory[keyPtr: keyPtr + keySize]

	fmt.Printf("read key [%s]\n", string(realKey))

	var size int;
	var valueStr string;
	v := store.Get(realKey)
	if v == nil {
		/*
		valueStr = ""
		size = 0;
		*/
		return -1
	} else {
		valueStr = string(v)
		size = len(valueStr)
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
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	realKey := memory[keyPtr: keyPtr + keySize]
	realValue := memory[valuePtr: valuePtr + valueSize]

	fmt.Printf("write key [%s], value [%s]\n", string(realKey), string(realValue))

	valueStr := string(realValue)
	Value := []byte(valueStr)

	store.Set(realKey, Value)
}

//export delete_db
func deleteDB(context unsafe.Pointer, keyPtr, keySize int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	realKey := memory[keyPtr: keyPtr + keySize]

	fmt.Printf("delete key [%s]\n", string(realKey))

	store.Delete(realKey)
}