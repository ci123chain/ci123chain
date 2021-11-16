package keeper

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/supply/meta"
	"strings"
)

func (k Keeper) DeployRegistryContract(ctx sdk.Context, moduleName string, params interface{}) (sdk.AccAddress, error) {
	ma := k.GetModuleAccount(ctx, moduleName)
	defer func(account exported.ModuleAccountI) {
		if err := account.SetSequence(account.GetSequence() + 1); err != nil {
			panic(err)
		}
		k.ak.SetAccount(ctx, account)
	}(ma)
	ctx.WithIsRootMsg(true)
	sender := ma.GetAddress()

	bin, err := hex.DecodeString(strings.TrimPrefix(meta.DefaultRegistryByteCode, "0x"))

	if err != nil {
		return sdk.AccAddress{}, err
	}

	msg := k.BuildParams(sender, nil, bin, DefaultGas, ma.GetSequence())

	result, err := k.evmKeeper.EvmTxExec(ctx, msg)

	if result != nil && err == nil{
		addStr := result.VMResult().Log
		fmt.Println("contract addr: ", addStr)
		return sdk.HexToAddress(addStr), err
	}
	return sdk.AccAddress{}, err
}
