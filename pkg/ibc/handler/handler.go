package handler

import (
	"encoding/hex"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/ibc/keeper"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
)

func NewHandler(k keeper.IBCKeeper) sdk.Handler {
	return func(ctx sdk.Context, tx sdk.Tx) sdk.Result {
		switch tx := tx.(type) {
		case *types.IBCTransfer:
			return handleMsgIBCTransfer(ctx, k, *tx)
		default:
			errMsg := fmt.Sprintf("unrecognized supply message type: %T", tx)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// 跨链消息
func handleMsgIBCTransfer(ctx sdk.Context, k keeper.IBCKeeper, tx types.IBCTransfer) sdk.Result {
	uuidStr := keeper.GenerateUniqueID()
	// 编码
	retbz, _ := hex.DecodeString(uuidStr)

	ibcMsg, err := makeIBCMsg(retbz, tx)
	if err != nil {
		return sdk.ErrInternal("make ibc msg failed").TraceSDK(err.Error()).Result()
	}
	k.SetIBCMsg(ctx, ibcMsg)
	return sdk.Result{Data: retbz}
}

func makeIBCMsg(uuidBz []byte, tx types.IBCTransfer) (types.IBCMsg, error) {
	txbz, err := types.IbcCdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		return types.IBCMsg{}, err
	}
	ibcMsg := types.IBCMsg{
		UniqueID: 	uuidBz,
		Raw: 		txbz,
		State: 		keeper.StateReady,
	}
	return ibcMsg, nil
}