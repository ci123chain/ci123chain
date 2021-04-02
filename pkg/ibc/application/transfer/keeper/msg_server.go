package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
)

var _ types.MsgServer = Keeper{}

// See createOutgoingPacket in spec:https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay

// Transfer defines a rpc handler method for MsgTransfer.
func (k Keeper) Transfer(ctx sdk.Context, msg *types.MsgTransfer) (*types.MsgTransferResponse, error) {
	sender := sdk.HexToAddress(msg.Sender)
	if err := k.SendTransfer(
		ctx, msg.SourcePort, msg.SourceChannel, msg.Token, sender, msg.Receiver, msg.TimeoutHeight, msg.TimeoutTimestamp,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC fungible token transfer", "token", msg.Token.Denom, "amount", msg.Token.Amount.String(), "sender", msg.Sender, "receiver", msg.Receiver)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttributeString(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttributeString(types.AttributeKeyReceiver, msg.Receiver),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgTransferResponse{}, nil
}
