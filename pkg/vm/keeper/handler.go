package keeper

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	evm "github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
)

const (
	InstantiateFuncName = "init"
)

func NewHandler(k *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch tx := msg.(type) {
		//case *wasm.MsgUploadContract:
		//	return handleMsgUploadContract(ctx, *k, *tx)
		//case *wasm.MsgInstantiateContract:
		//	return handleMsgInstantiateContract(ctx, *k, *tx)
		//case *wasm.MsgExecuteContract:
		//	return handleMsgExecuteContract(ctx, *k, *tx)
		//case *wasm.MsgMigrateContract:
		//	return handleMsgMigrateContract(ctx, *k, *tx)
		case *evm.MsgEvmTx:
			return handleMsgEvmTx(ctx, k, *tx)
		case types.MsgEthereumTx:
			return handleMsgEthereumTx(ctx, k, tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return nil, errors.New(errMsg)
		}
	}
}
//
//func handleMsgUploadContract(ctx sdk.Context, k Keeper, msg wasm.MsgUploadContract) (res *sdk.Result, err error) {
//	codeHash, err := k.Upload(ctx, msg.Code, msg.FromAddress)
//	if err != nil {
//		return nil, err
//	}
//	res = &sdk.Result{
//		Data:  codeHash,
//		Events: ctx.EventManager().Events(),
//	}
//	return
//}
//
//func handleMsgInstantiateContract(ctx sdk.Context, k Keeper, msg wasm.MsgInstantiateContract) (res *sdk.Result, err error) {
//	gasLimit := ctx.GasLimit()
//	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
//
//	args, err := wasm.CallData2Input(msg.Args)
//	if err != nil {
//		return nil, err
//	}
//	contractAddr, err := k.Instantiate(ctx, msg.CodeHash, msg.FromAddress, args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe, wasm.EmptyAddress, gasWanted)
//	if err != nil {
//		return nil, err
//	}
//	res = &sdk.Result{
//		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
//		Events: ctx.EventManager().Events(),
//	}
//	return
//}
//
//func handleMsgExecuteContract(ctx sdk.Context, k Keeper, msg wasm.MsgExecuteContract) (res *sdk.Result, err error){
//	gasLimit := ctx.GasLimit()
//	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
//
//	args, err := wasm.CallData2Input(msg.Args)
//	if err != nil {
//		return nil, err
//	}
//	result, Err := k.Execute(ctx, msg.Contract, msg.FromAddress, args, gasWanted)
//	if Err != nil {
//		return nil, Err
//	}
//	res = &result
//	res.Events = ctx.EventManager().Events()
//	return
//}
//
//func handleMsgMigrateContract(ctx sdk.Context, k Keeper, msg wasm.MsgMigrateContract) (res *sdk.Result, err error) {
//	gasLimit := ctx.GasLimit()
//	gasWanted := gasLimit - ctx.GasMeter().GasConsumed()
//
//	args, err := wasm.CallData2Input(msg.Args)
//	if err != nil {
//		return nil, err
//	}
//	contractAddr, err := k.Migrate(ctx, msg.CodeHash, msg.FromAddress, msg.Contract, args, msg.Name, msg.Version, msg.Author, msg.Email, msg.Describe, gasWanted)
//	if err != nil {
//		return nil, err
//	}
//	res = &sdk.Result{
//		Data:  []byte(fmt.Sprintf("%s", contractAddr.String())),
//		Events: ctx.EventManager().Events(),
//	}
//	return
//}

// handleMsgEvmTx handles an Ethereum specific tx
func handleMsgEvmTx(ctx sdk.Context, k *Keeper, msg evm.MsgEvmTx) (*sdk.Result, error) {
	// parse the chainID from a string to a base-10 integer
	chainIDEpoch := big.NewInt(util.CHAINID)

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
		return nil, errors.New("chain configs not found")
	}

	executionResult, err := st.TransitionDb(ctx, config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("err transitionDb: %s", err.Error()))
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
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Data.Amount.String())),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(evm.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(sender.String())),
		),
	})

	if msg.Data.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				evm.EventTypeEvmTx,
				sdk.NewAttribute([]byte(evm.AttributeKeyRecipient), []byte(msg.Data.Recipient.String())),
			),
		)
	}
	
	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events()
	return executionResult.Result, nil
}

func handleMsgEthereumTx(ctx sdk.Context, k *Keeper, msg types.MsgEthereumTx) (*sdk.Result, error) {
	// parse the chainID from a string to a base-10 integer
	chainIDEpoch := big.NewInt(util.CHAINID)

	// Verify signature and retrieve sender address
	sender, err := msg.VerifySig(chainIDEpoch)
	if err != nil {
		return nil, err
	}

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
		return nil, errors.New("chain configs not found")
	}

	executionResult, err := st.TransitionDb(ctx, config)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("err transitionDb: %s", err.Error()))
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
			sdk.NewAttribute([]byte(sdk.AttributeKeyAmount), []byte(msg.Data.Amount.String())),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute([]byte(sdk.AttributeKeyModule), []byte(evm.AttributeValueCategory)),
			sdk.NewAttribute([]byte(sdk.AttributeKeySender), []byte(sender.String())),
		),
	})

	if msg.Data.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				evm.EventTypeEvmTx,
				sdk.NewAttribute([]byte(evm.AttributeKeyRecipient), []byte(msg.Data.Recipient.String())),
			),
		)
	}

	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events()
	return executionResult.Result, nil
}



