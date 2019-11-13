package handler

import (
	"encoding/json"
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/ibc/keeper"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
)

func NewHandler(k keeper.IBCKeeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		switch tx := tx.(type) {
		case *types.IBCTransfer:
			return handleMsgIBCTransfer(ctx, k, *tx)
		case *types.ApplyIBCTx:
			return handleMsgApplyIBCTx(ctx, k, *tx)
		case *types.IBCMsgBankSend:
			return handleMsgIBCSendTx(ctx, k, *tx)
		case *types.IBCReceiveReceiptMsg:
			return handleMsgReceiveReceipt(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// 跨链交易第二步 (fabric -> ci)
func handleMsgIBCSendTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCMsgBankSend) sdk.Result {
	ibcMsg, err := keeper.ValidateRawIBCMessage(tx)
	if err != nil {
		return sdk.ErrUnknownRequest("Bank pkg invalid").TraceSDK(err.Error()).Result()
	}
	// 去重
	//ibc := k.GetIBCByUniqueID(ctx, ibcMsg.UniqueID)
	//if ibc != nil {
	//	return sdk.ErrUnknownRequest("ibcTx already exist with uniqueID " + string(ibc.UniqueID)).Result()
	//}

	// todo bank action
	err = k.BankSend(ctx, *ibcMsg)
	if err != nil {
		return sdk.ErrInsufficientCoins(err.Error()).Result()
	}

	receipt, err := k.MakeBankReceipt(ctx, *ibcMsg)
	if err != nil {
		return sdk.ErrUnknownRequest("Get bank receipt error").TraceSDK(err.Error()).Result()
	}

	// 保存该交易
	err = k.SetIBCMsg(ctx, *ibcMsg)
	if err != nil {
		return sdk.ErrUnknownRequest("Save ibcMsg error").TraceSDK(err.Error()).Result()
	}
	receiptBz, _ := json.Marshal(receipt)
	return sdk.Result{Data: receiptBz}
}


func handleMsgApplyIBCTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.ApplyIBCTx) sdk.Result {
	signedIBCMsg, err := k.ApplyIBCMsg(ctx, tx.UniqueID, tx.ObserverID)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	signedIBCMsgBz, _ := json.Marshal(signedIBCMsg)

	return sdk.Result{Data: signedIBCMsgBz}
}


// 跨链消息
func handleMsgIBCTransfer(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCTransfer) sdk.Result {
	uuidStr := keeper.GenerateUniqueID(tx.Bytes())

	retbz := []byte(uuidStr)
	ibcMsg, err := makeIBCMsg([]byte(uuidStr), tx)
	if err != nil {
		return sdk.ErrInternal("make ibc msg failed").TraceSDK(err.Error()).Result()
	}
	k.SetIBCMsg(ctx, ibcMsg)
	return sdk.Result{Data: retbz}
}

// 接收到回执消息
func handleMsgReceiveReceipt(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCReceiveReceiptMsg) sdk.Result  {
	receiveObj, err := keeper.ValidateRawReceiptMessage(tx)
	if err != nil {
		return sdk.ErrUnknownRequest("Receipt pkg invalid").TraceSDK(err.Error()).Result()
	}
	err = k.ReceiveReceipt(ctx, *receiveObj)
	if err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return sdk.Result{Data: []byte(receiveObj.UniqueID)}
}

func makeIBCMsg(uuidBz []byte, tx types.IBCTransfer) (types.IBCMsg, error) {
	ibcMsg := types.IBCMsg{
		UniqueID: 		uuidBz,
		FromAddress: 	tx.From,
		ToAddress: 		tx.ToAddress,
		Amount: 		tx.Coin,
		State: 			types.StateReady,
	}
	return ibcMsg, nil
}