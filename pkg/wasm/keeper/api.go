package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/wasm/types"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"math/big"
	"strconv"
	"strings"
	"unsafe"
)

const (
	AddressSize = 20
	WASMDIR = "/wasm/"
)

type Address [AddressSize]byte

func NewAddress(raw []byte) (addr Address) {
	if len(raw) != AddressSize {
		panic("mismatch size")
	}

	copy(addr[:], raw)
	return
}

func (addr *Address) ToString() string {
	return hex.EncodeToString(addr[:])
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

func getInputLength(_ unsafe.Pointer, token int32) int32 {
	return int32(len(inputData[token]))
}

func getInput(context unsafe.Pointer, token, ptr int32, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[ptr:ptr+size], inputData[token])
}

func performSend(context unsafe.Pointer, to int32, amount int64) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	var toAddr Address
	copy(toAddr[:], memory[to:to+AddressSize])

	coinUint, err := strconv.ParseUint(strconv.FormatInt(amount, 10),10,64)
	if err != nil {
		return 0
	}

	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	fromAcc := runtimeCfg.Creator

	toAcc, err := helper.StrToAddress(toAddr.ToString())
	if err != nil {
		return 0
	}
	coin := sdk.NewUInt64Coin(coinUint)
	err = runtimeCfg.Keeper.AccountKeeper.Transfer(*runtimeCfg.Context, fromAcc, toAcc, coin)

	if err != nil {
		return 0
	}
	return 1
}

func getPreCaller(context unsafe.Pointer, callerPtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}

	precallerAddress := Address{}
	copy(precallerAddress[:], runtimeCfg.PreCaller.Bytes())

	var memory = instanceContext.Memory().Data()

	copy(memory[callerPtr:callerPtr+AddressSize], precallerAddress[:])
}

func getCreator(context unsafe.Pointer, CreatorPtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}

	creatorAddr := Address{} //contractAddress
	copy(creatorAddr[:], runtimeCfg.Creator.Bytes())

	var memory = instanceContext.Memory().Data()

	copy(memory[CreatorPtr:CreatorPtr+AddressSize], creatorAddr[:])
}

func getInvoker(context unsafe.Pointer, invokerPtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}

	invokerAddr := Address{} //contractAddress
	copy(invokerAddr[:], runtimeCfg.Invoker.Bytes())

	var memory = instanceContext.Memory().Data()

	copy(memory[invokerPtr:invokerPtr+AddressSize], invokerAddr[:])
}

func selfAddress(context unsafe.Pointer, contractPtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}

	contractAddress := Address{}
	copy(contractAddress[:], runtimeCfg.SelfAddress.Bytes())

	var memory = instanceContext.Memory().Data()

	copy(memory[contractPtr:contractPtr+AddressSize], contractAddress[:])
}

func getBlockHeader(context unsafe.Pointer, valuePtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}

	var memory = instanceContext.Memory().Data()

	var height = runtimeCfg.Context.BlockHeader().Height
	var now = runtimeCfg.Context.BlockHeader().Time.Unix()

	sink := NewSink([]byte{})
	sink.WriteU64(uint64(height))                      // 高度
	sink.WriteU64(uint64(now)) // 区块头时间

	copy(memory[valuePtr:valuePtr+8*2], sink.Bytes())
}

func notifyContract(context unsafe.Pointer, ptr, size int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}

	var memory = instanceContext.Memory().Data()

	event, err := NewEventFromSlice(memory[ptr : ptr+size])
	if err != nil {
		panic(err)
	}

	attrs := []sdk.Attribute{}
	for key, value := range event.Attr {
		attrs = append(attrs, sdk.NewAttribute(key, toString(value)))
	}
	if runtimeCfg.Context != nil {
		runtimeCfg.Context.EventManager().EmitEvent(
			sdk.NewEvent(event.Type, attrs...),
		)
	}
}

func returnContract(context unsafe.Pointer, ptr, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	result := memory[ptr : ptr+size]

	sink := NewSink(result)
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
}

func callContract(context unsafe.Pointer, addrPtr, inputPtr, inputSize int32) int32 {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()

	var addr Address
	copy(addr[:], memory[addrPtr: addrPtr + AddressSize])

	input := memory[inputPtr : inputPtr+inputSize]

	contractAddress := sdk.ToAccAddress(addr[:])
	if contractAddress == runtimeCfg.SelfAddress {
		panic(errors.New("Cannot call contract self"))
	}

	codeInfo, err := runtimeCfg.Keeper.contractInstance(*runtimeCfg.Context, contractAddress)
	if err != nil {
		panic(err)
	}
	ccstore := runtimeCfg.Context.KVStore(runtimeCfg.Keeper.storeKey)
	var code []byte
	codeHash, _ := hex.DecodeString(codeInfo.CodeHash)
	wc, err := runtimeCfg.Keeper.wasmer.GetWasmCode(runtimeCfg.Keeper.homeDir, codeHash)
	if err != nil {
		wc = ccstore.Get(codeHash)

		fileName := runtimeCfg.Keeper.wasmer.FilePathMap[fmt.Sprintf("%x",codeInfo.CodeHash)]
		err = ioutil.WriteFile(runtimeCfg.Keeper.homeDir + WASMDIR + fileName, wc, wasmtypes.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	code = wc

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
	res, err := wasmRuntime.Call(code, input, INVOKE, newRuntimeCfg)
	if err != nil {
		panic(err)
	}

	token := int32(InputDataTypeContractResult)
	inputData[token] = res
	return token
}

func newContract(context unsafe.Pointer, newContractPtr, codeHashPtr, codeHashSize, argsPtr, argsSize int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()

	args := memory[argsPtr : argsPtr + argsSize]
	codeHash := memory[codeHashPtr : codeHashPtr+codeHashSize]
	hash, err := hex.DecodeString(strings.ToLower(string(codeHash)))
	if err != nil {
		panic(err)
	}

	newContractAddress, err := runtimeCfg.Keeper.Instantiate(*runtimeCfg.Context, hash, runtimeCfg.Invoker, args, "", "", "", "", "", wasmtypes.EmptyAddress, runtimeCfg.GasWanted)
	if err != nil {
		panic(err)
	}

	contractAddress := Address{}
	copy(contractAddress[:], newContractAddress.Bytes())
	copy(memory[newContractPtr:newContractPtr+AddressSize], contractAddress[:])
}

func destroyContract(context unsafe.Pointer) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	fmt.Printf("destroy contract :%s", runtimeCfg.SelfAddress.String())

	contractAddr := runtimeCfg.SelfAddress
	ccstore := runtimeCfg.Context.KVStore(runtimeCfg.Keeper.storeKey)
	ccstore.Delete(wasmtypes.GetContractAddressKey(contractAddr))
	//contractBz := ccstore.Get(wasmtypes.GetContractAddressKey(contractAddr))
	//if contractBz == nil {
	//	panic(errors.New("get contract address failed"))
	//}
	//var contract wasmtypes.ContractInfo
	//keeper.cdc.MustUnmarshalBinaryBare(contractBz, &contract)
	//
	//codeHash, _ := hex.DecodeString(contract.CodeInfo.CodeHash)
	//
	//ccstore.Delete(wasmtypes.GetCodeKey(codeHash))


	//var wasmer Wasmer
	//wasmerBz := ccstore.Get(wasmtypes.GetWasmerKey())
	//if wasmerBz != nil {
	//	keeper.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
	//	if wasmer.LastFileID == 0 {
	//		panic(errors.New("unexpected wasmer info"))
	//	}
	//	keeper.wasmer = wasmer
	//	err := keeper.wasmer.DeleteCode(codeHash)
	//	if err != nil {
	//		panic(err)
	//	}
	//	bz := keeper.cdc.MustMarshalJSON(keeper.wasmer)
	//	ccstore.Set(wasmtypes.GetWasmerKey(), bz)
	//} else {
	//	panic(errors.New("no wasmer"))
	//}
	return
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

func panicContract(context unsafe.Pointer, dataPtr, dataSize int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	data := memory[dataPtr: dataPtr + dataSize]
	panic(string(data))
}

func debugPrint(context unsafe.Pointer, msgPtr, msgSize int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	data := memory[msgPtr : msgPtr+msgSize]
	println(string(data))
}

func getValidatorPower(context unsafe.Pointer, dataPtr, dataSize, valuePtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()

	source := NewSink(memory[dataPtr : dataPtr+dataSize])

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
	value := make([]*sdk.RustU128, len(validators))
	for _, v := range validators {
		i := 0
		val, ok := runtimeCfg.Keeper.StakingKeeper.GetValidator(*runtimeCfg.Context, sdk.HexToAddress(v.ToString()))
		if !ok {
			value[i] = sdk.NewRustU128(big.NewInt(0))
		}else {
			value[i] = sdk.NewRustU128(big.NewInt(val.DelegatorShares.TruncateInt64()))
		}
		i++
	}

	//根据链上信息返回验证者的 delegate shares
	/*value := make([]uint64, len(validators))
	for i := range value {
		value[i] = uint64(i)
	}*/

	sink := NewSink([]byte{})
	for i := range value {
		sink.WriteU128(value[i])
	}

	res := sink.Bytes()
	copy(memory[valuePtr:int(valuePtr)+len(res)], res)
}

func totalPower(context unsafe.Pointer, valuePtr int32) {
	instanceContext := wasm.IntoInstanceContext(context)
	data := instanceContext.Data()
	runtimeCfg, ok := data.(*runtimeConfig)
	if !ok {
		panic(fmt.Sprintf("%#v", data))
	}
	var memory = instanceContext.Memory().Data()
	bondedPool := runtimeCfg.Keeper.StakingKeeper.GetBondedPool(*runtimeCfg.Context)
	u128 := sdk.NewRustU128(bondedPool.GetCoin().Amount.BigInt())
	copy(memory[valuePtr:valuePtr+16], u128.Bytes())
	//根据链上信息返回总权益
	//return 123456789
}