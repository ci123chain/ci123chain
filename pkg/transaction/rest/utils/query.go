package utils

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/abci/codec"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/context"
	"encoding/hex"
	"time"

	sdk "gitlab.oneitfarm.com/blockchain/ci123chain/pkg/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func QueryTx(cliCtx context.Context, hashHexStr string) (sdk.TxResponse, error) {
	hash, err := hex.DecodeString(hashHexStr)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	node, err := cliCtx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}

	resTx, err := node.Tx(hash, true)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	resBlocks, err := getBlocksForTxResults(cliCtx, []*ctypes.ResultTx{resTx})
	if err != nil {
		return sdk.TxResponse{}, err
	}
	out, err := formatTxResult(cliCtx.Cdc, resTx, resBlocks[resTx.Height])
	if err != nil {
		return out, err
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
	//tx, err := parseTx(cdc, resTx.Tx)
	//if err != nil {
	//	return sdk.TxResponse{}, err
	//}
	//fmt.Println(resTx.Tx)
	return sdk.NewResponseResultTx(resTx, nil, resBlock.Block.Time.Format(time.RFC3339)), nil
}


//func parseTx(cdc *codec.Codec, txBytes []byte) (sdk.Tx, error) {
	//var tx types.StdTx
	//
	//err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return tx, nil
//}

