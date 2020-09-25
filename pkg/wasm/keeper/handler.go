package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/ci123chain/ci123chain/pkg/wasm/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch tx := msg.(type) {
		case *wasm.MsgUploadContract:
			return handleMsgUploadContract(ctx, k, *tx)
		case *wasm.MsgInstantiateContract:
			return handleMsgInstantiateContract(ctx, k, *tx)
		case *wasm.MsgExecuteContract:
			return handleMsgExecuteContract(ctx, k, *tx)
		case *wasm.MsgMigrateContract:
			return handleMsgMigrateContract(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgUploadContract(ctx sdk.Context, k Keeper, msg wasm.MsgUploadContract) (res sdk.Result) {
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	codeHash, err := k.Upload(ctx, msg.Code, msg.FromAddress)
	if err != nil {
		return wasm.ErrUploadFailed(wasm.DefaultCodespace, err).Result()
	}
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", codeHash)),
		Events: ctx.EventManager().Events(),
	}
	return
}


func handleMsgInstantiateContract(ctx sdk.Context, k Keeper, msg wasm.MsgInstantiateContract) (res sdk.Result) {
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)
	nonce := ctx.Nonce()

	contractAddr, err := k.Instantiate(ctx, msg.Code, msg.FromAddress, nonce, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe, wasm.EmptyAddress)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: ctx.EventManager().Events(),
	}
	return
}

func handleMsgExecuteContract(ctx sdk.Context, k Keeper, msg wasm.MsgExecuteContract) (res sdk.Result){
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	res, err := k.Execute(ctx, msg.Contract, msg.FromAddress, msg.Args)
	if err != nil {
		return wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	}
	res.Events = ctx.EventManager().Events()
	return
}

func handleMsgMigrateContract(ctx sdk.Context, k Keeper, msg wasm.MsgMigrateContract) (res sdk.Result) {
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)
	nonce := ctx.Nonce()
	contractAddr, err := k.Migrate(ctx, msg.Code, msg.FromAddress, msg.Contract, nonce, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: ctx.EventManager().Events(),
	}
	return
}