package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
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

// Implements KVStore
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
func (s Store) Has(key []byte) bool {
	res := s.parent.Get(s.key(key))
	return len(res) != 0
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
//
//func readDB(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize, offset int32) int32 {
//	instanceContext := wasm.IntoInstanceContext(context)
//	data := instanceContext.Data()
//	runtimeCfg, ok := data.(*runtimeConfig)
//	if !ok {
//		panic(fmt.Sprintf("%#v", data))
//	}
//	var memory = instanceContext.Memory().Data()
//
//	realKey := memory[keyPtr: keyPtr + keySize]
//
//	//fmt.Printf("read key [%s]\n", string(realKey))
//
//	var size int;
//	v := runtimeCfg.Store.Get(realKey)
//	if v == nil {
//		/*
//		valueStr = ""
//		size = 0;
//		*/
//		return -1
//	} else {
//		size = len(v)
//	}
//	if offset >= int32(size) {
//		return 0
//	}
//
//	index := offset + valueSize
//	if index > int32(size) {
//		index = int32(size)
//	}
//
//	copiedData := v[offset: index]
//	copy(memory[valuePtr:valuePtr+valueSize], copiedData)
//
//	return int32(size)
//}
//
//func writeDB(context unsafe.Pointer, keyPtr, keySize, valuePtr, valueSize int32) {
//	instanceContext := wasm.IntoInstanceContext(context)
//	data := instanceContext.Data()
//	runtimeCfg, ok := data.(*runtimeConfig)
//	if !ok {
//		panic(fmt.Sprintf("%#v", data))
//	}
//	var memory = instanceContext.Memory().Data()
//
//	realKey := memory[keyPtr: keyPtr + keySize]
//	realValue := memory[valuePtr: valuePtr + valueSize]
//
//	//fmt.Printf("key [%s], value [%v]\n", string(realKey), realValue)
//
//	var Value = make([]byte, len(realValue))
//	copy(Value[:], realValue[:])
//
//	runtimeCfg.Store.Set(realKey, Value)
//}
//
//func deleteDB(context unsafe.Pointer, keyPtr, keySize int32) {
//	instanceContext := wasm.IntoInstanceContext(context)
//	data := instanceContext.Data()
//	runtimeCfg, ok := data.(*runtimeConfig)
//	if !ok {
//		panic(fmt.Sprintf("%#v", data))
//	}
//	var memory = instanceContext.Memory().Data()
//
//	realKey := memory[keyPtr: keyPtr + keySize]
//
//	//fmt.Printf("delete key [%s]\n", string(realKey))
//
//	runtimeCfg.Store.Delete(realKey)
//}
//
//func newDBIter(context unsafe.Pointer, prefixPtr, prefixSize int32) int32 {
//	var instanceContext = wasm.IntoInstanceContext(context)
//	data := instanceContext.Data()
//	runtimeCfg, ok := data.(*runtimeConfig)
//	if !ok {
//		panic(fmt.Sprintf("%#v", data))
//	}
//
//	var memory = instanceContext.Memory().Data()
//	realPrefix := string(memory[prefixPtr : prefixPtr+prefixSize])
//
//	fmt.Printf("new iter prefix: [%s]\n", realPrefix)
//
//	iter := runtimeCfg.Store.parent.RemoteIterator([]byte(realPrefix), types.PrefixEndBytes([]byte(realPrefix)))
//
//	mockKV := [][2]string{}
//	for iter.Valid() {
//		key := string(iter.Key())
//		realKey := strings.Split(key, realPrefix)
//		value := iter.Value()
//		mockKV = append(mockKV, [2]string{realKey[1], string(value)})
//		iter.Next()
//	}
//
//	token := int32(len(iterToken))
//
//	// get store iterator
//	iterToken[token] = &DatabaseIter{
//		Prefix: realPrefix,
//		Index:  -1,
//		MockKV: mockKV,
//	}
//
//	return token
//}
//
//func dbIterNext(context unsafe.Pointer, token int32) int32 {
//	iter := iterToken[token]
//
//	if iter.Index+1 >= len(iter.MockKV) {
//		// next不存在
//		return -1
//	}
//
//	kvToken := int32(len(inputData))
//
//	iter.Index++
//	inputData[kvToken] = []byte(iter.MockKV[iter.Index][1])
//
//	return kvToken
//}
//
//func dbIterKey(context unsafe.Pointer, token int32) int32 {
//	iter := iterToken[token]
//
//	if iter.Index < 0 || iter.Index >= len(iter.MockKV) {
//		// 不存在
//		return -1
//	}
//
//	kvToken := int32(len(inputData))
//
//	inputData[kvToken] = []byte(iter.MockKV[iter.Index][0])
//
//	return kvToken
//}
//
//func dbIterValue(context unsafe.Pointer, token int32) int32 {
//	iter := iterToken[token]
//
//	if iter.Index < 0 || iter.Index >= len(iter.MockKV) {
//		// 不存在
//		return -1
//	}
//
//	kvToken := int32(len(inputData))
//
//	inputData[kvToken] = []byte(iter.MockKV[iter.Index][1])
//
//	return kvToken
//}