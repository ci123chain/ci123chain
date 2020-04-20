package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"unsafe"
)


//要在import函数中调用export函数，需要中转函数
type middle struct {
	fun map[string]func(...interface{}) (wasm.Value, error)
}

const RegionSize = 12

type Store struct {
	parent types.KVStore
	prefix []byte
}

var store Store
var middleIns = middle{fun: make(map[string]func(...interface{}) (wasm.Value, error))}
func NewStore(parent types.KVStore, prefix []byte) Store {
	return Store{
		parent: parent,
		prefix: prefix,
	}
}

// Implements KVStore
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
func readDB(context unsafe.Pointer, key, value int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	keyAddr := NewRegion(memory[key : key+RegionSize])
	realKey := memory[keyAddr.Offset : keyAddr.Offset+keyAddr.Length]

	fmt.Printf("read key [%s]\n", string(realKey))

	allocate, exist := middleIns.fun["allocate"]
	if !exist {
		panic("allocate not found")
	}

	/*if _, exist := store[string(realKey)]; !exist {
		panic(string(realKey) + " not found")
	}

	size := len(store[string(realKey)])*/
	var size int;
	var valueStr string;
	v := store.Get(realKey)
	if v == nil {
		valueStr = ""
		size = 0;
	} else {
		valueStr = string(v)
		size = len(valueStr)
	}

	valueOffset, err := allocate(size)
	if err != nil {
		panic(err)
	}
	//copy(memory[valueOffset.ToI32():valueOffset.ToI32()+int32(size)], store[string(realKey)])
	copy(memory[valueOffset.ToI32():valueOffset.ToI32()+int32(size)], valueStr)

	region := Region{
		Offset:   uint32(valueOffset.ToI32()),
		Capacity: uint32(size),
		Length:   uint32(size),
	}
	copy(memory[value:value+RegionSize], region.ToBytes())

	return 0
}

//export write_db
func writeDB(context unsafe.Pointer, key, value int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	keyAddr := NewRegion(memory[key : key+RegionSize])
	realKey := memory[keyAddr.Offset : keyAddr.Offset+keyAddr.Length]
	valueAddr := NewRegion(memory[value : value+RegionSize])
	realValue := memory[valueAddr.Offset : valueAddr.Offset+valueAddr.Length]

	fmt.Printf("write key [%s], value [%s]\n", string(realKey), string(realValue))

	//store[string(realKey)] = string(realValue)
	valueStr := string(realValue)
	Value := []byte(valueStr)

	store.Set(realKey, Value)
}

//export delete_db
func deleteDB(context unsafe.Pointer, key int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	keyAddr := NewRegion(memory[key : key+RegionSize])
	realKey := memory[keyAddr.Offset : keyAddr.Offset+keyAddr.Length]

	fmt.Printf("delete key [%s]\n", string(realKey))

	//delete(store, string(realKey))

	store.Delete(realKey)
}

//Region 内存指针
type Region struct {
	Offset   uint32
	Capacity uint32
	Length   uint32
}

func NewRegion(b []byte) Region {
	var ret Region
	bytesBuffer := bytes.NewBuffer(b)
	_ = binary.Read(bytesBuffer, binary.LittleEndian, &ret.Offset)
	_ = binary.Read(bytesBuffer, binary.LittleEndian, &ret.Capacity)
	_ = binary.Read(bytesBuffer, binary.LittleEndian, &ret.Length)
	return ret
}

func (region Region) ToBytes() []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.LittleEndian, region.Offset)
	_ = binary.Write(bytesBuffer, binary.LittleEndian, region.Capacity)
	_ = binary.Write(bytesBuffer, binary.LittleEndian, region.Length)
	return bytesBuffer.Bytes()
}