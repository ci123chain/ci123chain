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
		case *wasm.StoreCodeTx:
			return handleStoreCodeTx(ctx, k, *tx)
		case *wasm.UninstallCodeTx:
			return handleUninstallCodeTx(ctx, k, *tx)
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

	codeHash, Err := k.Install(ctx, msg.Sender, msg.WASMByteCode)
	if Err != nil {
		return wasm.ErrCreateFailed(wasm.DefaultCodespace, Err).Result()
	}

	return sdk.Result{
		Data:   codeHash,
		Events: ctx.EventManager().Events(),
	}
}

func handleUninstallCodeTx(ctx sdk.Context, k Keeper, msg wasm.UninstallCodeTx) sdk.Result {
	err := msg.ValidateBasic()
	if err != nil {
		return wasm.ErrInvalidMsg(wasm.DefaultCodespace, err).Result()
	}

	Err := k.Uninstall(ctx, msg.Sender, msg.CodeHash)
	if Err != nil {
		return wasm.ErrUninstallFailed(wasm.DefaultCodespace, Err).Result()
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func handleInstantiateContractTx(ctx sdk.Context, k Keeper, msg wasm.InstantiateContractTx) (res sdk.Result) {
	gasLimit := msg.GetGas()
	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
	SetGasWanted(gasWanted)

	//defer func() {
	//	if r := recover(); r != nil{
	//		var err error
	//		switch x := r.(type) {
	//		case string:
	//			err = errors.New(x)
	//			res = wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	//		case error:
	//			err = x
	//			res = wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	//		case sdk.ErrorOutOfGas:
	//			err = errors.New(x.Descriptor)
	//			res = wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	//			res.GasUsed = gasLimit
	//			res.GasWanted = gasLimit
	//		default:
	//			err = errors.New("")
	//			res = wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	//		}
	//	}
	//}()

	contractAddr, err := k.Instantiate(ctx, msg.CodeHash, msg.Sender, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe)
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

	//defer func() {
	//	if r := recover(); r != nil{
	//		var err error
	//		switch x := r.(type) {
	//		case string:
	//			err = errors.New(x)
	//			res = wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	//		case error:
	//			err = x
	//			res = wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	//		case sdk.ErrorOutOfGas:
	//			err = errors.New(x.Descriptor)
	//			res = wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	//			res.GasUsed = gasLimit
	//			res.GasWanted = gasLimit
	//		default:
	//			err = errors.New("")
	//			res = wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	//		}
	//	}
	//}()

	res, err := k.Execute(ctx, msg.Contract, msg.Sender,msg.Args)
	if err != nil {
		return wasm.ErrExecuteFailed(wasm.DefaultCodespace, err).Result()
	}
	res.Events = ctx.EventManager().Events()
	res.GasUsed = ctx.GasMeter().GasConsumed()
	return
}