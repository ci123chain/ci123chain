package handler

import (
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/types"
)

func NewHandler(k keeper.IBCKeeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		switch tx := tx.(type) {
		case *types.IBCTransfer:
			return handleMsgIBCTransfer(ctx, k, *tx)
		case *types.ApplyIBCTx:
			return handleMsgApplyIBCTx(ctx, k, *tx)
		case *types.IBCMsgBankSend:
			return handleMsgIBCBankSendTx(ctx, k, *tx)
		case *types.IBCReceiveReceiptMsg:
			return handleMsgReceiveReceipt(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}



// 新增跨链消息
func handleMsgIBCTransfer(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCTransfer) sdk.Result {
	uuidStr := keeper.GenerateUniqueID(tx.Bytes())

	retbz := []byte(uuidStr)
	ibcMsg, err := makeIBCMsg([]byte(uuidStr), tx)
	if err != nil {
		return types.ErrMakeIBCMsg(types.DefaultCodespace, err).Result()
	}
	err = k.SetIBCMsg(ctx, ibcMsg)
	if err != nil {
		return types.ErrSetIBCMsg(types.DefaultCodespace, err).Result()
	}
	ctx.Logger().Info("Create IBCTransaction successed")

	return sdk.Result{Data: retbz}
}

// 第一步: 申请处理跨链交易
func handleMsgApplyIBCTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.ApplyIBCTx) sdk.Result {
	signedIBCMsg, err := k.ApplyIBCMsg(ctx, tx)
	if err != nil {
		return types.ErrApplyIBCMsg(types.DefaultCodespace, err).Result()
	}
	signedIBCMsgBz, _ := json.Marshal(signedIBCMsg)

	ctx.Logger().Info("Apply IBCTransaction successed")
	return sdk.Result{Data: signedIBCMsgBz}
}

// 第二步: 跨链交易 bank 转账 (fabric -> ci)
func handleMsgIBCBankSendTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCMsgBankSend) sdk.Result {

	ibcMsg, err := keeper.ValidateRawIBCMessage(tx)
	if err != nil {
		return err.Result()
	}

	// todo warning
	//ibc := k.GetIBCByUniqueID(ctx, ibcMsg.UniqueID)
	//if ibc != nil {
	//	return sdk.ErrUnknownRequest("ibcTx already exist with uniqueID " + string(ibc.UniqueID)).Result()
	//}
	old_ibcMsg := k.GetIBCByUniqueID(ctx, ibcMsg.UniqueID)
	if old_ibcMsg != nil && old_ibcMsg.State == types.StateDone {
		receipt, err := k.MakeBankReceipt(ctx, *ibcMsg)
		if err != nil {
			return types.ErrMakeBankReceipt(types.DefaultCodespace, err).Result()
		}

		receiptBz, _ := json.Marshal(*receipt)
		return sdk.Result{Data: receiptBz}
	}

	ibcMsg.State = types.StateDone

	receipt, err2 := k.MakeBankReceipt(ctx, *ibcMsg)
	if err2 != nil {
		return types.ErrMakeBankReceipt(types.DefaultCodespace, err2).Result()
	}

	// todo bank action
	err2 = k.BankSend(ctx, *ibcMsg)
	if err2 != nil {
		return types.ErrBankSend(types.DefaultCodespace, err2).Result()
	}

	// 保存该交易
	err2 = k.SetIBCMsg(ctx, *ibcMsg)
	if err2 != nil {
		return types.ErrSetIBCMsg(types.DefaultCodespace, err2).Result()
	}
	receiptBz, _ := json.Marshal(*receipt)

	//bank转账成功，observer nonce+1
	//account := k.AccountKeeper.GetAccount(ctx, tx.From)
	//saveErr := account.SetSequence(account.GetSequence() + 1)
	//if saveErr != nil {
	//	return sdk.ErrInternal("Failed to set sequence").Result()
	//}
	//k.AccountKeeper.SetAccount(ctx, account)
	//

	ctx.Logger().Info("Handle IBCTransaction successed")

	return sdk.Result{Data: receiptBz}
}


// 接收到回执消息
func handleMsgReceiveReceipt(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCReceiveReceiptMsg) sdk.Result  {

	receiveObj, err := keeper.ValidateRawReceiptMessage(tx)
	if err != nil {
		return err.Result()
	}
	err2 := k.ReceiveReceipt(ctx, *receiveObj)
	if err2 != nil {
		return types.ErrReceiveReceipt(types.DefaultCodespace, err2).Result()
	}
	//交易成功，observer nonce+1
	//account := k.AccountKeeper.GetAccount(ctx, tx.From)
	//saveErr := account.SetSequence(account.GetSequence() + 1)
	//if saveErr != nil {
	//	return transaction.ErrSetSequence(types.DefaultCodespace, saveErr.Error()).Result()
	//}
	//k.AccountKeeper.SetAccount(ctx, account)
	//
	ctx.Logger().Info("Handle receipt successed")

	return sdk.Result{Data: []byte(receiveObj.UniqueID)}
}

// 生成 IbcInfo
func makeIBCMsg(uuidBz []byte, tx types.IBCTransfer) (types.IBCInfo, error) {
	ibcMsg := types.IBCInfo{
		UniqueID: 		uuidBz,
		FromAddress: 	tx.From,
		ToAddress: 		tx.ToAddress,
		Amount: 		tx.Coin,
		State: 			types.StateReady,
	}
	return ibcMsg, nil
}