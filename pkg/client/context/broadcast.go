package context

import (
	//ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

// Broadcast the transfer bytes to Tendermint
func (ctx *Context) BroadcastTx(tx []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}
	res, err := node.BroadcastTxCommit(tx)
	if err != nil {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}
	if res.CheckTx.Code != uint32(0) {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}
	if res.DeliverTx.Code != uint32(0) {
		return sdk.NewResponseFormatBroadcastTxCommit(res), err
	}

	return sdk.NewResponseFormatBroadcastTxCommit(res), nil
}

func (ctx *Context) BroadcastTxAsync(tx []byte) (sdk.TxResponse, error) {
	node, err := ctx.GetNode()
	if err != nil {
		return sdk.TxResponse{}, err
	}
	res, err := node.BroadcastTxAsync(tx)
	if err != nil {
		return sdk.TxResponse{Info:"something error", Code:res.Code}, err
	}

	if res.Code != uint32(0) {
		return sdk.TxResponse{Info:"error code", Code:res.Code}, err
	}

	return sdk.TxResponse{Info:"success",Code:res.Code}, nil
}