package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/ci123chain/ci123chain/pkg/wasm/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch tx := tx.(type) {
		case *wasm.InstantiateContractTx:
			return handleInstantiateContractTx(ctx, k, *tx)
		case *wasm.ExecuteContractTx:
			return handleExecuteContractTx(ctx, k, *tx)
		case *wasm.MigrateContractTx:
			return handleMigrateContractTx(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleInstantiateContractTx(ctx sdk.Context, k Keeper, msg wasm.InstantiateContractTx) (res sdk.Result) {
	gasLimit := msg.GetGas()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	contractAddr, err := k.Instantiate(ctx, msg.Code, msg.Sender, msg.Nonce, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: ctx.EventManager().Events(),
	}
	res.GasUsed = ctx.GasMeter().GasConsumed()
	return
}

func handleExecuteContractTx(ctx sdk.Context, k Keeper, msg wasm.ExecuteContractTx) (res sdk.Result){
	gasLimit := msg.GetGas()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	res, err := k.Execute(ctx, msg.Contract, msg.Sender,msg.Args)
	if err != nil {
		return wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	}
	res.Events = ctx.EventManager().Events()
	res.GasUsed = ctx.GasMeter().GasConsumed()
	return
}

func handleMigrateContractTx(ctx sdk.Context, k Keeper, msg wasm.MigrateContractTx) (res sdk.Result) {
	gasLimit := msg.GetGas()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	contractAddr, err := k.Migrate(ctx, msg.Code, msg.Sender, msg.Contract, msg.Nonce, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: ctx.EventManager().Events(),
	}
	res.GasUsed = ctx.GasMeter().GasConsumed()
	return
}