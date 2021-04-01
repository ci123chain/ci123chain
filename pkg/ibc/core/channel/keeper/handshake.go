package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
)

// ChanOpenInit is called by a module to initiate a channel opening handshake with
// a module on another chain. The counterparty channel identifier is validated to be
// empty in msg validation.
func (k Keeper) ChanOpenInit(
	ctx sdk.Context,
	order types.Order,
	connectionHops []string,
	portID string,
	portCap *capabilitytypes.Capability,
	counterparty types.Counterparty,
	version string,
) (string, *capabilitytypes.Capability, error) {
	// connection hop length checked on msg.ValidateBasic()
	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, connectionHops[0])
	if !found {
		return "", nil, errors.Errorf("connection not found: %s", connectionHops[0])
	}

	getVersions := connectionEnd.GetVersions()
	if len(getVersions) != 1 {
		return "", nil, errors.Errorf(
			"single version must be negotiated on connection before opening channel, got: %v",
			getVersions,
		)
	}

	if !connectiontypes.VerifySupportedFeature(getVersions[0], order.String()) {
		return "", nil, errors.Errorf(
			"connection version %s does not support channel ordering: %s",
			getVersions[0], order.String(),
		)
	}

	if !k.portKeeper.Authenticate(ctx, portCap, portID) {
		return "", nil, errors.Errorf("caller does not own port capability for port ID %s", portID)
	}

	channelID := k.GenerateChannelIdentifier(ctx)
	channel := types.NewChannel(types.INIT, order, counterparty, connectionHops, version)
	k.SetChannel(ctx, portID, channelID, channel)

	capKey, err := k.scopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
	if err != nil {
		return "", nil, errors.Wrapf(err, "could not create channel capability for port ID %s and channel ID %s", portID, channelID)
	}

	k.SetNextSequenceSend(ctx, portID, channelID, 1)
	k.SetNextSequenceRecv(ctx, portID, channelID, 1)
	k.SetNextSequenceAck(ctx, portID, channelID, 1)

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", "NONE", "new-state", "INIT")

	//defer func() {
	//	telemetry.IncrCounter(1, "ibc", "channel", "open-init")
	//}()

	return channelID, capKey, nil
}



// CounterpartyHops returns the connection hops of the counterparty channel.
// The counterparty hops are stored in the inverse order as the channel's.
// NOTE: Since connectionHops only supports single connection channels for now,
// this function requires that connection hops only contain a single connection id
func (k Keeper) CounterpartyHops(ctx sdk.Context, ch types.Channel) ([]string, bool) {
	// Return empty array if connection hops is more than one
	// ConnectionHops length should be verified earlier
	if len(ch.ConnectionHops) != 1 {
		return []string{}, false
	}
	counterpartyHops := make([]string, 1)
	hop := ch.ConnectionHops[0]
	conn, found := k.connectionKeeper.GetConnection(ctx, hop)
	if !found {
		return []string{}, false
	}

	counterpartyHops[0] = conn.GetCounterparty().GetConnectionID()
	return counterpartyHops, true
}

// ChanOpenTry is called by a module to accept the first step of a channel opening
// handshake initiated by a module on another chain.
func (k Keeper) ChanOpenTry(
	ctx sdk.Context,
	order types.Order,
	connectionHops []string,
	portID,
	previousChannelID string,
	portCap *capabilitytypes.Capability,
	counterparty types.Counterparty,
	version,
	counterpartyVersion string,
	proofInit []byte,
	proofHeight exported.Height,
) (string, *capabilitytypes.Capability, error) {
	var (
		previousChannel      types.Channel
		previousChannelFound bool
	)

	channelID := previousChannelID

	// empty channel identifier indicates continuing a previous channel handshake
	if previousChannelID != "" {
		// channel identifier and connection hop length checked on msg.ValidateBasic()
		// ensure that the previous channel exists
		previousChannel, previousChannelFound = k.GetChannel(ctx, portID, previousChannelID)
		if !previousChannelFound {
			return "", nil, errors.Errorf("previous channel does not exist for supplied previous channelID %s", previousChannelID)
		}
		// previous channel must use the same fields
		if !(previousChannel.Ordering == order &&
			previousChannel.Counterparty.PortId == counterparty.PortId &&
			previousChannel.Counterparty.ChannelId == "" &&
			previousChannel.ConnectionHops[0] == connectionHops[0] &&
			previousChannel.Version == version) {
			return "", nil, errors.New("channel fields mismatch previous channel fields")
		}

		if previousChannel.State != types.INIT {
			return "", nil, errors.Errorf("previous channel state is in %s, expected INIT", previousChannel.State)
		}

	} else {
		// generate a new channel
		channelID = k.GenerateChannelIdentifier(ctx)
	}

	if !k.portKeeper.Authenticate(ctx, portCap, portID) {
		return "", nil, errors.Errorf("caller does not own port capability for port ID %s", portID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, connectionHops[0])
	if !found {
		return "", nil, errors.Errorf("connection not found", connectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return "", nil, errors.Errorf(
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	getVersions := connectionEnd.GetVersions()
	if len(getVersions) != 1 {
		return "", nil, errors.Errorf(
			"single version must be negotiated on connection before opening channel, got: %v",
			getVersions,
		)
	}

	if !connectiontypes.VerifySupportedFeature(getVersions[0], order.String()) {
		return "", nil, errors.Errorf(
			"connection version %s does not support channel ordering: %s",
			getVersions[0], order.String(),
		)
	}

	// NOTE: this step has been switched with the one below to reverse the connection
	// hops
	channel := types.NewChannel(types.TRYOPEN, order, counterparty, connectionHops, version)

	counterpartyHops, found := k.CounterpartyHops(ctx, channel)
	if !found {
		// should not reach here, connectionEnd was able to be retrieved above
		panic("cannot find connection")
	}

	// expectedCounterpaty is the counterparty of the counterparty's channel end
	// (i.e self)
	expectedCounterparty := types.NewCounterparty(portID, "")
	expectedChannel := types.NewChannel(
		types.INIT, channel.Ordering, expectedCounterparty,
		counterpartyHops, counterpartyVersion,
	)

	if err := k.connectionKeeper.VerifyChannelState(
		ctx, connectionEnd, proofHeight, proofInit,
		counterparty.PortId, counterparty.ChannelId, expectedChannel,
	); err != nil {
		return "", nil, err
	}

	var (
		capKey *capabilitytypes.Capability
		err    error
	)

	if !previousChannelFound {
		capKey, err = k.scopedKeeper.NewCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
		if err != nil {
			return "", nil, errors.Wrapf(err, "could not create channel capability for port ID %s and channel ID %s", portID, channelID)
		}

		k.SetNextSequenceSend(ctx, portID, channelID, 1)
		k.SetNextSequenceRecv(ctx, portID, channelID, 1)
		k.SetNextSequenceAck(ctx, portID, channelID, 1)
	} else {
		// capability initialized in ChanOpenInit
		capKey, found = k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(portID, channelID))
		if !found {
			return "", nil, errors.Errorf(
				"capability not found for existing channel, portID (%s) channelID (%s)", portID, channelID,
			)
		}
	}

	k.SetChannel(ctx, portID, channelID, channel)

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", previousChannel.State.String(), "new-state", "TRYOPEN")

	//defer func() {
	//	telemetry.IncrCounter(1, "ibc", "channel", "open-try")
	//}()

	return channelID, capKey, nil
}



// ChanOpenAck is called by the handshake-originating module to acknowledge the
// acceptance of the initial request by the counterparty module on the other chain.
func (k Keeper) ChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterpartyVersion,
	counterpartyChannelID string,
	proofTry []byte,
	proofHeight exported.Height,
) error {
	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return errors.Errorf("port ID (%s) channel ID (%s)", portID, channelID)
	}

	if !(channel.State == types.INIT || channel.State == types.TRYOPEN) {
		return errors.Errorf(
			"channel state should be INIT or TRYOPEN (got %s)", channel.State.String(),
		)
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return errors.Errorf("caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errors.Errorf("connection not found %s", channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errors.Errorf(
			"connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	counterpartyHops, found := k.CounterpartyHops(ctx, channel)
	if !found {
		// should not reach here, connectionEnd was able to be retrieved above
		panic("cannot find connection")
	}

	// counterparty of the counterparty channel end (i.e self)
	expectedCounterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.TRYOPEN, channel.Ordering, expectedCounterparty,
		counterpartyHops, counterpartyVersion,
	)

	if err := k.connectionKeeper.VerifyChannelState(
		ctx, connectionEnd, proofHeight, proofTry,
		channel.Counterparty.PortId, counterpartyChannelID,
		expectedChannel,
	); err != nil {
		return err
	}

	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", channel.State.String(), "new-state", "OPEN")

	//defer func() {
	//	telemetry.IncrCounter(1, "ibc", "channel", "open-ack")
	//}()

	channel.State = types.OPEN
	channel.Version = counterpartyVersion
	channel.Counterparty.ChannelId = counterpartyChannelID
	k.SetChannel(ctx, portID, channelID, channel)

	return nil
}



// ChanOpenConfirm is called by the counterparty module to close their end of the
//  channel, since the other end has been closed.
func (k Keeper) ChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	proofAck []byte,
	proofHeight exported.Height,
) error {
	channel, found := k.GetChannel(ctx, portID, channelID)
	if !found {
		return errors.Errorf("channel not found: port ID (%s) channel ID (%s)", portID, channelID)
	}

	if channel.State != types.TRYOPEN {
		return errors.Errorf(
			"invalid channel state, channel state is not TRYOPEN (got %s)", channel.State.String(),
		)
	}

	if !k.scopedKeeper.AuthenticateCapability(ctx, chanCap, host.ChannelCapabilityPath(portID, channelID)) {
		return errors.Errorf("channel capability not found: caller does not own capability for channel, port ID (%s) channel ID (%s)", portID, channelID)
	}

	connectionEnd, found := k.connectionKeeper.GetConnection(ctx, channel.ConnectionHops[0])
	if !found {
		return errors.Errorf("connection not found", channel.ConnectionHops[0])
	}

	if connectionEnd.GetState() != int32(connectiontypes.OPEN) {
		return errors.Errorf(
			"invalid connection state, connection state is not OPEN (got %s)", connectiontypes.State(connectionEnd.GetState()).String(),
		)
	}

	counterpartyHops, found := k.CounterpartyHops(ctx, channel)
	if !found {
		// Should not reach here, connectionEnd was able to be retrieved above
		panic("cannot find connection")
	}

	counterparty := types.NewCounterparty(portID, channelID)
	expectedChannel := types.NewChannel(
		types.OPEN, channel.Ordering, counterparty,
		counterpartyHops, channel.Version,
	)

	if err := k.connectionKeeper.VerifyChannelState(
		ctx, connectionEnd, proofHeight, proofAck,
		channel.Counterparty.PortId, channel.Counterparty.ChannelId,
		expectedChannel,
	); err != nil {
		return err
	}

	channel.State = types.OPEN
	k.SetChannel(ctx, portID, channelID, channel)
	k.Logger(ctx).Info("channel state updated", "port-id", portID, "channel-id", channelID, "previous-state", "TRYOPEN", "new-state", "OPEN")

	//defer func() {
	//	telemetry.IncrCounter(1, "ibc", "channel", "open-confirm")
	//}()
	return nil
}
