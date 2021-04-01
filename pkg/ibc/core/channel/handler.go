package channel

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	"github.com/pkg/errors"
)

// HandleMsgChannelOpenInit defines the sdk.Handler for MsgChannelOpenInit
func HandleMsgChannelOpenInit(ctx sdk.Context, k keeper.Keeper, portCap *capabilitytypes.Capability, msg *types.MsgChannelOpenInit) (*sdk.Result, string, *capabilitytypes.Capability, error) {
	channelID, capKey, err := k.ChanOpenInit(
		ctx, msg.Channel.Ordering, msg.Channel.ConnectionHops, msg.PortId,
		portCap, msg.Channel.Counterparty, msg.Channel.Version,
	)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "channel handshake open init failed")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeChannelOpenInit,
			sdk.NewAttributeString(types.AttributeKeyPortID, msg.PortId),
			sdk.NewAttributeString(types.AttributeKeyChannelID, channelID),
			sdk.NewAttributeString(types.AttributeCounterpartyPortID, msg.Channel.Counterparty.PortId),
			sdk.NewAttributeString(types.AttributeCounterpartyChannelID, msg.Channel.Counterparty.ChannelId),
			sdk.NewAttributeString(types.AttributeKeyConnectionID, msg.Channel.ConnectionHops[0]),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, channelID, capKey, nil
}

// HandleMsgChannelOpenTry defines the sdk.Handler for MsgChannelOpenTry
func HandleMsgChannelOpenTry(ctx sdk.Context, k keeper.Keeper, portCap *capabilitytypes.Capability, msg *types.MsgChannelOpenTry) (*sdk.Result, string, *capabilitytypes.Capability, error) {
	channelID, capKey, err := k.ChanOpenTry(ctx, msg.Channel.Ordering, msg.Channel.ConnectionHops, msg.PortId, msg.PreviousChannelId,
		portCap, msg.Channel.Counterparty, msg.Channel.Version, msg.CounterpartyVersion, msg.ProofInit, msg.ProofHeight,
	)
	if err != nil {
		return nil, "", nil, errors.Wrap(err, "channel handshake open try failed")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeChannelOpenTry,
			sdk.NewAttributeString(types.AttributeKeyPortID, msg.PortId),
			sdk.NewAttributeString(types.AttributeKeyChannelID, channelID),
			sdk.NewAttributeString(types.AttributeCounterpartyPortID, msg.Channel.Counterparty.PortId),
			sdk.NewAttributeString(types.AttributeCounterpartyChannelID, msg.Channel.Counterparty.ChannelId),
			sdk.NewAttributeString(types.AttributeKeyConnectionID, msg.Channel.ConnectionHops[0]),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, channelID, capKey, nil
}

// HandleMsgChannelOpenAck defines the sdk.Handler for MsgChannelOpenAck
func HandleMsgChannelOpenAck(ctx sdk.Context, k keeper.Keeper, channelCap *capabilitytypes.Capability, msg *types.MsgChannelOpenAck) (*sdk.Result, error) {
	err := k.ChanOpenAck(
		ctx, msg.PortId, msg.ChannelId, channelCap, msg.CounterpartyVersion, msg.CounterpartyChannelId, msg.ProofTry, msg.ProofHeight,
	)
	if err != nil {
		return nil, errors.Wrap(err, "channel handshake open ack failed")
	}

	channel, _ := k.GetChannel(ctx, msg.PortId, msg.ChannelId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeChannelOpenAck,
			sdk.NewAttributeString(types.AttributeKeyPortID, msg.PortId),
			sdk.NewAttributeString(types.AttributeKeyChannelID, msg.ChannelId),
			sdk.NewAttributeString(types.AttributeCounterpartyPortID, channel.Counterparty.PortId),
			sdk.NewAttributeString(types.AttributeCounterpartyChannelID, channel.Counterparty.ChannelId),
			sdk.NewAttributeString(types.AttributeKeyConnectionID, channel.ConnectionHops[0]),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}

// HandleMsgChannelOpenConfirm defines the sdk.Handler for MsgChannelOpenConfirm
func HandleMsgChannelOpenConfirm(ctx sdk.Context, k keeper.Keeper, channelCap *capabilitytypes.Capability, msg *types.MsgChannelOpenConfirm) (*sdk.Result, error) {
	err := k.ChanOpenConfirm(ctx, msg.PortId, msg.ChannelId, channelCap, msg.ProofAck, msg.ProofHeight)
	if err != nil {
		return nil, errors.Wrap(err, "channel handshake open confirm failed")
	}

	channel, _ := k.GetChannel(ctx, msg.PortId, msg.ChannelId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeChannelOpenConfirm,
			sdk.NewAttributeString(types.AttributeKeyPortID, msg.PortId),
			sdk.NewAttributeString(types.AttributeKeyChannelID, msg.ChannelId),
			sdk.NewAttributeString(types.AttributeCounterpartyPortID, channel.Counterparty.PortId),
			sdk.NewAttributeString(types.AttributeCounterpartyChannelID, channel.Counterparty.ChannelId),
			sdk.NewAttributeString(types.AttributeKeyConnectionID, channel.ConnectionHops[0]),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().Events(),
	}, nil
}
