package keeper

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	wasmtypes "github.com/ci123chain/ci123chain/pkg/wasm/types"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode/utf8"
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

func returnContract(context unsafe.Pointer, ptr, size int32) {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	result := memory[ptr : ptr+size]

	sink := NewSink(result)
	_, err := sink.ReadBool()
	if err != nil {
		panic(err)
	}
	msg, _, err := sink.ReadBytes()
	if err != nil {
		panic(err)
	}

	panic(VMRes{res: msg})
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

	res, err := keeper.wasmer.Call(code, input)
	SetStore(tempStore)
	SetCreator(tempCreator)
	if err != nil {
		panic(err)
	}

	token := int32(InputDataTypeContractResult)
	inputData[token] = res
	return token
}

func destroyContract(context unsafe.Pointer) {
	fmt.Printf("destroy contract :%s", creator.String())

	contractAddr := creator
	ccstore := ctx.KVStore(keeper.storeKey)
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

func migrateContract(context unsafe.Pointer, codePtr, codeSize, namePtr, nameSize, verPtr, verSize,
	authorPtr, authorSize, emailPtr, emailSize, descPtr, descSize, initPtr, initSize, newAddrPtr int32) int32 {
	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	var code, name, version, author, email, desc, init = memory[codePtr : codePtr+codeSize],
		memory[namePtr : namePtr+nameSize],
		memory[verPtr : verPtr+verSize],
		memory[authorPtr : authorPtr+authorSize],
		memory[emailPtr : emailPtr+emailSize],
		memory[descPtr : descPtr+descSize],
		memory[initPtr : initPtr+initSize]

	if !utf8.ValidString(string(name)) {
		panic("invalid utf8 name")
	}
	if !utf8.ValidString(string(version)) {
		panic("invalid utf8 version")
	}
	if !utf8.ValidString(string(author)) {
		panic("invalid utf8 author")
	}
	if !utf8.ValidString(string(email)) {
		panic("invalid utf8 email")
	}
	if !utf8.ValidString(string(desc)) {
		panic("invalid utf8 desc")
	}
	if !utf8.ValidString(string(init)) {
		panic("invalid utf8 init")
	}

	newCodeHash, err := hex.DecodeString(string(code))
	if err != nil {
		panic(err)
	}

	oldContractAddr := creator
	ccstore := ctx.KVStore(keeper.storeKey)
	contractBz := ccstore.Get(wasmtypes.GetContractAddressKey(oldContractAddr))
	if contractBz == nil {
		panic(errors.New("get contract address failed"))
	}
	var contract wasmtypes.ContractInfo
	keeper.cdc.MustUnmarshalBinaryBare(contractBz, &contract)

	newContractAddr := keeper.generateContractAddress(newCodeHash)
	existingAcct := keeper.AccountKeeper.GetAccount(*ctx, newContractAddr)
	if existingAcct != nil {
		panic("Contract account exists")
	}

	prefix := "s/k:" + keeper.storeKey.Name() + "/"
	oldKey := wasmtypes.GetContractStorePrefixKey(oldContractAddr)

	startKey := append([]byte(prefix), oldKey...)
	endKey := EndKey(startKey)

	iter := keeper.cdb.Iterator(startKey, endKey)
	defer iter.Close()

	prefixStoreKey := wasmtypes.GetContractStorePrefixKey(newContractAddr)
	prefixStore := NewStore(ctx.KVStore(keeper.storeKey), prefixStoreKey)

	for iter.Valid() {
		key := string(iter.Key())
		realKey := strings.Split(key, string(startKey))
		value := iter.Value()
		prefixStore.Set([]byte(realKey[1]), value)
		iter.Next()
	}

	codeInfo := wasmtypes.NewCodeInfo(strings.ToUpper(hex.EncodeToString(newCodeHash)), invoker)
	ccstore.Set(wasmtypes.GetCodeKey(newCodeHash), keeper.cdc.MustMarshalBinaryBare(codeInfo))

	var initMsg json.RawMessage
	var createdAt *wasmtypes.CreatedAt
	if len(init) != 0 {
		initMsg = init
		wc, err := keeper.wasmer.GetWasmCode(newCodeHash)
		if err != nil {
			wc = store.Get(newCodeHash)
			fileName := keeper.wasmer.FilePathMap[fmt.Sprintf("%x", newCodeHash)]
			err = ioutil.WriteFile(keeper.wasmer.HomeDir + "/" + fileName, wc, wasmtypes.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		code = wc
		input, err := handleArgs(initMsg)
		if err != nil {
			panic(err)
		}
		_, err = keeper.wasmer.Call(code, input)
		if err != nil {
			panic(err)
		}
		createdAt = wasmtypes.NewCreatedAt(*ctx)
	} else {
		initMsg = contract.InitMsg
		createdAt = contract.Created
	}
	contractInfo := wasmtypes.NewContractInfo(newCodeHash, invoker, initMsg, string(name), string(version), string(author), string(email), string(desc), createdAt)
	ccstore.Set(wasmtypes.GetContractAddressKey(newContractAddr), keeper.cdc.MustMarshalBinaryBare(contractInfo))

	Account := keeper.AccountKeeper.GetAccount(*ctx, invoker)
	Account.AddContract(newContractAddr)
	keeper.AccountKeeper.SetAccount(*ctx, Account)
	//ccstore.Delete(wasmtypes.GetContractAddressKey(oldContractAddr))
	//ccstore.Delete(wasmtypes.GetCodeKey(codeHash))

	var addr = memory[newAddrPtr : newAddrPtr+AddressSize]
	copy(addr, newContractAddr.Bytes())

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

func EndKey(startKey []byte) (endKey []byte){
	key := string(startKey)
	length := len(key)
	last := []rune(key[length-1:])
	end := key[:length-1] + string(last[0] + 1)
	endKey = []byte(end)
	return
}