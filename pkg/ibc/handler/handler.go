package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/types"
)

func NewHandler(k keeper.IBCKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.IBCTransfer:
			return handleMsgIBCTransfer(ctx, k, *msg)
		case *types.MsgApplyIBC:
			return handleMsgApplyIBCTx(ctx, k, *msg)
		case *types.IBCMsgBankSend:
			return handleMsgIBCBankSendTx(ctx, k, *msg)
		case *types.IBCReceiveReceiptMsg:
			return handleMsgReceiveReceipt(ctx, k, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", msg)
			return nil, errors.New(errMsg)
		}
	}
}

// 新增跨链消息
func handleMsgIBCTransfer(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCTransfer) (*sdk.Result, error) {
	uuidStr := keeper.GenerateUniqueID(tx.Bytes())

	retbz := []byte(uuidStr)
	ibcMsg, err := makeIBCMsg([]byte(uuidStr), tx)
	if err != nil {
		return nil, err
	}
	err = k.SetIBCMsg(ctx, ibcMsg)
	if err != nil {
		return nil, err
	}
	ctx.Logger().Info("Create IBCTransaction successed")

	return &sdk.Result{Data: retbz}, nil
}

// 第一步: 申请处理跨链交易
func handleMsgApplyIBCTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.MsgApplyIBC) (*sdk.Result, error) {
	signedIBCMsg, err := k.ApplyIBCMsg(ctx, tx)
	if err != nil {
		return nil, err
	}
	signedIBCMsgBz, _ := json.Marshal(signedIBCMsg)

	ctx.Logger().Info("Apply IBCTransaction successed")
	return &sdk.Result{Data: signedIBCMsgBz}, nil
}

// 第二步: 跨链交易 bank 转账 (fabric -> ci)
func handleMsgIBCBankSendTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCMsgBankSend) (*sdk.Result, error) {

	ibcMsg, err := keeper.ValidateRawIBCMessage(tx)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		receiptBz, _ := json.Marshal(*receipt)
		return &sdk.Result{Data: receiptBz}, nil
	}

	ibcMsg.State = types.StateDone

	receipt, err2 := k.MakeBankReceipt(ctx, *ibcMsg)
	if err2 != nil {
		return nil, err2
	}

	// todo bank action
	err2 = k.BankSend(ctx, *ibcMsg)
	if err2 != nil {
		return nil, err2
	}

	// 保存该交易
	err2 = k.SetIBCMsg(ctx, *ibcMsg)
	if err2 != nil {
		return nil, err2
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

	return &sdk.Result{Data: receiptBz}, nil
}

// 接收到回执消息
func handleMsgReceiveReceipt(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCReceiveReceiptMsg) (*sdk.Result, error)  {

	receiveObj, err := keeper.ValidateRawReceiptMessage(tx)
	if err != nil {
		return nil, err
	}
	err2 := k.ReceiveReceipt(ctx, *receiveObj)
	if err2 != nil {
		return nil, err2
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

	return &sdk.Result{Data: []byte(receiveObj.UniqueID)}, nil
}

// 生成 IbcInfo
func makeIBCMsg(uuidBz []byte, tx types.IBCTransfer) (types.IBCInfo, error) {
	ibcMsg := types.IBCInfo{
		UniqueID: 		uuidBz,
		FromAddress: 	tx.FromAddress,
		ToAddress: 		tx.ToAddress,
		Amount: 		tx.Coin,
		State: 			types.StateReady,
	}
	return ibcMsg, nil
}