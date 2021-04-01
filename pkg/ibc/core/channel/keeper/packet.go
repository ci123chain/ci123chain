package keeper

import (
	"bytes"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"time"
)

// SendPacket is called by a module in order to send an IBC packet on a channel
// end owned by the calling module to the corresponding module on the counterparty
// chain.
func (k Keeper) SendPacket(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	packet exported.PacketI,
) error {
	if err := packet.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "packet failed basic validation")
	}

	channel, found := k.GetChannel(ctx, packet.GetSourcePort(), packet.GetSourceChannel())
	if !found {
		return sdkerrors.Wrap(types.ErrChannelNotFound, packet.GetSourceChannel())
	}

	if channel.State == types.CLOSED {
		return sdkerrors.Wrapf(
			types.ErrInvalidChannelState,
			"channel is CLOSED (got %s)", channel.State.String(),
		)
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, channelCap, host.ChannelCapabilityPath(packet.GetSourcePort(), packet.GetSourceChannel())) {
		return sdkerrors.Wrapf(types.ErrChannelCapabilityNotFound, "caller does not own capability for channel, port ID (%s) channel ID (%s)", packet.GetSourcePort(), packet.GetSourceChannel())
	}

	if packet.GetDestPort() != channel.Counterparty.PortId {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet destination port doesn't match the counterparty's port (%s ≠ %s)", packet.GetDestPort(), channel.Counterparty.PortId,
		)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", packet.GetDestChannel(), channel.Counterparty.ChannelId,
		)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	clientState, found := k.clientKeeper.GetClientState(ctx, connectionEnd.GetClientID())
	if !found {
		return clienttypes.ErrConsensusStateNotFound
	}

	// prevent accidental sends with clients that cannot be updated
	if clientState.IsFrozen() {
		return sdkerrors.Wrapf(clienttypes.ErrClientFrozen, "cannot send packet on a frozen client with ID %s", connectionEnd.GetClientID())
	}

	// check if packet timeouted on the receiving chain
	latestHeight := clientState.GetLatestHeight()
	timeoutHeight := packet.GetTimeoutHeight()
	if !timeoutHeight.IsZero() && latestHeight.GTE(timeoutHeight) {
		return sdkerrors.Wrapf(
			types.ErrPacketTimeout,
			"receiving chain block height >= packet timeout height (%s >= %s)", latestHeight, timeoutHeight,
		)
	}

	latestTimestamp, err := k.connectionKeeper.GetTimestampAtHeight(ctx, connectionEnd, latestHeight)
	if err != nil {
		return err
	}

	if packet.GetTimeoutTimestamp() != 0 && latestTimestamp >= packet.GetTimeoutTimestamp() {
		return sdkerrors.Wrapf(
			types.ErrPacketTimeout,
			"receiving chain block timestamp >= packet timeout timestamp (%s >= %s)", time.Unix(0, int64(latestTimestamp)), time.Unix(0, int64(packet.GetTimeoutTimestamp())),
		)
	}

	nextSequenceSend, found := k.GetNextSequenceSend(ctx, packet.GetSourcePort(), packet.GetSourceChannel())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", packet.GetSourcePort(), packet.GetSourceChannel(),
		)
	}

	if packet.GetSequence() != nextSequenceSend {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet sequence ≠ next send sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceSend,
		)
	}

	commitment := types.CommitPacket(k.cdc, packet)

	nextSequenceSend++
	k.SetNextSequenceSend(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), nextSequenceSend)
	k.SetPacketCommitment(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence(), commitment)

	// Emit Event with Packet data along with other packet information for relayer to pick up
	// and relay to other chain
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSendPacket,
			sdk.NewAttributeString(types.AttributeKeyData, string(packet.GetData())),
			sdk.NewAttributeString(types.AttributeKeyTimeoutHeight, timeoutHeight.String()),
			sdk.NewAttributeString(types.AttributeKeyTimeoutTimestamp, fmt.Sprintf("%d", packet.GetTimeoutTimestamp())),
			sdk.NewAttributeString(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttributeString(types.AttributeKeySrcPort, packet.GetSourcePort()),
			sdk.NewAttributeString(types.AttributeKeySrcChannel, packet.GetSourceChannel()),
			sdk.NewAttributeString(types.AttributeKeyDstPort, packet.GetDestPort()),
			sdk.NewAttributeString(types.AttributeKeyDstChannel, packet.GetDestChannel()),
			sdk.NewAttributeString(types.AttributeKeyChannelOrdering, channel.Ordering.String()),
			// we only support 1-hop packets now, and that is the most important hop for a relayer
			// (is it going to a chain I am connected to)
			sdk.NewAttributeString(types.AttributeKeyConnection, channel.ConnectionHops[0]),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	k.Logger(ctx).Info("packet sent", "packet", fmt.Sprintf("%v", packet))
	return nil
}



// AcknowledgePacket is called by a module to process the acknowledgement of a
// packet previously sent by the calling module on a channel to a counterparty
// module on the counterparty chain. Its intended usage is within the ante
// handler. AcknowledgePacket will clean up the packet commitment,
// which is no longer necessary since the packet has been received and acted upon.
// It will also increment NextSequenceAck in case of ORDERED channels.
func (k Keeper) AcknowledgePacket(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet exported.PacketI,
	acknowledgement []byte,
	proof []byte,
	proofHeight exported.Height,
) error {
	channel, found := k.GetChannel(ctx, packet.GetSourcePort(), packet.GetSourceChannel())
	if !found {
		return sdkerrors.Wrapf(
			types.ErrChannelNotFound,
			"port ID (%s) channel ID (%s)", packet.GetSourcePort(), packet.GetSourceChannel(),
		)
	}

	if channel.State != types.OPEN {
		return sdkerrors.Wrapf(
			types.ErrInvalidChannelState,
			"channel state is not OPEN (got %s)", channel.State.String(),
		)
	}

	// Authenticate capability to ensure caller has authority to receive packet on this channel
	capName := host.ChannelCapabilityPath(packet.GetSourcePort(), packet.GetSourceChannel())
	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, capName) {
		return sdkerrors.Wrapf(
			types.ErrInvalidChannelCapability,
			"channel capability failed authentication for capability name %s", capName,
		)
	}

	// packet must have been sent to the channel's counterparty
	if packet.GetDestPort() != channel.Counterparty.PortId {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet destination port doesn't match the counterparty's port (%s ≠ %s)", packet.GetDestPort(), channel.Counterparty.PortId,
		)
	}

	if packet.GetDestChannel() != channel.Counterparty.ChannelId {
		return sdkerrors.Wrapf(
			types.ErrInvalidPacket,
			"packet destination channel doesn't match the counterparty's channel (%s ≠ %s)", packet.GetDestChannel(), channel.Counterparty.ChannelId,
		)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return sdkerrors.Wrap(connectiontypes.ErrConnectionNotFound, channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return sdkerrors.Wrapf(
			connectiontypes.ErrInvalidConnectionState,
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	commitment := k.GetPacketCommitment(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())

	packetCommitment := types.CommitPacket(k.cdc, packet)

	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(types.ErrInvalidPacket, "commitment bytes are not equal: got (%v), expected (%v)", packetCommitment, commitment)
	}

	if err := k.connectionKeeper.VerifyPacketAcknowledgement(
		ctx, connectionEnd, proofHeight, proof, packet.GetDestPort(), packet.GetDestChannel(),
		packet.GetSequence(), acknowledgement,
	); err != nil {
		return sdkerrors.Wrap(err, "packet acknowledgement verification failed")
	}

	// assert packets acknowledged in order
	if channel.Ordering == types.ORDERED {
		nextSequenceAck, found := k.GetNextSequenceAck(ctx, packet.GetSourcePort(), packet.GetSourceChannel())
		if !found {
			return sdkerrors.Wrapf(
				types.ErrSequenceAckNotFound,
				"source port: %s, source channel: %s", packet.GetSourcePort(), packet.GetSourceChannel(),
			)
		}

		if packet.GetSequence() != nextSequenceAck {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidSequence,
				"packet sequence ≠ next ack sequence (%d ≠ %d)", packet.GetSequence(), nextSequenceAck,
			)
		}

		// All verification complete, in the case of ORDERED channels we must increment nextSequenceAck
		nextSequenceAck++

		// incrementing NextSequenceAck and storing under this chain's channelEnd identifiers
		// Since this is the original sending chain, our channelEnd is packet's source port and channel
		k.SetNextSequenceAck(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), nextSequenceAck)

	}

	// Delete packet commitment, since the packet has been acknowledged, the commitement is no longer necessary
	k.deletePacketCommitment(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())

	// log that a packet has been acknowledged
	k.Logger(ctx).Info("packet acknowledged", "packet", fmt.Sprintf("%v", packet))

	// emit an event marking that we have processed the acknowledgement
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAcknowledgePacket,
			sdk.NewAttributeString(types.AttributeKeyTimeoutHeight, packet.GetTimeoutHeight().String()),
			sdk.NewAttributeString(types.AttributeKeyTimeoutTimestamp, fmt.Sprintf("%d", packet.GetTimeoutTimestamp())),
			sdk.NewAttributeString(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttributeString(types.AttributeKeySrcPort, packet.GetSourcePort()),
			sdk.NewAttributeString(types.AttributeKeySrcChannel, packet.GetSourceChannel()),
			sdk.NewAttributeString(types.AttributeKeyDstPort, packet.GetDestPort()),
			sdk.NewAttributeString(types.AttributeKeyDstChannel, packet.GetDestChannel()),
			sdk.NewAttributeString(types.AttributeKeyChannelOrdering, channel.Ordering.String()),
			// we only support 1-hop packets now, and that is the most important hop for a relayer
			// (is it going to a chain I am connected to)
			sdk.NewAttributeString(types.AttributeKeyConnection, channel.ConnectionHops[0]),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return nil
}
