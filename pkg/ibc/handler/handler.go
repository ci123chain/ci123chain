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
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgApplyIBCTx(ctx sdk.Context, k keeper.IBCKeeper, tx types.ApplyIBCTx) sdk.Result {
	signedIBCMsg, err := k.ApplyIBCMsg(ctx, tx.UniqueID, tx.ObserverID)
	if err != nil {
		return sdk.ErrInternal("Get SignedMsg Error ").TraceSDK(err.Error()).Result()
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

func makeIBCMsg(uuidBz []byte, tx types.IBCTransfer) (types.IBCMsg, error) {
	txbz, err := types.IbcCdc.MarshalJSON(tx)
	if err != nil {
		return types.IBCMsg{}, err
	}
	ibcMsg := types.IBCMsg{
		UniqueID: 	uuidBz,
		Raw: 		txbz,
		State: 		types.StateReady,
	}
	return ibcMsg, nil
}