package context

import (
	//ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
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
		return sdk.NewResponseFormatBroadcastTx(res), err
	}

	return sdk.NewResponseFormatBroadcastTx(res), nil
}