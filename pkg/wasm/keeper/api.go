package keeper

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/wasm/types"
	"io/ioutil"
	"strconv"
	"time"
	"unicode/utf8"
	"unsafe"
)

type Param struct {
	Method string   `json:"method"`
	Args   []string `json:"args"`
}

func NewParamFromSlice(raw []byte) (Param, error) {
	var param Param

	sink := NewSink(raw)
	method, err := sink.ReadString()
	if err != nil {
		return param, err
	}
	param.Method = method

	size, err := sink.ReadU32()
	if err != nil {
		return param, err
	}

	for i := 0; i < int(size); i++ {
		arg, err := sink.ReadString()
		if err != nil {
			return param, err
		}
		param.Args = append(param.Args, arg)
	}

	return param, nil
}

func (param Param) Serialize() []byte {
	// 参数必须是合法的UTF8字符串
	if !utf8.ValidString(param.Method) {
		panic("invalid string")
	}
	for i := range param.Args {
		if !utf8.ValidString(param.Args[i]) {
			panic("invalid string")
		}
	}

	sink := NewSink([]byte{})
	sink.WriteString(param.Method)
	sink.WriteU32(uint32(len(param.Args)))
	for i := range param.Args {
		sink.WriteString(param.Args[i])
	}

	return sink.Bytes()
}

const AddressSize = 20

type Address [AddressSize]byte

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

func getInputLength(context unsafe.Pointer) int32 {
	return int32(len([]byte(inputData)))
}

func getInput(context unsafe.Pointer, ptr int32, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[ptr:ptr+size], inputData)
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
		return 1
	}

	fromAcc := creator

	toAcc, err := helper.StrToAddress(toAddr.ToString())
	if err != nil {
		return 1
	}
	coin := types.NewUInt64Coin(coinUint)
	err = accountKeeper.Transfer(*ctx, fromAcc, toAcc, coin)
	if err != nil {
		return 1
	}
	return 0
}

func getCreator(context unsafe.Pointer, CreatorPtr int32) {
	creatorAddr := Address{} //contractAddress
	copy(creatorAddr[:], "addr1111111111111111")

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[CreatorPtr:CreatorPtr+AddressSize], creatorAddr[:])
}

func getInvoker(context unsafe.Pointer, invokerPtr int32) {
	creatorAddr := Address{} //contractAddress
	copy(creatorAddr[:], "addr2222222222222222")

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[invokerPtr:invokerPtr+AddressSize], creatorAddr[:])
}

func getTime(context unsafe.Pointer) int64 {
	now := time.Now() //blockHeader.Time
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
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	result := memory[ptr : ptr+size]

	sink := NewSink(result)
	ok, err := sink.ReadBool()
	if err != nil {
		fmt.Println(err)
		return
	}
	length, err := sink.ReadU32()
	if err != nil {
		fmt.Println(err)
		return
	}
	msg, _, err := sink.ReadBytes(int(length))
	if err != nil {
		fmt.Println(err)
		return
	}
	if ok {
		invokeResult = fmt.Sprintf("ok msg: %s\n", string(msg))
	} else {
		invokeResult = fmt.Sprintf("error msg: %s\n", string(msg))
	}
}

func callContract(context unsafe.Pointer, addrPtr, paramPtr, paramSize int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	var addr Address
	copy(addr[:], memory[addrPtr: addrPtr + AddressSize])

	param, err := NewParamFromSlice(memory[paramPtr: paramPtr+ paramSize])
	if err != nil {
		fmt.Println(err)
		return 0
	}

	contractAddress := sdk.HexToAddress(addr.ToString())
	if contractAddress == creator {
		fmt.Println(errors.New("don't call yourself"))
		return 0
	}

	codeInfo, err := keeper.contractInstance(*ctx, contractAddress)
	if err != nil {
		fmt.Println(err)
		return 0
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
			fmt.Println(err)
			return 0
		}
	}
	code = wc

	prefixStoreKey := wasmtypes.GetContractStorePrefixKey(contractAddress)
	prefixStore := NewStore(ctx.KVStore(keeper.storeKey), prefixStoreKey)

	tempCreator := creator
	tempStore := store
	SetStore(prefixStore)
	SetCreator(contractAddress)

	err = keeper.wasmer.Call(code, param.Serialize())
	SetStore(tempStore)
	SetCreator(tempCreator)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	
	return 1
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
