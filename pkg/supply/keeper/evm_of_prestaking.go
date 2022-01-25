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
