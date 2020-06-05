package keeper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"strconv"
	"unsafe"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

const AddressSize = 20
var inputData []byte
type Address [AddressSize]byte

func (addr *Address) ToString() string {
	return hex.EncodeToString(addr[:])
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
	copy(toAddr[:], memory[to: to + AddressSize])

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
	copy(creatorAddr[:], creator.String())

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[CreatorPtr: CreatorPtr + AddressSize], creatorAddr[:])
}

func getInvoker(context unsafe.Pointer, invokerPtr int32) {
	invokerAddr := Address{}//contractAddress
	copy(invokerAddr[:], invoker.String())

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	copy(memory[invokerPtr: invokerPtr + AddressSize], invokerAddr[:])
}

func getTime(context unsafe.Pointer) int64 {
	blockTime := blockHeader.Time //blockHeader.Time
	return blockTime.Unix()
}

func notifyContract(context unsafe.Pointer, ptr, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	type Event struct {
		Type string                 `json:"type"`
		Attr map[string]interface{} `json:"attr"`
	}

	var event Event
	err := json.Unmarshal(memory[ptr: ptr + size], &event)
	if err != nil {
		panic(err)
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

var InvokeResult string
func returnContract(context unsafe.Pointer, ptr, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	result := memory[ptr: ptr + size]

	var resp RespW
	err := json.Unmarshal(result, &resp)
	if err != nil {
		fmt.Println(err)
	}
	if resp.Err != "" {
		InvokeResult = resp.Err
	}else{
		InvokeResult = string(resp.Ok.Data)
	}

}

type RespW struct {
	Ok  RespN   `json:"ok"`
	Err string 	`json:"err"`
}

type RespN struct {
	Data []byte `json:"data"`
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