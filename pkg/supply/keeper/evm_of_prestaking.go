package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ci123chain/ci123chain/pkg/supply/meta"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ethereum/go-ethereum/common"
	web3 "github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"math"
	"math/big"
)

///supply/keeper/evm_of_prestaking.go
const DefaultGas = math.MaxUint64 / 2


//// MintCoinsFromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
//func (k Keeper) MintCoinsFromModuleToEvmAccount(ctx sdk.Context,
//	recipientAddr sdk.AccAddress, wlkContract string, amt *big.Int) error {
//	//param := []interface{}{metaData.Name, metaData.Symbol, metaData.Symbol, 0, true}
//	//denomAddr, err := a.DeployERC20Contract(ctx, owner, param)
//	//if err != nil {
//	//	return err
//	//}
//	err := k.Mint(ctx, sdk.HexToAddress(wlkContract), recipientAddr, types.ModuleName, amt)
//	return err
//}

// TransferFromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) TransferFromModuleToEvmAccount(ctx sdk.Context,
	recipientAddr sdk.AccAddress, wlkContract string, amt *big.Int) error {
	return k.SendCoinsFromModuleToEVMAccount(ctx, recipientAddr, types.ModuleName, sdk.HexToAddress(wlkContract), amt)
}

func (k Keeper) BuildParams(sender sdk.AccAddress, to *common.Address, payload []byte, gasLimit, nonce uint64) evmtypes.MsgEvmTx {
	return evmtypes.MsgEvmTx{
		From: sender,
		Data: evmtypes.TxData{
			Payload: payload,
			Amount: big.NewInt(0),
			Recipient: to,
			GasLimit: gasLimit,
			AccountNonce: nonce,
		},
	}
}
//
//func (k Keeper) QueryTokenManagerAddress (ctx sdk.Context, weelinkDao sdk.AccAddress, moduleName string, params interface{}) (sdk.AccAddress, error) {
//	sender := k.GetModuleAddress(moduleName)
//	abiIns, err := abi.NewABI(meta.DefaultDaoABI)
//	if err != nil {
//		return sdk.AccAddress{}, err
//	}
//	m, ok := abiIns.Methods["tokenManager"]
//	if !ok {
//		return sdk.AccAddress{}, fmt.Errorf("invalid method")
//	}
//	data, err := abi.Encode(nil, m.Inputs)
//
//	data = append(m.ID(), data...)
//
//	msg := k.BuildParams(sender, &weelinkDao.Address, data, DefaultGas,  k.getNonce(ctx, sender))
//
//	// for simulator
//	ctx.WithIsCheckTx(true)
//	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
//	if err != nil {
//		return sdk.AccAddress{}, err
//	}
//	resData, err := evmtypes.DecodeResultData(result.VMResult().Data)
//	if err != nil {
//		return sdk.AccAddress{}, err
//	}
//	if len(resData.Ret) > 0 {
//		respInterface, err := abi.Decode(m.Outputs, resData.Ret)
//		if err != nil {
//			return sdk.AccAddress{}, err
//		}
//		resp := respInterface.(map[string]interface{})
//		v, ok := resp["0"].(web3.Address)
//		if ok {
//			return sdk.HexToAddress(v.String()), nil
//		}
//		return sdk.AccAddress{}, err
//	}
//	return sdk.AccAddress{}, err
//}
//
//func (k Keeper) DeployDaoContract(ctx sdk.Context, moduleName string, params interface{}) (sdk.AccAddress, error) {
//	ma := k.GetModuleAccount(ctx, moduleName)
//	defer func(account exported.ModuleAccountI) {
//		if err := account.SetSequence(account.GetSequence() + 1); err != nil {
//			panic(err)
//		}
//		k.ak.SetAccount(ctx, account)
//	}(ma)
//	ctx.WithIsRootMsg(true)
//	sender := ma.GetAddress()
//	abiIns, err := abi.NewABI(meta.DefaultDaoABI)
//
//	if err != nil {
//		return sdk.AccAddress{}, err
//	}
//
//	bin, err := hex.DecodeString(strings.TrimPrefix(DefaultDaoByteCode, "0x"))
//
//	if err != nil {
//		return sdk.AccAddress{}, err
//	}
//
//	data, err := abi.Encode(params, abiIns.Constructor.Inputs)
//	if err != nil {
//		return sdk.AccAddress{}, err
//	}
//
//	data = append(bin, data...)
//	msg := k.BuildParams(sender, nil, data, DefaultGas, ma.GetSequence())
//
//	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
//
//	if result != nil && err == nil{
//		addStr := result.VMResult().Log
//		return sdk.HexToAddress(addStr), err
//	}
//	return sdk.AccAddress{}, err
//}


func (k Keeper) BurnEVMCoin(ctx sdk.Context, moduleName string, wlkContract, to sdk.AccAddress, amount *big.Int) error {
	from := k.GetModuleAddress(moduleName)

	abiIns, err := abi.NewABI(meta.DefaultTokenManagerABI)
	if err != nil {
		return err
	}
	m, ok := abiIns.Methods["burn"]
	if !ok {
		return fmt.Errorf("invalid method")
	}
	params := []interface{}{web3.HexToAddress(to.Address.String()), amount}
	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(from, &wlkContract.Address, data, DefaultGas, k.getNonce(ctx, from))
	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Burn action result", "value", result)
	return err
}

func (k Keeper) Mint(ctx sdk.Context, contract, to sdk.AccAddress, moduleName string, amount *big.Int) error {
	sender := k.GetModuleAddress(moduleName)
	fmt.Println(sender.String())
	abiIns, err := abi.NewABI(meta.DefaultTokenManagerABI)
	if err != nil {
		return err
	}

	m, ok := abiIns.Methods["mint"]
	if !ok {
		return fmt.Errorf("invalid method")
	}

	params := []interface{}{web3.HexToAddress(to.Address.String()), amount}

	data, err := abi.Encode(params, m.Inputs)
	data = append(m.ID(), data...)

	ctx = ctx.WithIsRootMsg(true)
	msg := k.BuildParams(sender, &contract.Address, data, DefaultGas,  k.getNonce(ctx, sender))

	_, err = k.evmKeeper.EvmTxExec(ctx, msg)
	if err != nil {
		k.Logger(ctx).Error("Mint action result", "error", err.Error())
	}
	return err
}


func (k Keeper)getNonce(ctx sdk.Context, address sdk.AccAddress) uint64 {
	return k.ak.GetAccount(ctx, address).GetSequence()
}
