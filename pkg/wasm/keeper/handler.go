package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
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
	em := ctx.EventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(wasm.EventType,
			sdk.NewAttribute(wasm.AttributeKeyMethod, wasm.AttributeValueUpload),
			sdk.NewAttribute(wasm.AttributeKeyCodeHash, string(codeHash)),
			sdk.NewAttribute(sdk.AttributeKeyModule, wasm.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	})
	res = sdk.Result{
		Data:  codeHash,
		Events: em.Events(),
	}
	return
}


func handleMsgInstantiateContract(ctx sdk.Context, k Keeper, msg wasm.MsgInstantiateContract) (res sdk.Result) {
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	contractAddr, err := k.Instantiate(ctx, msg.CodeHash, msg.FromAddress, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe, wasm.EmptyAddress)

	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	em := ctx.EventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(wasm.EventType,
			sdk.NewAttribute(wasm.AttributeKeyMethod, wasm.AttributeValueInitiate),
			sdk.NewAttribute(wasm.AttributeKeyAddress, contractAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, wasm.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	})
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: em.Events(),
	}
	return
}

func handleMsgExecuteContract(ctx sdk.Context, k Keeper, msg wasm.MsgExecuteContract) (res sdk.Result){
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	logger.GetLogger().Info(("!!!!!!!!!!!!!Enter HandleMsgExecuteContract!!!!!!!!!!!!! "))
	res, err := k.Execute(ctx, msg.Contract, msg.FromAddress, msg.Args)
	if err != nil {
		return wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	}
	em := ctx.EventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(wasm.EventType,
			sdk.NewAttribute(wasm.AttributeKeyMethod, wasm.AttributeValueInvoke),
			sdk.NewAttribute(wasm.AttributeKeyAddress, msg.Contract.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, wasm.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	})
	res.Events = em.Events()
	return
}

func handleMsgMigrateContract(ctx sdk.Context, k Keeper, msg wasm.MsgMigrateContract) (res sdk.Result) {
	gasLimit := ctx.GasLimit()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	contractAddr, err := k.Migrate(ctx, msg.CodeHash, msg.FromAddress, msg.Contract, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	em := ctx.EventManager()
	em.EmitEvents(sdk.Events{
		sdk.NewEvent(wasm.EventType,
			sdk.NewAttribute(wasm.AttributeKeyMethod, wasm.AttributeValueMigrate),
			sdk.NewAttribute(wasm.AttributeKeyAddress, contractAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyModule, wasm.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress.String()),
		),
	})
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: em.Events(),
	}
	return
}