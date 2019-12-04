package utils

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"time"

	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func QueryTx(cliCtx context.Context, hashHexStr string) (sdk.TxResponse, sdk.Error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return sdk.TxResponse{}, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	resTx, err := node.Tx(hash, true)
	if err != nil {
		return sdk.TxResponse{}, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resTx})
	if err != nil {
		return sdk.TxResponse{}, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}
	out, err := formatTxResult(cliCtx.Cdc, resTx, resBlocks[resTx.Height])
	if err != nil {
		return out, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	return out, nil
}


func getBlocksForTxResults(cliCtx context.Context, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*ctypes.ResultBlock)

	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(&resTx.Height)
			if err != nil {
				return nil, err
			}
			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}


func formatTxResult(cdc *codec.Codec, resTx *ctypes.ResultTx, resBlock *ctypes.ResultBlock) (sdk.TxResponse, error) {
	tx, err := parseTx(cdc, resTx.Tx)
	if err != nil {
		return sdk.TxResponse{}, err
	}
	return sdk.NewResponseResultTx(resTx, tx, resBlock.Block.Time.Format(time.RFC3339)), nil
}


func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {

	// todo: only TransferTx implement

	tx := new(transfer.TransferTx)
	rlp.DecodeBytes(txBytes, &tx)
	//err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	//if err != nil {
	//	return nil, err
	//}

	return tx, nil
}

