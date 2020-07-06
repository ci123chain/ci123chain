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
	"strconv"
	"unsafe"
)

const AddressSize = 20

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

	fmt.Println("send to: " + toAddr.ToString())
	fmt.Printf("send amount: %d\n", amount)

	coinUint, err := strconv.ParseUint(string(amount),10,64)
	if err != nil {
		return 0
	}

	fromAcc := creator

	toAcc, err := helper.StrToAddress(toAddr.ToString())
	if err != nil {
		return 0
	}
	coin := sdk.NewUInt64Coin(coinUint)
	err = accountKeeper.Transfer(*ctx, fromAcc, toAcc, coin)
	if err != nil {
		return 0
	}
	return 1
}

func getCreator(context unsafe.Pointer, CreatorPtr int32) {
	creatorAddr := Address{} //contractAddress
	copy(creatorAddr[:], creator.String())

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[CreatorPtr:CreatorPtr+AddressSize], creatorAddr[:])
}

func getInvoker(context unsafe.Pointer, invokerPtr int32) {
	invokerAddr := Address{} //contractAddress
	copy(invokerAddr[:], invoker.String())

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[invokerPtr:invokerPtr+AddressSize], invokerAddr[:])
}

func getTime(_ unsafe.Pointer) int64 {
	now := ctx.BlockHeader().Time //blockHeader.Time
	return now.Unix()
}

func notifyContract(context unsafe.Pointer, ptr, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	event, err := NewEventFromSlice(memory[ptr : ptr+size])
	if err != nil {
		fmt.Println(err)
	}

	attrs := []sdk.Attribute{}
	for key, value := range event.Attr {
		attrs = append(attrs, sdk.NewAttribute(key, toString(value)))
	}
	if ctx != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(event.Type, attrs...),
		)
	}
}

func returnContract(context unsafe.Pointer, ptr, size int32) {
	ResetCallResult()
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	result := memory[ptr : ptr+size]

	sink := NewSink(result)
	ok, err := sink.ReadBool()
	if err != nil {
		fmt.Println(err)
		return
	}
	msg, _, err := sink.ReadBytes()
	if err != nil {
		fmt.Println(err)
		return
	}

	callResult = msg

	if ok {
		invokeResult = fmt.Sprintf("ok msg: %s\n", string(msg))
	} else {
		invokeResult = fmt.Sprintf("error msg: %s\n", string(msg))
	}
}

func callContract(context unsafe.Pointer, addrPtr, inputPtr, inputSize int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	var addr Address
	copy(addr[:], memory[addrPtr: addrPtr + AddressSize])

	input := memory[inputPtr : inputPtr+inputSize]

	fmt.Println("call contract: " + addr.ToString())
	fmt.Print("call param: ")
	fmt.Println(input)

	contractAddress := sdk.HexToAddress(addr.ToString())
	if contractAddress == creator {
		panic(errors.New("don't call yourself"))
	}

	codeInfo, err := keeper.contractInstance(*ctx, contractAddress)
	if err != nil {
		panic(err)
	}
	ccstore := ctx.KVStore(keeper.storeKey)
	var code []byte
	codeHash, _ := hex.DecodeString(codeInfo.CodeHash)
	wc, err := keeper.wasmer.GetWasmCode(codeHash)
	if err != nil {
		wc = ccstore.Get(codeHash)

		fileName := keeper.wasmer.FilePathMap[fmt.Sprintf("%x",codeInfo.CodeHash)]
		err = ioutil.WriteFile(keeper.wasmer.HomeDir + "/" + fileName, wc, wasmtypes.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	code = wc

	prefixStoreKey := wasmtypes.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(keeper.storeKey), prefixStoreKey)

	tempCreator := creator
	tempStore := store
	SetStore(prefixStore)
	SetCreator(contractAddress)

	err = keeper.wasmer.IndirectCall(code, input)
	SetStore(tempStore)
	SetCreator(tempCreator)
	if err != nil {
		panic(err)
	}

	token := int32(InputDataTypeContractResult)
	inputData[token] = callResult
	return token
}

func destroyContract(context unsafe.Pointer) {
	fmt.Printf("destroy contract :%s", creator.String())

	contractAddr := creator
	var wasmer Wasmer
	store := ctx.KVStore(keeper.storeKey)
	contractBz := store.Get(wasmtypes.GetContractAddressKey(contractAddr))
	if contractBz == nil {
		panic(errors.New("get contract address failed"))
	}
	var contract wasmtypes.ContractInfo
	keeper.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	codeHash, _ := hex.DecodeString(contract.CodeInfo.CodeHash)
	store.Delete(wasmtypes.GetContractAddressKey(contractAddr))
	store.Delete(wasmtypes.GetCodeKey(codeHash))

	wasmerBz := store.Get(wasmtypes.GetWasmerKey())
	if wasmerBz != nil {
		keeper.cdc.MustUnmarshalJSON(wasmerBz, &wasmer)
		if wasmer.LastFileID == 0 {
			panic(errors.New("unexpected wasmer info"))
		}
		keeper.wasmer = wasmer
		err := keeper.wasmer.DeleteCode(codeHash)
		if err != nil {
			panic(err)
		}
		bz := keeper.cdc.MustMarshalJSON(keeper.wasmer)
		store.Set(wasmtypes.GetWasmerKey(), bz)
	} else {
		panic(errors.New("no wasmer"))
	}
	return
}

func migrateContract(context unsafe.Pointer, codePtr, codeSize, namePtr, nameSize, verPtr, verSize,
	authorPtr, authorSize, emailPtr, emailSize, descPtr, descSize, newAddrPtr int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	var code, name, version, author, email, desc = memory[codePtr : codePtr+codeSize],
		memory[namePtr : namePtr+nameSize],
		memory[verPtr : verPtr+verSize],
		memory[authorPtr : authorPtr+authorSize],
		memory[emailPtr : emailPtr+emailSize],
		memory[descPtr : descPtr+descSize]

	fmt.Printf("code len: %d\n", len(code))
	fmt.Printf("name: %s\n", string(name)) //实际需要判断utf8, 下同
	fmt.Printf("version: %s\n", string(version))
	fmt.Printf("author: %s\n", string(author))
	fmt.Printf("email: %s\n", string(email))
	fmt.Printf("desc: %s\n", string(desc))

	var addr = memory[newAddrPtr : newAddrPtr+AddressSize]
	copy(addr, "contract000000000002")

	return 1 // bool
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
	fmt.Printf("panic: %s\n", string(data))
	panic(string(data))
}