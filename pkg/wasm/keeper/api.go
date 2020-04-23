package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
	"strconv"
	"unsafe"
)

func perform_transfer(context unsafe.Pointer, fromPtr int32, toPtr int32, amountPtr int32) int32 {
	var err error

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()

	fromAddr := NewRegion(memory[fromPtr : fromPtr + RegionSize])
	from := memory[fromAddr.Offset : fromAddr.Offset + fromAddr.Length]

	toAddr := NewRegion(memory[toPtr : toPtr + RegionSize])
	to := memory[toAddr.Offset : toAddr.Offset + toAddr.Length]

	amountAddr := NewRegion(memory[amountPtr : amountPtr + RegionSize])
	amount := memory[amountAddr.Offset : amountAddr.Offset + amountAddr.Length]

	coinUint, err := strconv.ParseUint(string(amount),10,64)
	if err != nil {
		return 1
	}
	fromAcc, err := helper.StrToAddress(string(from))
	if err != nil {
		return 1
	}
	toAcc, err := helper.StrToAddress(string(to))
	if err != nil {
		return 1
	}
	coin := types.NewUInt64Coin(coinUint)
	err = accountKeeper.Transfer(ctx, fromAcc, toAcc, coin)
	if err != nil {
		return 1
	}
	return 0
}

func getCreator(context unsafe.Pointer, creatorPtr int32) {
	creatorStr := creator //contractAddress
	size := len(creatorStr)

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	allocate, exist := middleIns.fun["allocate"]
	if !exist {
		panic("allocate not found")
	}
	valueOffset, err := allocate(size)
	if err != nil {
		panic(err)
	}
	copy(memory[valueOffset.ToI32():valueOffset.ToI32()+int32(size)], creatorStr)
	region := Region{
		Offset:   uint32(valueOffset.ToI32()),
		Capacity: uint32(size),
		Length:   uint32(size),
	}
	copy(memory[creatorPtr:creatorPtr+RegionSize], region.ToBytes())
}

func getInvoker(context unsafe.Pointer, invokerPtr int32) {
	invokerStr := invoker //invokerAddress
	size := len(invokerStr)

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	allocate, exist := middleIns.fun["allocate"]
	if !exist {
		panic("allocate not found")
	}
	valueOffset, err := allocate(size)
	if err != nil {
		panic(err)
	}
	copy(memory[valueOffset.ToI32():valueOffset.ToI32()+int32(size)], invokerStr)
	region := Region{
		Offset:   uint32(valueOffset.ToI32()),
		Capacity: uint32(size),
		Length:   uint32(size),
	}
	copy(memory[invokerPtr:invokerPtr+RegionSize], region.ToBytes())
}

func getTime(context unsafe.Pointer, timePtr int32) {
	tStr := blockHeader.Time.Format("2006-01-02 15:04:05")
	size := len(tStr)

	var instanceContext = wasm.IntoInstanceContext(context)
	var memory = instanceContext.Memory().Data()
	allocate, exist := middleIns.fun["allocate"]
	if !exist {
		panic("allocate not found")
	}
	valueOffset, err := allocate(size)
	if err != nil {
		panic(err)
	}
	copy(memory[valueOffset.ToI32():valueOffset.ToI32()+int32(size)], tStr)
	region := Region{
		Offset:   uint32(valueOffset.ToI32()),
		Capacity: uint32(size),
		Length:   uint32(size),
	}
	copy(memory[timePtr:timePtr+RegionSize], region.ToBytes())
}
