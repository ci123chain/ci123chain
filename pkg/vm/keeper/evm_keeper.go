package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	evm "github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/types"
	"github.com/ethereum/go-ethereum/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
)


var _ types.Keeper = (*Keeper)(nil)

func (k Keeper) EvmTxExec(ctx sdk.Context,m sdk.Msg) (types.VMResult, error) {
	msg, ok := m.(evm.MsgEvmTx)
	if !ok {
		return nil, evm.ErrContractMsgInvalid
	}
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
		return nil, evm.ErrEVMChainConfigInvalid
	}

	executionResult, err := st.TransitionDb(ctx, config)
	if err != nil {
		return nil,  evm.ErrExecTransactionInvalid.Wrap(err.Error())
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

	return executionResult, nil
}
