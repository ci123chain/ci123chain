package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/pkg/errors"
)

// VerifyConnectionState verifies a proof of the connection state of the
// specified connection end stored on the target machine.
func (k Keeper) VerifyConnectionState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	connectionID string,
	connectionEnd exported.ConnectionI, // opposite connection
) error {
	clientState, found := k.clientKeeper.GetClientState(ctx, connection.GetClientID())
	if !found {
		return errors.Errorf("client not found for client_id %s", connection.GetClientID())
	}

	if err := clientState.VerifyConnectionState(
		k.clientKeeper.ClientStore(ctx, connection.GetClientID()), k.cdc, height,
		connection.GetCounterparty().GetPrefix(), proof, connectionID, connectionEnd,
	); err != nil {
		return errors.Errorf("failed connection state verification for client (%s)", connection.GetClientID())
	}

	return nil
}


func (k Keeper) VerifyChannelState(ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID, channelID string,
	channel exported.ChannelI) error {
	clientState, found := k.clientKeeper.GetClientState(ctx, connection.GetClientID())
	if !found {
		return errors.Errorf("client state not found in verify channel state: %s", connection.GetClientID())
	}
	if err := clientState.VerifyChannelState(k.clientKeeper.ClientStore(ctx, connection.GetClientID()), k.cdc, height,
		connection.GetCounterparty().GetPrefix(), proof, portID, channelID, channel); err != nil {
		return errors.Wrapf(err, "failed channel state verification for client (%s)", connection.GetClientID())
	}
	return nil
}


// VerifyClientState verifies a proof of a client state of the running machine
// stored on the target machine
func (k Keeper) VerifyClientState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	clientState exported.ClientState,
) error {
	clientID := connection.GetClientID()
	targetClient, found := k.clientKeeper.GetClientState(ctx, clientID)
	if !found {
		return errors.Errorf("client not found for client_id %s", connection.GetClientID())
	}

	if err := targetClient.VerifyClientState(
		k.clientKeeper.ClientStore(ctx, clientID), k.cdc, height,
		connection.GetCounterparty().GetPrefix(), connection.GetCounterparty().GetClientID(), proof, clientState); err != nil {
		return errors.Wrapf(err, "failed client state verification for target client: %s", connection.GetClientID())
	}

	return nil
}



// VerifyClientConsensusState verifies a proof of the consensus state of the
// specified client stored on the target machine.
func (k Keeper) VerifyClientConsensusState(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	consensusHeight exported.Height,
	proof []byte,
	consensusState exported.ConsensusState,
) error {
	clientID := connection.GetClientID()
	clientState, found := k.clientKeeper.GetClientState(ctx, clientID)
	if !found {
		return errors.Errorf("client not found for client_id %s", connection.GetClientID())
	}

	if err := clientState.VerifyClientConsensusState(
		k.clientKeeper.ClientStore(ctx, clientID), k.cdc, height,
		connection.GetCounterparty().GetClientID(), consensusHeight, connection.GetCounterparty().GetPrefix(), proof, consensusState,
	); err != nil {
		return errors.Wrapf(err, "failed consensus state verification for client (%s)", connection.GetClientID())
	}

	return nil
}


// VerifyPacketAcknowledgement verifies a proof of an incoming packet
// acknowledgement at the specified port, specified channel, and specified sequence.
func (k Keeper) VerifyPacketAcknowledgement(
	ctx sdk.Context,
	connection exported.ConnectionI,
	height exported.Height,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	acknowledgement []byte,
) error {
	clientState, found := k.clientKeeper.GetClientState(ctx, connection.GetClientID())
	if !found {
		return errors.Errorf("client not found for client_id %s", connection.GetClientID())
	}

	if err := clientState.VerifyPacketAcknowledgement(
		k.clientKeeper.ClientStore(ctx, connection.GetClientID()), k.cdc, height,
		uint64(ctx.BlockHeader().Time.UnixNano()), connection.GetDelayPeriod(),
		connection.GetCounterparty().GetPrefix(), proof, portID, channelID,
		sequence, acknowledgement,
	); err != nil {
		return errors.Wrapf(err, "failed packet acknowledgement verification for client (%s)", connection.GetClientID())
	}

	return nil
}