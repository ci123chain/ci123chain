package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/transfer/types"
	"time"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func QueryTx(cliCtx context.Context, hashHexStr string) (sdk.TxResponse, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return sdk.TxResponse{}, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, types.ErrQueryTx(types.DefaultCodespace, err.Error())
	}

	resTx, err := node.Tx(cliCtx.Context(), hash, true)
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

func QueryTxsWithHeight(cliCtx context.Context, heights []int64) ([]sdk.TxsResult, error) {
	var results = make([]sdk.TxsResult, 0)
	for _, v := range heights {
		var result sdk.TxsResult
		res, getErr := getTxsInfo(cliCtx, v)
		if getErr != nil {
			result.Txs = nil
			result.Height = v
			result.Error = getErr.Error()
		}else {
			result.Error = ""
			result.Height = v
			result.Txs = res
		}
		results = append(results, result)
	}
	return results, nil
}


func getBlocksForTxResults(cliCtx context.Context, resTxs []*ctypes.ResultTx) (map[int64]*ctypes.ResultBlock, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	resBlocks := make(map[int64]*ctypes.ResultBlock)

	for _, resTx := range resTxs {
		if _, ok := resBlocks[resTx.Height]; !ok {
			resBlock, err := node.Block(cliCtx.Context(), &resTx.Height)
			if err != nil {
				return nil, err
			}
			resBlocks[resTx.Height] = resBlock
		}
	}

	return resBlocks, nil
}

func getTxsInfo(cliCtx context.Context, height int64) ([]sdk.TxInfo, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}
	infoes := make([]sdk.TxInfo, 0)
	block, err := node.Block(cliCtx.Context(), &height)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, errors.New(fmt.Sprintf("the height %d, has no block", height))
	}
	for _, v := range block.Block.Data.Txs{
		///hashes = append(hashes, hex.EncodeToString(v.Hash()))
		resTx, _ := node.Tx(cliCtx.Context(), v.Hash(), true)
		index := resTx.Index
		info := sdk.TxInfo{
			Hash:  hex.EncodeToString(v.Hash()),
			Index: index,
		}
		infoes = append(infoes, info)
	}
	return infoes, nil
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

	//tx := new(transfer.TransferTx)
	var tx sdk.Tx
	//rlp.DecodeBytes(txBytes, &tx)
	err := cdc.UnmarshalBinaryBare(txBytes, &tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

