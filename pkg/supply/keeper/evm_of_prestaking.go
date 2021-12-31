package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ci123chain/ci123chain/pkg/supply/meta"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/abi"
	"math"
	"math/big"
	"strings"
)

///supply/keeper/evm_of_prestaking.go
const (
	DefaultGas = math.MaxUint64 / 2
)

// MintCoinsFromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) MintCoinsFromModuleToEvmAccount(ctx sdk.Context,
	recipientAddr sdk.AccAddress, wlkContract string, amt *big.Int) error {
	//param := []interface{}{metaData.Name, metaData.Symbol, metaData.Symbol, 0, true}
	//denomAddr, err := a.DeployERC20Contract(ctx, owner, param)
	//if err != nil {
	//	return err
	//}
	err := k.Mint(ctx, sdk.HexToAddress(wlkContract), recipientAddr, types.ModuleName, amt)
	return err
}

// TransferFromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) TransferFromModuleToEvmAccount(ctx sdk.Context,
	recipientAddr sdk.AccAddress, wlkContract string, amt *big.Int) error {
	return k.SendCoinsFromModuleToEVMAccount(ctx, recipientAddr, types.ModuleName, sdk.HexToAddress(wlkContract), amt)
}

// Transfer721FromModuleToEvmAccount transfers coins from a ModuleAccount to an AccAddress
func (k Keeper) Transfer721FromModuleToEvmAccount(ctx sdk.Context,
	recipientAddr sdk.AccAddress, wlkContract string, tokenId *big.Int) error {
	return k.Send721CoinsFromModuleToEVMAccount(ctx, recipientAddr, types.ModuleName, sdk.HexToAddress(wlkContract), tokenId)
}

func (k Keeper) BuildParams(sender sdk.AccAddress, to *common.Address, payload []byte, gasLimit, nonce uint64) evmtypes.MsgEvmTx {
	return evmtypes.MsgEvmTx{
		From: sender,
		Data: evmtypes.TxData{
			Payload:      payload,
			Amount:       big.NewInt(0),
			Recipient:    to,
			GasLimit:     gasLimit,
			AccountNonce: nonce,
		},
	}
}

func (k Keeper) DeployDaoContract(ctx sdk.Context, moduleName string, params interface{}) (sdk.AccAddress, error) {
	ma := k.GetModuleAccount(ctx, moduleName)
	defer func(account exported.ModuleAccountI) {
		if err := account.SetSequence(account.GetSequence() + 1); err != nil {
			panic(err)
		}
		k.ak.SetAccount(ctx, account)
	}(ma)
	ctx.WithIsRootMsg(true)
	sender := ma.GetAddress()
	abiIns, err := abi.NewABI(meta.DefaultDaoABI)

	if err != nil {
		return sdk.AccAddress{}, err
	}

	bin, err := hex.DecodeString(strings.TrimPrefix(DefaultDaoByteCode, "0x"))

	if err != nil {
		return sdk.AccAddress{}, err
	}

	data, err := abi.Encode(params, abiIns.Constructor.Inputs)
	if err != nil {
		return sdk.AccAddress{}, err
	}

	data = append(bin, data...)
	msg := k.BuildParams(sender, nil, data, DefaultGas, ma.GetSequence())

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)

	if result != nil && err == nil {
		addStr := result.VMResult().Log
		return sdk.HexToAddress(addStr), err
	}
	return sdk.AccAddress{}, err
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
	msg := k.BuildParams(sender, &contract.Address, data, DefaultGas, k.getNonce(ctx, sender))

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)
	k.Logger(ctx).Info("Mint action result", "value", result)
	return err
}

func (k Keeper) getNonce(ctx sdk.Context, address sdk.AccAddress) uint64 {
	return k.ak.GetAccount(ctx, address).GetSequence()
}
