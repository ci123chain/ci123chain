package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	sdkerrors "github.com/pkg/errors"
)

// TimeoutPacket is called by a module which originally attempted to send a
// packet to a counterparty module, where the timeout height has passed on the
// counterparty chain without the packet being committed, to prove that the
// packet can no longer be executed and to allow the calling module to safely
// perform appropriate state transitions. Its intended usage is within the
// ante handler.
func (k Keeper) TimeoutPacket(
	ctx sdk.Context,
	packet exported.PacketI,
	proof []byte,
	proofHeight exported.Height,
	nextSequenceRecv uint64,
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

	// NOTE: TimeoutPacket is called by the AnteHandler which acts upon the packet.Route(),
	// so the capability authentication can be omitted here

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
		return sdkerrors.Wrap(
			connectiontypes.ErrConnectionNotFound,
			channel.ConnectionHops[0],
		)
	}

	// check that timeout height or timeout timestamp has passed on the other end
	proofTimestamp, err := k.connectionKeeper.GetTimestampAtHeight(ctx, connectionEnd, proofHeight)
	if err != nil {
		return err
	}

	timeoutHeight := packet.GetTimeoutHeight()
	if (timeoutHeight.IsZero() || proofHeight.LT(timeoutHeight)) &&
		(packet.GetTimeoutTimestamp() == 0 || proofTimestamp < packet.GetTimeoutTimestamp()) {
		return sdkerrors.Wrap(types.ErrPacketTimeout, "packet timeout has not been reached for height or timestamp")
	}

	commitment := k.GetPacketCommitment(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())

	packetCommitment := types.CommitPacket(k.cdc, packet)

	x1, _ := json.Marshal(commitment)
	x2, _ := json.Marshal(packetCommitment)
	k.Logger(ctx).Info("x1: ", x1, ";x2: ", x2)
	// verify we sent the packet and haven't cleared it out yet
	if !bytes.Equal(commitment, packetCommitment) {
		return sdkerrors.Wrapf(types.ErrInvalidPacket, "packet commitment bytes are not equal: got (%v), expected (%v)", commitment, packetCommitment)
	}

	switch channel.Ordering {
	case types.ORDERED:
		// check that packet has not been received
		if nextSequenceRecv > packet.GetSequence() {
			return sdkerrors.Wrapf(
				types.ErrInvalidPacket,
				"packet already received, next sequence receive > packet sequence (%d > %d)", nextSequenceRecv, packet.GetSequence(),
			)
		}

		// check that the recv sequence is as claimed
		err = k.connectionKeeper.VerifyNextSequenceRecv(
			ctx, connectionEnd, proofHeight, proof,
			packet.GetDestPort(), packet.GetDestChannel(), nextSequenceRecv,
		)
	case types.UNORDERED:
		err = k.connectionKeeper.VerifyPacketReceiptAbsence(
			ctx, connectionEnd, proofHeight, proof,
			packet.GetDestPort(), packet.GetDestChannel(), packet.GetSequence(),
		)
	default:
		panic(sdkerrors.Wrapf(types.ErrInvalidChannelOrdering, channel.Ordering.String()))
	}

	if err != nil {
		return err
	}

	// NOTE: the remaining code is located in the TimeoutExecuted function
	return nil
}

// TimeoutExecuted deletes the commitment send from this chain after it verifies timeout.
// If the timed-out packet came from an ORDERED channel then this channel will be closed.
//
// CONTRACT: this function must be called in the IBC handler
func (k Keeper) TimeoutExecuted(
	ctx sdk.Context,
	chanCap *capabilitytypes.Capability,
	packet exported.PacketI,
) error {
	channel, found := k.GetChannel(ctx, packet.GetSourcePort(), packet.GetSourceChannel())
	if !found {
		return sdkerrors.Wrapf(types.ErrChannelNotFound, "port ID (%s) channel ID (%s)", packet.GetSourcePort(), packet.GetSourceChannel())
	}

	capName := host.ChannelCapabilityPath(packet.GetSourcePort(), packet.GetSourceChannel())
	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, capName) {
		return sdkerrors.Wrapf(
			types.ErrChannelCapabilityNotFound,
			"caller does not own capability for channel with capability name %s", capName,
		)
	}

	k.deletePacketCommitment(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())

	if channel.Ordering == types.ORDERED {
		channel.State = types.CLOSED
		k.SetChannel(ctx, packet.GetSourcePort(), packet.GetSourceChannel(), channel)
	}

	k.Logger(ctx).Info("packet timed-out", "packet", fmt.Sprintf("%v", packet))

	// emit an event marking that we have processed the timeout
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTimeoutPacket,
			sdk.NewAttributeString(types.AttributeKeyTimeoutHeight, packet.GetTimeoutHeight().String()),
			sdk.NewAttributeString(types.AttributeKeyTimeoutTimestamp, fmt.Sprintf("%d", packet.GetTimeoutTimestamp())),
			sdk.NewAttributeString(types.AttributeKeySequence, fmt.Sprintf("%d", packet.GetSequence())),
			sdk.NewAttributeString(types.AttributeKeySrcPort, packet.GetSourcePort()),
			sdk.NewAttributeString(types.AttributeKeySrcChannel, packet.GetSourceChannel()),
			sdk.NewAttributeString(types.AttributeKeyDstPort, packet.GetDestPort()),
			sdk.NewAttributeString(types.AttributeKeyDstChannel, packet.GetDestChannel()),
			sdk.NewAttributeString(types.AttributeKeyChannelOrdering, channel.Ordering.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	})

	return nil
}
