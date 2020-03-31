package keeper

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	wasm "github.com/tanhuiya/ci123chain/pkg/wasm/types"
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

	codeID, Err := k.Create(ctx, msg.Sender, msg.WASMByteCode, msg.Source, msg.Builder)
	if Err != nil {
		return wasm.ErrCreateFailed(wasm.DefaultCodespace, Err).Result()
	}

	//交易成功，nonce+1
	account := k.AccountKeeper.GetAccount(ctx, msg.From)
	saveErr := account.SetSequence(account.GetSequence() + 1)
	if saveErr != nil {
		return wasm.ErrSetNewAccountSequence(wasm.DefaultCodespace, saveErr).Result()
	}
	k.AccountKeeper.SetAccount(ctx, account)
	//

	return sdk.Result{
		Data:   []byte(fmt.Sprintf("codeID:%d", codeID)),
	}
}

func handleInstantiateContractTx(ctx sdk.Context, k Keeper, msg wasm.InstantiateContractTx) sdk.Result {

	contractAddr , err := k.Instantiate(ctx, msg.CodeID, msg.Sender, msg.InitMsg, msg.Label, msg.InitFunds)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}

	//交易成功，nonce+1
	account := k.AccountKeeper.GetAccount(ctx, msg.From)
	saveErr := account.SetSequence(account.GetSequence() + 1)
	if saveErr != nil {
		return wasm.ErrSetNewAccountSequence(wasm.DefaultCodespace, saveErr).Result()
	}
	k.AccountKeeper.SetAccount(ctx, account)

	return sdk.Result{
		Data:  []byte(fmt.Sprintf("contractAddress:%s", contractAddr.String())),
	}
}

func handleExecuteContractTx(ctx sdk.Context, k Keeper, msg wasm.ExecuteContractTx) sdk.Result{
	res, err := k.Execute(ctx, msg.Contract, msg.Sender,msg.Msg, msg.SendFunds)
	if err != nil {
		return wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	}

	//交易成功，nonce+1
	account := k.AccountKeeper.GetAccount(ctx, msg.From)
	saveErr := account.SetSequence(account.GetSequence() + 1)
	if saveErr != nil {
		return wasm.ErrSetNewAccountSequence(wasm.DefaultCodespace, saveErr).Result()
	}
	k.AccountKeeper.SetAccount(ctx, account)
	//
	//TODO
	return res
}