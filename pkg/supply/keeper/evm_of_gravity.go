package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/supply/meta"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/pkg/errors"
	"github.com/umbracle/go-web3/abi"
	"math/big"
	"strings"
)

func (k Keeper) SendCoinsFromModuleToEVMAccount(ctx sdk.Context, to sdk.AccAddress,
	moduleName string, wlkContract sdk.AccAddress, amount *big.Int) error {

	from := k.GetModuleAddress(moduleName)


	abiIns, err := abi.NewABI(meta.DefaultERC20ABI)
	if err != nil {
		return err
	}
	m, ok := abiIns.Methods["transfer"]
	if !ok {
		return fmt.Errorf("invalid method")
	}
	params := []interface{}{to.Address, amount}
	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(from, &wlkContract.Address, data, DefaultGas, k.getNonce(ctx, from))

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Transfer action result", "value", result)
	return err
}


func (k Keeper) SendCoinsFromEVMAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
	moduleName string, wlkContract sdk.AccAddress, amount *big.Int) error {

	to := k.GetModuleAddress(moduleName)

	abiIns, err := abi.NewABI(meta.DefaultERC20ABI)
	if err != nil {
		return err
	}
	m, ok := abiIns.Methods["transfer"]
	if !ok {
		return fmt.Errorf("invalid method")
	}
	params := []interface{}{to.Address, amount}
	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(senderAddr, &wlkContract.Address, data, DefaultGas, k.getNonce(ctx, senderAddr))

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Transfer action result", "value", result)
	return err
}



func (k Keeper) WRC20DenomValueForFunc(ctx sdk.Context, moduleName string, contract sdk.AccAddress, funcName string) (value interface{}, err error) {
	sender := k.GetModuleAddress(moduleName)
	abiIns, err := abi.NewABI(meta.DefaultERC20ABI)
	if err != nil {
		return nil, err
	}
	m, ok := abiIns.Methods[funcName]
	if !ok {
		return nil, fmt.Errorf("invalid method")
	}
	data, err := abi.Encode(nil, m.Inputs)

	data = append(m.ID(), data...)

	msg := k.BuildParams(sender, &contract.Address, data, DefaultGas,  k.getNonce(ctx, sender))

	// for simulator
	ctx.WithIsCheckTx(true)
	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	if err != nil {
		return nil, err
	}
	resData, err := evmtypes.DecodeResultData(result.VMResult().Data)
	if err != nil {
		return nil, err
	}
	if len(resData.Ret) > 0 {
		respInterface, err := abi.Decode(m.Outputs, resData.Ret)
		if err != nil {
			return nil, err
		}
		resp := respInterface.(map[string]interface{})
		v, ok := resp["0"]
		if ok {
			return v, nil
		}
		return nil, err
	}

	return nil, errors.New(fmt.Sprintf("Empty result for func: %s", funcName))
}

func (k Keeper) DeployWRC20ForGivenERC20(ctx sdk.Context, moduleName string, params interface{}) (address sdk.AccAddress, err error) {
	ma := k.GetModuleAccount(ctx, moduleName)
	sender := ma.GetAddress()
	abiIns, err := abi.NewABI(meta.DefaultERC20ABI)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	data, err := abi.Encode(params, abiIns.Constructor.Inputs)
	if err != nil {
		return sdk.AccAddress{}, err
	}
	bin, err := hex.DecodeString(strings.TrimPrefix(meta.DefaultERC20Bytecode, "0x"))
	data = append(bin, data...)
	msg := k.BuildParams(sender, nil, data, DefaultGas, ma.GetSequence())

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)

	if result != nil && err == nil{
		addStr := result.VMResult().Log
		return sdk.HexToAddress(addStr), err
	}
	return sdk.AccAddress{}, err
}
