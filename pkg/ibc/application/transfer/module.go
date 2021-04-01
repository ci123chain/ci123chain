package transfer

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/application/transfer/types"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
)

// AppModuleBasic is the IBC Transfer AppModuleBasic
type AppModuleBasic struct{}
// AppModule represents the AppModule for this module
type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// OnRecvPacket implements the IBCModule interface
func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {
	var data types.FungibleTokenPacketData
	if err := types.IBCTransferCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	acknowledgement := channeltypes.NewResultAcknowledgement([]byte{byte(1)})

	err := am.keeper.OnRecvPacket(ctx, packet, data)
	if err != nil {
		acknowledgement = channeltypes.NewErrorAcknowledgement(err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttributeString(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttributeString(types.AttributeKeyDenom, data.Denom),
			sdk.NewAttributeString(types.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttributeString(types.AttributeKeyAckSuccess, fmt.Sprintf("%t", err != nil)),
		),
	)

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, acknowledgement.GetBytes(), nil
}


// OnAcknowledgementPacket implements the IBCModule interface
func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {
	var ack channeltypes.Acknowledgement
	if err := types.IBCTransferCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet acknowledgement: %v", err)
	}
	var data types.FungibleTokenPacketData
	if err := types.IBCTransferCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-20 transfer packet data: %s", err.Error())
	}

	if err := am.keeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypePacket,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttributeString(types.AttributeKeyReceiver, data.Receiver),
			sdk.NewAttributeString(types.AttributeKeyDenom, data.Denom),
			sdk.NewAttributeString(types.AttributeKeyAmount, fmt.Sprintf("%d", data.Amount)),
			sdk.NewAttributeString(types.AttributeKeyAck, ack.String()),
		),
	)

	switch resp := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Result:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttributeString(types.AttributeKeyAckSuccess, string(resp.Result)),
			),
		)
	case *channeltypes.Acknowledgement_Error:
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypePacket,
				sdk.NewAttributeString(types.AttributeKeyAckError, resp.Error),
			),
		)
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
