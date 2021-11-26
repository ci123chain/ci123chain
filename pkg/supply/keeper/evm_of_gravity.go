package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/supply/meta"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/pkg/errors"
	"github.com/umbracle/go-web3/abi"
	"math/big"
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
