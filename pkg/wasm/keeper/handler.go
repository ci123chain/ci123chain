package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/ci123chain/ci123chain/pkg/wasm/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		switch tx := tx.(type) {
		case *wasm.StoreCodeTx:
			return handleStoreCodeTx(ctx, k, *tx)
		case *wasm.InstantiateContractTx:
			return handleInstantiateContractTx(ctx, k, *tx)
		case *wasm.ExecuteContractTx:
			return handleExecuteContractTx(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleStoreCodeTx(ctx sdk.Context, k Keeper, msg wasm.StoreCodeTx) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return wasm.ErrInvalidMsg(wasm.DefaultCodespace, err).Result()
	}

	codeHash, Err := k.Create(ctx, msg.Sender, msg.WASMByteCode)
	if Err != nil {
		return wasm.ErrCreateFailed(wasm.DefaultCodespace, Err).Result()
	}

	return sdk.Result{
		Data:   codeHash,
	}
}

func handleInstantiateContractTx(ctx sdk.Context, k Keeper, msg wasm.InstantiateContractTx) sdk.Result {
	contractAddr, err := k.Instantiate(ctx, msg.CodeHash, msg.Sender, msg.Args, msg.Label)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}

	return sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
	}
}

func handleExecuteContractTx(ctx sdk.Context, k Keeper, msg wasm.ExecuteContractTx) sdk.Result{
	res, err := k.Execute(ctx, msg.Contract, msg.Sender,msg.Args)
	if err != nil {
		return wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	}

	return res
}