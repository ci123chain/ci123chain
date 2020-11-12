package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	evm "github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	wasm "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	"github.com/ethereum/go-ethereum/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
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
		case *evm.MsgEvmTx:
			return handleMsgEvmTx(ctx, k, *tx)
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
		Data:  codeHash,
		Events: ctx.EventManager().Events(),
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

	contractAddr, err := k.Migrate(ctx, msg.CodeHash, msg.FromAddress, msg.Contract, msg.Args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe)
	if err != nil {
		return wasm.ErrInstantiateFailed(wasm.DefaultCodespace, err).Result()
	}
	res = sdk.Result{
		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
		Events: ctx.EventManager().Events(),
	}
	return
}

// handleMsgEvmTx handles an Ethereum specific tx
func handleMsgEvmTx(ctx sdk.Context, k Keeper, msg evm.MsgEvmTx) sdk.Result {
	// parse the chainID from a string to a base-10 integer
	//todo
	chainIDEpoch := big.NewInt(123)

	// Verify signature and retrieve sender address
	sender := common.BytesToAddress(msg.GetFromAddress().Bytes())

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)

	st := evm.StateTransition{
		AccountNonce: msg.Data.AccountNonce,
		Price:        msg.Data.Price,
		GasLimit:     msg.Data.GasLimit,
		Recipient:    msg.Data.Recipient,
		Amount:       msg.Data.Amount,
		Payload:      msg.Data.Payload,
		Csdb:         k.CommitStateDB.WithContext(ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &ethHash,
		Sender:       sender,
		Simulate:     ctx.IsCheckTx(),
	}

	// since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different than the
	// other nodes, causing a consensus error
	if !st.Simulate {
		// Prepare db for logs
		// TODO: block hash
		k.CommitStateDB.Prepare(ethHash, common.Hash{}, k.TxCount)
		k.TxCount++
	}

	config, found := k.GetChainConfig(ctx)
	if !found {
		return evm.ErrChainConfigNotFound(evm.DefaultCodespace, "chain config not found").Result()
	}

	executionResult, err := st.TransitionDb(ctx, config)
	if err != nil {
		return evm.ErrTransitionDb(evm.DefaultCodespace, fmt.Sprintf("err transitionDb: %s", err.Error())).Result()
	}

	if !st.Simulate {
		// update block bloom filter
		k.Bloom.Or(k.Bloom, executionResult.Bloom)

		// update transaction logs in KVStore
		err = k.SetLogs(ctx, common.BytesToHash(txHash), executionResult.Logs)
		if err != nil {
			panic(err)
		}
	}

	// log successful execution
	k.Logger(ctx).Info(executionResult.Result.Log)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			evm.EventTypeEvmTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Data.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, evm.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		),
	})

	if msg.Data.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				evm.EventTypeEvmTx,
				sdk.NewAttribute(evm.AttributeKeyRecipient, msg.Data.Recipient.String()),
			),
		)
	}

	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events()
	return *executionResult.Result
}