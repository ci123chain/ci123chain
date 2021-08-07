package keeper


import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk_types "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"io/ioutil"
	"strconv"

	"github.com/wasmerio/wasmer-go/wasmer"
	"math/big"
	"strings"
)

const (
	WASMDIR = "/wasm/"
	AddressSize = 20
)
var iterToken = map[int32]*DatabaseIter{}

type DatabaseIter struct {
	Prefix string
	MockKV [][2]string
	Index  int
}

type Address [AddressSize]byte


func NewAddress(raw []byte) (addr Address) {
	if len(raw) != AddressSize {
		panic("mismatch size")
	}

	copy(addr[:], raw)
	return
}

func (addr *Address) ToString() string {
	return "0x" + hex.EncodeToString(addr[:])
}

type Event struct {
	Type string                 `json:"type"`
	Attr map[string]interface{} `json:"attr"`
}

const (
	EventAttrValueTypeInt64  = 0
	EventAttrValueTypeString = 1
)

func NewEventFromSlice(raw []byte) (Event, error) {
	event := Event{
		Attr: map[string]interface{}{},
	}

	sink := NewSink(raw)

	tp, err := sink.ReadString()
	if err != nil {
		return event, err
	}
	event.Type = tp

	sizeOfMap, err := sink.ReadU32()
	if err != nil {
		return event, err
	}

	for i := 0; i < int(sizeOfMap); i++ {
		key, err := sink.ReadString()
		if err != nil {
			return event, err
		}
		typeOfValue, err := sink.ReadByte()
		if err != nil {
			return event, err
		}

		var value interface{}
		switch typeOfValue {
		case EventAttrValueTypeInt64:
			value, err = sink.ReadI64()
		case EventAttrValueTypeString:
			value, err = sink.ReadString()
		default:
			return event, errors.New(fmt.Sprintf("unexpected event attr type: %b", typeOfValue))
		}
		if err != nil {
			return event, err
		}
		event.Attr[key] = value
	}

	return event, nil
}

func (wc *WasmerContext) getInputLength(args []wasmer.Value) ([]wasmer.Value, error) {
	token := args[0].I32()
	length := int32(len(inputData[token]))
	return []wasmer.Value{wasmer.NewI32(length)}, nil
}

func (wc *WasmerContext) getInput(args []wasmer.Value) ([]wasmer.Value, error) {
	token := args[0].I32()
	ptr := args[1].I32()
	size := args[2].I32()
	memory := wc.getMemory()

	copy(memory[ptr:ptr+size], inputData[token])

	if token > 0 {
		// 非0（0为合约输入流）token被读取后可以释放内存
		// 合法的token一定大于0
		// 分配的token不应超过int32上限，若发生则应终止合约调用并抛出异常
		delete(inputData, token)
	}
	return nil, nil
}

func (wc *WasmerContext) performSend(args []wasmer.Value) ([]wasmer.Value, error) {
	to := args[0].I32()
	amount := args[1].I64()
	memory := wc.getMemory()

	var toAddr Address
	copy(toAddr[:], memory[to:to+AddressSize])

	coinUint, err := strconv.ParseUint(strconv.FormatInt(amount, 10),10,64)
	if err != nil {
		return []wasmer.Value{wasmer.NewI32(0)}, nil
	}

	fromAcc := wc.cfg.Creator

	toAcc, err := helper.StrToAddress(toAddr.ToString())
	if err != nil {
		return []wasmer.Value{wasmer.NewI32(0)}, nil
	}
	coin := sdk_types.NewUInt64Coin(sdk_types.ChainCoinDenom, coinUint)
	err = wc.cfg.Keeper.AccountKeeper.Transfer(*wc.cfg.Context, fromAcc, toAcc, sdk_types.NewCoins(coin))
	if err != nil {
		return []wasmer.Value{wasmer.NewI32(0)}, nil
	}
	return []wasmer.Value{wasmer.NewI32(1)}, nil // 1 代表 bool true
}

func (wc *WasmerContext) getCreator(args []wasmer.Value) ([]wasmer.Value, error) {
	CreatorPtr := args[0].I32()
	memory := wc.getMemory()
	creatorAddr := Address{} //contractAddress
	copy(creatorAddr[:], wc.cfg.Creator.Bytes())
	copy(memory[CreatorPtr:CreatorPtr+AddressSize], creatorAddr[:])
	return nil, nil
}

func (wc *WasmerContext) getInvoker(args []wasmer.Value) ([]wasmer.Value, error) {
	invokerPtr := args[0].I32()
	memory := wc.getMemory()
	invokerAddr := Address{} //contractAddress
	copy(invokerAddr[:], wc.cfg.Invoker.Bytes())
	copy(memory[invokerPtr:invokerPtr+AddressSize], invokerAddr[:])
	return nil, nil
}

func (wc *WasmerContext) selfAddress(args []wasmer.Value) ([]wasmer.Value, error) {
	contractPtr := args[0].I32()
	memory := wc.getMemory()
	contractAddress := Address{}
	copy(contractAddress[:], wc.cfg.SelfAddress.Bytes())
	copy(memory[contractPtr:contractPtr+AddressSize], contractAddress[:])
	return nil, nil
}

func (wc *WasmerContext) getPreCaller(args []wasmer.Value) ([]wasmer.Value, error) {
	callerPtr := args[0].I32()
	preCallerAddress := Address{}
	copy(preCallerAddress[:], wc.cfg.PreCaller.Bytes())

	memory := wc.getMemory()
	copy(memory[callerPtr:callerPtr+AddressSize], preCallerAddress[:])
	return nil, nil
}

func (wc *WasmerContext) getBlockHeader(args []wasmer.Value) ([]wasmer.Value, error) {
	valuePtr := args[0].I32()
	memory := wc.getMemory()
	var height = wc.cfg.Context.BlockHeader().Height
	var now = wc.cfg.Context.BlockHeader().Time.Unix()
	sink := wasmtypes.NewSink([]byte{})
	sink.WriteU64(uint64(height))                      // 高度
	sink.WriteU64(uint64(now)) // 区块头时间

	copy(memory[valuePtr:valuePtr+8*2], sink.Bytes())
	return nil, nil
}

func (wc *WasmerContext) notifyContract(args []wasmer.Value) ([]wasmer.Value, error) {
	ptr := args[0].I32()
	size := args[1].I32()
	memory := wc.getMemory()

	event, err := NewEventFromSlice(memory[ptr : ptr+size])
	if err != nil {
		panic(err)
	}

	var attrs []sdk_types.Attribute
	for key, value := range event.Attr {
		attrs = append(attrs, sdk_types.NewAttribute([]byte(key), []byte(toString(value))))
	}
	if wc.cfg.Context != nil {
		wc.cfg.Context.EventManager().EmitEvent(
			sdk_types.NewEvent(event.Type, attrs...),
		)
	}
	return nil, nil
}

func (wc *WasmerContext) returnContract(args []wasmer.Value) ([]wasmer.Value, error) {
	ptr := args[0].I32()
	size := args[1].I32()
	memory := wc.getMemory()

	result := memory[ptr : ptr+size]

	sink := wasmtypes.NewSink(result)
	success, err := sink.ReadBool()
	if err != nil {
		panic(err)
	}
	msg, _, err := sink.ReadBytes()
	if err != nil {
		panic(err)
	}
	if success {
		panic(VMRes{
			res: msg,
			err: nil,
		})
	} else {
		panic(VMRes{
			res: nil,
			err: msg,
		})
	}
	return nil, nil
}

func (wc *WasmerContext) callContract(args []wasmer.Value) ([]wasmer.Value, error) {
	addrPtr := args[0].I32()
	inputPtr := args[1].I32()
	inputSize := args[2].I32()
	memory := wc.getMemory()
	runtimeCfg := wc.cfg

	var addr Address
	copy(addr[:], memory[addrPtr: addrPtr + AddressSize])

	input := memory[inputPtr : inputPtr+inputSize]

	contractAddress := sdk_types.ToAccAddress(addr[:])
	if contractAddress == wc.cfg.SelfAddress {
		panic(errors.New("Cannot call contract self"))
	}

	codeInfo, err := runtimeCfg.Keeper.contractInstance(*runtimeCfg.Context, contractAddress)
	if err != nil {
		panic(err)
	}
	ccstore := runtimeCfg.Context.KVStore(runtimeCfg.Keeper.storeKey)

	codeHash, _ := hex.DecodeString(codeInfo.CodeHash)
	wasmBC, err := runtimeCfg.Keeper.wasmer.GetWasmCode(runtimeCfg.Keeper.homeDir, codeHash)
	if err != nil {
		wasmBC = ccstore.Get(codeHash)
		fileName := runtimeCfg.Keeper.wasmer.FilePathMap[fmt.Sprintf("%x",codeInfo.CodeHash)]
		err = ioutil.WriteFile(runtimeCfg.Keeper.homeDir + WASMDIR + fileName, wasmBC, wasmtypes.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	code := wasmBC

	prefixStoreKey := wasmtypes.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(runtimeCfg.Context.KVStore(runtimeCfg.Keeper.storeKey), prefixStoreKey)

	newCreator := runtimeCfg.Keeper.GetCreator(*runtimeCfg.Context, contractAddress)
	newRuntimeCfg := &runtimeConfig{
		Store:       prefixStore,
		GasUsed:     runtimeCfg.GasUsed,
		GasWanted:   runtimeCfg.GasWanted,
		PreCaller:   runtimeCfg.SelfAddress,
		Creator:     newCreator,
		Invoker:     runtimeCfg.Invoker,
		SelfAddress: contractAddress,
		Keeper:      runtimeCfg.Keeper,
		Context:     runtimeCfg.Context,
	}
	wasmRuntime := new(wasmRuntime)
	sink := NewSink(input)
	method, err := sink.ReadString()
	if err != nil {
		panic(err)
	}
	res, err := wasmRuntime.Call(code, sink.Bytes(), method, newRuntimeCfg)
	if err != nil {
		panic(err)
	}

	token := int32(InputDataTypeContractResult)
	inputData[token] = res
	return []wasmer.Value{wasmer.NewI32(token)}, nil
}

func (wc *WasmerContext) newContract(args []wasmer.Value) ([]wasmer.Value, error) {
	codeHashPtr := args[0].I32()
	codeHashSize := args[1].I32()
	argsPtr := args[2].I32()
	argsSize := args[3].I32()
	newContractPtr := args[4].I32()
	memory := wc.getMemory()

	runtimeCfg := wc.cfg

	newArgs := memory[argsPtr : argsPtr + argsSize]
	codeHash := memory[codeHashPtr : codeHashPtr+codeHashSize]
	hash, err := hex.DecodeString(strings.ToLower(string(codeHash)))
	if err != nil {
		panic(err)
	}

	sink := NewSink(newArgs)
	method, err := sink.ReadString()
	if err != nil {
		panic(err)
	}

	input := utils.WasmInput{
		Method: method,
		Sink:   sink.Bytes(),
	}
	newContractAddress, err := runtimeCfg.Keeper.Instantiate(*runtimeCfg.Context, hash, runtimeCfg.Invoker, input, "", "", "", "", "", wasmtypes.EmptyAddress, runtimeCfg.GasWanted)
	if err != nil {
		panic(err)
	}

	contractAddress := Address{}
	copy(contractAddress[:], newContractAddress.Bytes())
	copy(memory[newContractPtr:newContractPtr+AddressSize], contractAddress[:])
	return nil, nil
}

func (wc *WasmerContext) destroyContract(args []wasmer.Value) ([]wasmer.Value, error) {
	runtimeCfg := wc.cfg
	contractAddr := runtimeCfg.SelfAddress
	ccstore := runtimeCfg.Context.KVStore(runtimeCfg.Keeper.storeKey)
	ccstore.Delete(wasmtypes.GetContractAddressKey(contractAddr))
	return nil, nil
}
// todo
func (wc *WasmerContext) panicContract(args []wasmer.Value) ([]wasmer.Value, error) {
	dataPtr := args[0].I32()
	dataSize := args[1].I32()
	memory := wc.getMemory()
	data := memory[dataPtr : dataPtr+dataSize]
	//panic("contract panic: " + string(data))
	return nil, errors.New("contract panic: " + string(data))
}

func (wc *WasmerContext) getValidatorPower(args []wasmer.Value) ([]wasmer.Value, error) {
	dataPtr := args[0].I32()
	dataSize := args[1].I32()
	valuePtr := args[2].I32()

	runtimeCfg := wc.cfg
	memory := wc.getMemory()
	source := wasmtypes.NewSink(memory[dataPtr : dataPtr+dataSize])

	var validators []Address
	{
		length, err := source.ReadU32()
		if err != nil {
			panic(err)
		}
		validators = make([]Address, 0, length)
		var i uint32 = 0
		for ; i < length; i++ {
			bytes, _, err := source.ReadBytes()
			if err != nil {
				panic(err)
			}
			validators = append(validators, NewAddress(bytes))
		}
	}
	value := make([]*wasmtypes.RustU128, len(validators))
	for _, v := range validators {
		i := 0
		val, ok := runtimeCfg.Keeper.StakingKeeper.GetValidator(*runtimeCfg.Context, sdk_types.HexToAddress(v.ToString()))
		if !ok {
			value[i] = wasmtypes.NewRustU128(big.NewInt(0))
		}else {
			value[i] = wasmtypes.NewRustU128(big.NewInt(val.DelegatorShares.TruncateInt64()))
		}
		i++
	}

	sink := wasmtypes.NewSink([]byte{})
	for i := range value {
		sink.WriteU128(value[i])
	}

	res := sink.Bytes()
	copy(memory[valuePtr:int(valuePtr)+len(res)], res)
	return nil, nil
}

func (wc *WasmerContext) totalPower(args []wasmer.Value) ([]wasmer.Value, error) {
	valuePtr := args[0].I32()
	memory := wc.getMemory()
	runtimeCfg := wc.cfg

	bondedPool := runtimeCfg.Keeper.StakingKeeper.GetBondedPool(*runtimeCfg.Context)
	u128 := wasmtypes.NewRustU128(bondedPool.GetCoins().AmountOf(sdk_types.ChainCoinDenom).BigInt())
	copy(memory[valuePtr:valuePtr+16], u128.Bytes())
	return nil, nil
}
func (wc *WasmerContext) getBalance(args []wasmer.Value) ([]wasmer.Value, error) {
	addrPtr := args[0].I32()
	balancePtr := args[1].I32()
	memory := wc.getMemory()
	runtimeCfg := wc.cfg

	var addr Address
	copy(addr[:], memory[addrPtr:addrPtr+AddressSize])

	balance := runtimeCfg.Keeper.AccountKeeper.GetBalance(*runtimeCfg.Context, sdk_types.HexToAddress(addr.ToString()))
	u128 := wasmtypes.NewRustU128(balance.AmountOf(sdk_types.ChainCoinDenom).BigInt())
	copy(memory[balancePtr:balancePtr+16], u128.Bytes())
	return nil, nil
}

func (wc *WasmerContext) debugPrint(args []wasmer.Value) ([]wasmer.Value, error) {
	msgPtr := args[0].I32()
	msgSize := args[1].I32()
	memory := wc.getMemory()

	data := memory[msgPtr : msgPtr+msgSize]
	println(string(data))
	return nil, nil
}

func (wc *WasmerContext) addGas(args []wasmer.Value) ([]wasmer.Value, error) {
	gas := args[0].I32()
	runtimeCfg := wc.cfg

	runtimeCfg.GasUsed += int64(gas)
	if uint64(runtimeCfg.GasUsed) > runtimeCfg.GasWanted {
		return nil, sdk_types.ErrorOutOfGas{Descriptor: "run vm"}
	}
	return nil, nil
}



// FOR STORE

func (wc *WasmerContext) readDB(args []wasmer.Value) ([]wasmer.Value, error) {
	// debug
	//fmt.Println(wc)

	memory := wc.getMemory()

	keyPtr := args[0].I32()
	keySize := args[1].I32()
	valuePtr := args[2].I32()
	valueSize := args[3].I32()
	offset := args[4].I32()


	realKey := memory[keyPtr: keyPtr + keySize]

	var size int
	v := wc.cfg.Store.Get(realKey)
	if v == nil {
		/*
			valueStr = ""
			size = 0;
		*/
		return []wasmer.Value{wasmer.NewI32(-1)}, nil
	} else {
		size = len(v)
	}
	if offset >= int32(size) {
		return []wasmer.Value{wasmer.NewI32(0)}, nil
	}

	index := offset + valueSize
	if index > int32(size) {
		index = int32(size)
	}

	copiedData := v[offset: index]
	copy(memory[valuePtr:valuePtr+valueSize], copiedData)

	return []wasmer.Value{wasmer.NewI32(size)}, nil
}

func (wc *WasmerContext) writeDB(args []wasmer.Value) ([]wasmer.Value, error) {
	memory := wc.getMemory()

	keyPtr := args[0].I32()
	keySize := args[1].I32()
	valuePtr := args[2].I32()
	valueSize := args[3].I32()

	realKey := memory[keyPtr: keyPtr + keySize]
	realValue := memory[valuePtr: valuePtr + valueSize]

	var Value = make([]byte, len(realValue))
	copy(Value[:], realValue[:])
	// todo currency
	wc.cfg.Store.Set(realKey, Value)
	return nil, nil
}

func (wc *WasmerContext) deleteDB(args []wasmer.Value) ([]wasmer.Value, error) {
	//func (wc *WasmerContext)  deleteDB(context unsafe.Pointer, keyPtr, keySize int32) {
	keyPtr := args[0].I32()
	keySize := args[1].I32()
	memory := wc.getMemory()

	realKey := memory[keyPtr: keyPtr + keySize]
	wc.cfg.Store.Delete(realKey)
	return nil, nil
}

func (wc *WasmerContext) newDBIter(args []wasmer.Value) ([]wasmer.Value, error) {
	//func (wc *WasmerContext)  newDBIter(context unsafe.Pointer, prefixPtr, prefixSize int32) int32 {
	prefixPtr := args[0].I32()
	prefixSize := args[1].I32()

	memory := wc.getMemory()

	realPrefix := string(memory[prefixPtr : prefixPtr+prefixSize])

	fmt.Printf("new iter prefix: [%s]\n", realPrefix)

	iter := wc.cfg.Store.parent.RemoteIterator([]byte(realPrefix), sdk_types.PrefixEndBytes([]byte(realPrefix)))

	mockKV := [][2]string{}
	for iter.Valid() {
		key := string(iter.Key())
		realKey := strings.Split(key, realPrefix)
		value := iter.Value()
		mockKV = append(mockKV, [2]string{realKey[1], string(value)})
		iter.Next()
	}

	token := int32(len(iterToken))

	// get store iterator
	iterToken[token] = &DatabaseIter{
		Prefix: realPrefix,
		Index:  -1,
		MockKV: mockKV,
	}

	return []wasmer.Value{wasmer.NewI32(token)}, nil
}

func (wc *WasmerContext) dbIterNext(args []wasmer.Value) ([]wasmer.Value, error) {
	//func (wc *WasmerContext)  dbIterNext(context unsafe.Pointer, token int32) int32 {
	token := args[0].I32()

	iter := iterToken[token]

	if iter.Index+1 >= len(iter.MockKV) {
		// next不存在
		return []wasmer.Value{wasmer.NewI32(-1)}, nil
	}

	kvToken := int32(len(inputData))

	iter.Index++
	inputData[kvToken] = []byte(iter.MockKV[iter.Index][1])

	return []wasmer.Value{wasmer.NewI32(kvToken)}, nil
}

func (wc *WasmerContext) dbIterKey(args []wasmer.Value) ([]wasmer.Value, error) {
	token := args[0].I32()
	iter := iterToken[token]

	if iter.Index < 0 || iter.Index >= len(iter.MockKV) {
		// 不存在
		return []wasmer.Value{wasmer.NewI32(-1)}, nil
	}

	kvToken := int32(len(inputData))

	inputData[kvToken] = []byte(iter.MockKV[iter.Index][0])

	return []wasmer.Value{wasmer.NewI32(kvToken)}, nil
}

func (wc *WasmerContext) dbIterValue(args []wasmer.Value) ([]wasmer.Value, error) {
	token := args[0].I32()
	iter := iterToken[token]

	if iter.Index < 0 || iter.Index >= len(iter.MockKV) {
		// 不存在
		return []wasmer.Value{wasmer.NewI32(-1)}, nil
	}

	kvToken := int32(len(inputData))

	inputData[kvToken] = []byte(iter.MockKV[iter.Index][1])

	return []wasmer.Value{wasmer.NewI32(kvToken)}, nil
}

func toString(a interface{}) string {
	if v, p := a.(int); p {
		return strconv.Itoa(v)
	}
	if v, p := a.(int16); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(int32); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(uint); p {
		return strconv.Itoa(int(v))
	}
	if v, p := a.(float32); p {
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	}
	if v, p := a.(float64); p {
		return strconv.FormatFloat(v, 'f', -1, 32)
	}
	if v, p := a.(string); p {
		return v
	}
	return ""
}
