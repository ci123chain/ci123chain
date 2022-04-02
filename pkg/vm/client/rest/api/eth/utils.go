package eth

import (
	"fmt"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	clientcontext "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/libs"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	vmmodule "github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"math/big"
	"reflect"
)

// EthTransactionsFromTendermint returns a slice of ethereum transaction hashes and the total gas usage from a set of
// tendermint block transactions.
func EthTransactionsFromTendermint(clientCtx clientcontext.Context, txs []tmtypes.Tx, blockHash common.Hash, blockNumber uint64) ([]common.Hash, *big.Int, []*Transaction, error) {
	var transactionHashes []common.Hash
	var transactions []*Transaction
	gasUsed := big.NewInt(0)
	index := uint64(0)

	for _, tx := range txs {
		ethTx, err := RawTxToEthTx(clientCtx, tx)
		if err != nil {
			// continue to next transaction in case it's not a MsgEthereumTx
			continue
		}
		// TODO: Remove gas usage calculation if saving gasUsed per block
		gasUsed.Add(gasUsed, big.NewInt(int64(ethTx.GetGas())))
		transactionHashes = append(transactionHashes, common.BytesToHash(tx.Hash()))
		tx, err := NewTransaction(ethTx, common.BytesToHash(tx.Hash()), blockHash, blockNumber, index)
		if err == nil {
			transactions = append(transactions, tx)
			index++
		}
	}

	return transactionHashes, gasUsed, transactions, nil
}

// RawTxToEthTx returns a evm MsgEthereum transaction from raw tx bytes.
func RawTxToEthTx(clientCtx clientcontext.Context, bz []byte) (*types.MsgEthereumTx, error) {
	tx, err := types.DefaultTxDecoder(clientCtx.Cdc)(bz)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	ethTx, ok := tx.(*types.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type %T, expected %T", tx, evmtypes.MsgEvmTx{})
	}
	return ethTx, nil
}



// EthBlockFromTendermint returns a JSON-RPC compatible Ethereum blockfrom a given Tendermint block.
func EthBlockFromTendermint(clientCtx clientcontext.Context, block *tmtypes.Block, gasUsedI *big.Int, transactions interface{}) (map[string]interface{}, error) {
	gasLimit := DefaultRPCGasLimit

	r, _ := libs.RetryI(0, func(retryTimes int) (res interface{}, err error) {
		res, _, _, err = clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", vmmodule.ModuleName, evmtypes.QueryBloom, block.Height), nil, false)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
	res := r.([]byte)
	var bloomRes evmtypes.QueryBloomFilter
	types.GetCodec().MustUnmarshalJSON(res, &bloomRes)

	bloom := bloomRes.Bloom

	return formatBlock(block.Header, block.Size(), int64(gasLimit), gasUsedI, transactions, bloom), nil
}

func formatBlock(
	header tmtypes.Header, size int, gasLimit int64,
	gasUsed *big.Int, transactions interface{}, bloom ethtypes.Bloom,
) map[string]interface{} {
	//if len(header.DataHash) == 0 {
	//	header.DataHash = common.Hash{}.Bytes()
	//}

	//var transactionRoot common.Hash
	//if len(header.DataHash) == 0 {
	//	transactionRoot = common.BytesToHash(ethtypes.EmptyRootHash.Bytes())
	//}else {
	//	transactionRoot = common.BytesToHash(header.DataHash.Bytes())
	//}

	res := map[string]interface{}{
		"number":           hexutil.Uint64(header.Height),
		"hash":             hexutil.Bytes(header.Hash()),
		"parentHash":       common.BytesToHash(header.LastBlockID.Hash.Bytes()),//hexutil.Bytes(header.LastBlockID.Hash),
		"nonce":            nil, // PoW specific
		"sha3Uncles":       ethtypes.EmptyUncleHash, //common.Hash{},     // No uncles in Tendermint
		"logsBloom":        bloom,
		"transactionsRoot": ethtypes.EmptyRootHash,//hexutil.Bytes(header.DataHash),
		"stateRoot":        common.BytesToHash(header.AppHash),//hexutil.Bytes(header.AppHash),
		"miner":            common.Address{},
		"mixHash":          common.Hash{},
		"difficulty":       hexutil.Uint64(0), //big.NewInt(0),
		//"totalDifficulty":  0,
		"extraData":        []byte(""),
		"size":             hexutil.Uint64(size),
		"gasLimit":         hexutil.Uint64(gasLimit), // Static gas limit
		"gasUsed":          hexutil.Uint64(gasUsed.Uint64()),//(*hexutil.Big)(gasUsed),
		"timestamp":        hexutil.Uint64(header.Time.Unix()),
		"uncles":           []string{},
		"receiptsRoot":     common.Hash{},
		"totalTx": 			header.TotalTxs,
		"numbTx": 			header.NumTxs,
	}

	if !reflect.ValueOf(transactions).IsNil() {
		switch transactions.(type) {
		case []common.Hash:
			res["transactions"] = transactions.([]common.Hash)
		case []*Transaction:
			res["transactions"] = transactions.([]*Transaction)
		}
	} else {
		res["transactions"] = []common.Hash{}
	}
	return res
}

