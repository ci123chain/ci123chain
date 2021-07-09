package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	capabilitykeeper "github.com/ci123chain/ci123chain/pkg/capability/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/channel"
	channelkeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/keeper"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clientkeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/keeper"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectionkeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/keeper"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	portkeeper "github.com/ci123chain/ci123chain/pkg/ibc/core/port/keeper"
	porttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/port/types"
	paramtypes "github.com/ci123chain/ci123chain/pkg/params/types"
	"github.com/pkg/errors"
)

type Keeper struct {
	cdc codec.BinaryMarshaler
	ClientKeeper     clientkeeper.Keeper
	ConnectionKeeper connectionkeeper.Keeper
	ChannelKeeper    channelkeeper.Keeper
	PortKeeper       portkeeper.Keeper
	Router           *porttypes.Router
}

func NewKeeper(cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper,scopedKeeper capabilitykeeper.ScopedKeeper,) Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	connectionKeeper := connectionkeeper.NewKeeper(cdc, key, clientKeeper)

	portKeeper := portkeeper.NewKeeper(scopedKeeper)
	channelKeeper := channelkeeper.NewKeeper(cdc, key, clientKeeper, connectionKeeper, portKeeper, scopedKeeper)

	return Keeper{
		cdc:              cdc,
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
	}
}

// SetRouter sets the Router in IBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *porttypes.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}

func (k Keeper) CreateClient(ctx sdk.Context, msg *clienttypes.MsgCreateClient) (*clienttypes.MsgCreateClientResponse, error) {
	clientState, err := clienttypes.UnpackClientState(msg.ClientState)
	if err != nil {
		return nil, err
	}

	consensusState, err := clienttypes.UnpackConsensusState(msg.ConsensusState)
	if err != nil {
		return nil, err
	}

	clientID, err := k.ClientKeeper.CreateClient(ctx, clientState, consensusState)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			clienttypes.EventTypeCreateClient,
			sdk.NewAttributeString(clienttypes.AttributeKeyClientID, clientID),
			sdk.NewAttributeString(clienttypes.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttributeString(clienttypes.AttributeKeyConsensusHeight, clientState.GetLatestHeight().String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, clienttypes.AttributeValueCategory),
		),
	})
	return &clienttypes.MsgCreateClientResponse{}, nil
}


// UpdateClient defines a rpc handler method for MsgUpdateClient.
func (k Keeper) UpdateClient(ctx sdk.Context, msg *clienttypes.MsgUpdateClient) (*clienttypes.MsgUpdateClientResponse, error) {

	header, err := clienttypes.UnpackHeader(msg.Header)
	if err != nil {
		return nil, err
	}

	if err := k.ClientKeeper.UpdateClient(ctx, msg.ClientId, header); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, clienttypes.AttributeValueCategory),
		),
	)

	return &clienttypes.MsgUpdateClientResponse{}, nil
}




func (k Keeper) ConnectionOpenInit(ctx sdk.Context ,
	msg *connectiontypes.MsgConnectionOpenInit,
	) (*connectiontypes.MsgConnectionOpenInitResponse, error) {

	connectionID, err := k.ConnectionKeeper.ConnOpenInit(ctx, msg.ClientId, msg.Counterparty, msg.Version, msg.DelayPeriod)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "connection handshake open init failed")
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			connectiontypes.EventTypeConnectionOpenInit,
			sdk.NewAttributeString(connectiontypes.AttributeKeyConnectionID, connectionID),
			sdk.NewAttributeString(connectiontypes.AttributeKeyClientID, msg.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyClientID, msg.Counterparty.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyConnectionID, msg.Counterparty.ConnectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, connectiontypes.AttributeValueCategory),
		),
	})

	return &connectiontypes.MsgConnectionOpenInitResponse{}, nil
}


// ConnectionOpenTry defines a rpc handler method for MsgConnectionOpenTry.
func (k Keeper) ConnectionOpenTry(ctx sdk.Context, msg *connectiontypes.MsgConnectionOpenTry) (*connectiontypes.MsgConnectionOpenTryResponse, error) {
	targetClient, err := clienttypes.UnpackClientState(msg.ClientState)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "client in msg is not exported.ClientState. invalid client: %v.", targetClient)
	}

	connectionID, err := k.ConnectionKeeper.ConnOpenTry(
		ctx, msg.PreviousConnectionId, msg.Counterparty, msg.DelayPeriod, msg.ClientId, targetClient,
		connectiontypes.VersionsToExported(msg.CounterpartyVersions), msg.ProofInit, msg.ProofClient, msg.ProofConsensus,
		msg.ProofHeight, msg.ConsensusHeight,
	)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "connection handshake open try failed")
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			connectiontypes.EventTypeConnectionOpenTry,
			sdk.NewAttributeString(connectiontypes.AttributeKeyConnectionID, connectionID),
			sdk.NewAttributeString(connectiontypes.AttributeKeyClientID, msg.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyClientID, msg.Counterparty.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyConnectionID, msg.Counterparty.ConnectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, connectiontypes.AttributeValueCategory),
		),
	})

	return &connectiontypes.MsgConnectionOpenTryResponse{}, nil
}


// ConnectionOpenAck defines a rpc handler method for MsgConnectionOpenAck.
func (k Keeper) ConnectionOpenAck(ctx sdk.Context, msg *connectiontypes.MsgConnectionOpenAck) (*connectiontypes.MsgConnectionOpenAckResponse, error) {
	targetClient, err := clienttypes.UnpackClientState(msg.ClientState)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "client in msg is not exported.ClientState. invalid client: %v", targetClient)
	}

	if err := k.ConnectionKeeper.ConnOpenAck(
		ctx, msg.ConnectionId, targetClient, msg.Version, msg.CounterpartyConnectionId,
		msg.ProofTry, msg.ProofClient, msg.ProofConsensus,
		msg.ProofHeight, msg.ConsensusHeight,
	); err != nil {
		return nil, sdkerrors.Wrap(err, "connection handshake open ack failed")
	}

	connectionEnd, _ := k.ConnectionKeeper.GetConnection(ctx, msg.ConnectionId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			connectiontypes.EventTypeConnectionOpenAck,
			sdk.NewAttributeString(connectiontypes.AttributeKeyConnectionID, msg.ConnectionId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyClientID, connectionEnd.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyClientID, connectionEnd.Counterparty.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyConnectionID, connectionEnd.Counterparty.ConnectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, connectiontypes.AttributeValueCategory),
		),
	})

	return &connectiontypes.MsgConnectionOpenAckResponse{}, nil
}




// ConnectionOpenConfirm defines a rpc handler method for MsgConnectionOpenConfirm.
func (k Keeper) ConnectionOpenConfirm(ctx sdk.Context, msg *connectiontypes.MsgConnectionOpenConfirm) (*connectiontypes.MsgConnectionOpenConfirmResponse, error) {

	if err := k.ConnectionKeeper.ConnOpenConfirm(
		ctx, msg.ConnectionId, msg.ProofAck, msg.ProofHeight,
	); err != nil {
		return nil, sdkerrors.Wrap(err, "connection handshake open confirm failed")
	}

	connectionEnd, _ := k.ConnectionKeeper.GetConnection(ctx, msg.ConnectionId)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			connectiontypes.EventTypeConnectionOpenConfirm,
			sdk.NewAttributeString(connectiontypes.AttributeKeyConnectionID, msg.ConnectionId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyClientID, connectionEnd.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyClientID, connectionEnd.Counterparty.ClientId),
			sdk.NewAttributeString(connectiontypes.AttributeKeyCounterpartyConnectionID, connectionEnd.Counterparty.ConnectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttributeString(sdk.AttributeKeyModule, connectiontypes.AttributeValueCategory),
		),
	})

	return &connectiontypes.MsgConnectionOpenConfirmResponse{}, nil
}





// ChannelOpenInit defines a rpc handler method for MsgChannelOpenInit.
func (k Keeper) ChannelOpenInit(ctx sdk.Context, msg *channeltypes.MsgChannelOpenInit) (*channeltypes.MsgChannelOpenInitResponse, error) {

	// Lookup module by port capability
	module, portCap, err := k.PortKeeper.LookupModuleByPort(ctx, msg.PortId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	_, channelID, cap, err := channel.HandleMsgChannelOpenInit(ctx, k.ChannelKeeper, portCap, msg)
	if err != nil {
		return nil, err
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	if err = cbs.OnChanOpenInit(ctx, msg.Channel.Ordering, msg.Channel.ConnectionHops, msg.PortId, channelID, cap, msg.Channel.Counterparty, msg.Channel.Version); err != nil {
		return nil, errors.Wrap(err, "channel open init callback failed")
	}

	return &channeltypes.MsgChannelOpenInitResponse{}, nil
}

// ChannelOpenTry defines a rpc handler method for MsgChannelOpenTry.
func (k Keeper) ChannelOpenTry(ctx sdk.Context, msg *channeltypes.MsgChannelOpenTry) (*channeltypes.MsgChannelOpenTryResponse, error) {
	// Lookup module by port capability
	module, portCap, err := k.PortKeeper.LookupModuleByPort(ctx, msg.PortId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	_, channelID, cap, err := channel.HandleMsgChannelOpenTry(ctx, k.ChannelKeeper, portCap, msg)
	if err != nil {
		return nil, err
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	if err = cbs.OnChanOpenTry(ctx, msg.Channel.Ordering, msg.Channel.ConnectionHops, msg.PortId, channelID, cap, msg.Channel.Counterparty, msg.Channel.Version, msg.CounterpartyVersion); err != nil {
		return nil, sdkerrors.Wrap(err, "channel open try callback failed")
	}

	return &channeltypes.MsgChannelOpenTryResponse{}, nil
}

// ChannelOpenAck defines a rpc handler method for MsgChannelOpenAck.
func (k Keeper) ChannelOpenAck(ctx sdk.Context, msg *channeltypes.MsgChannelOpenAck) (*channeltypes.MsgChannelOpenAckResponse, error) {
	// Lookup module by channel capability
	module, cap, err := k.ChannelKeeper.LookupModuleByChannel(ctx, msg.PortId, msg.ChannelId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	_, err = channel.HandleMsgChannelOpenAck(ctx, k.ChannelKeeper, cap, msg)
	if err != nil {
		return nil, err
	}

	if err = cbs.OnChanOpenAck(ctx, msg.PortId, msg.ChannelId, msg.CounterpartyVersion); err != nil {
		return nil, sdkerrors.Wrap(err, "channel open ack callback failed")
	}

	return &channeltypes.MsgChannelOpenAckResponse{}, nil
}

// ChannelOpenConfirm defines a rpc handler method for MsgChannelOpenConfirm.
func (k Keeper) ChannelOpenConfirm(ctx sdk.Context, msg *channeltypes.MsgChannelOpenConfirm) (*channeltypes.MsgChannelOpenConfirmResponse, error) {
	// Lookup module by channel capability
	module, cap, err := k.ChannelKeeper.LookupModuleByChannel(ctx, msg.PortId, msg.ChannelId)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	_, err = channel.HandleMsgChannelOpenConfirm(ctx, k.ChannelKeeper, cap, msg)
	if err != nil {
		return nil, err
	}

	if err = cbs.OnChanOpenConfirm(ctx, msg.PortId, msg.ChannelId); err != nil {
		return nil, sdkerrors.Wrap(err, "channel open confirm callback failed")
	}

	return &channeltypes.MsgChannelOpenConfirmResponse{}, nil
}




// RecvPacket defines a rpc handler method for MsgRecvPacket.
func (k Keeper) RecvPacket(ctx sdk.Context, msg *channeltypes.MsgRecvPacket) (*channeltypes.MsgRecvPacketResponse, error) {
	// Lookup module by channel capability
	module, cap, err := k.ChannelKeeper.LookupModuleByChannel(ctx, msg.Packet.DestinationPort, msg.Packet.DestinationChannel)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	// Perform TAO verification
	if err := k.ChannelKeeper.RecvPacket(ctx, cap, msg.Packet, msg.ProofCommitment, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(err, "receive packet verification failed")
	}

	// Perform application logic callback
	_, ack, err := cbs.OnRecvPacket(ctx, msg.Packet)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "receive packet callback failed")
	}

	// Set packet acknowledgement only if the acknowledgement is not nil.
	// NOTE: IBC applications modules may call the WriteAcknowledgement asynchronously if the
	// acknowledgement is nil.
	if ack != nil {
		if err := k.ChannelKeeper.WriteAcknowledgement(ctx, cap, msg.Packet, ack); err != nil {
			return nil, err
		}
	}

	//defer func() {
	//	telemetry.IncrCounterWithLabels(
	//		[]string{"tx", "msg", "ibc", msg.Type()},
	//		1,
	//		[]metrics.Label{
	//			telemetry.NewLabel("source-port", msg.Packet.SourcePort),
	//			telemetry.NewLabel("source-channel", msg.Packet.SourceChannel),
	//			telemetry.NewLabel("destination-port", msg.Packet.DestinationPort),
	//			telemetry.NewLabel("destination-channel", msg.Packet.DestinationChannel),
	//		},
	//	)
	//}()

	return &channeltypes.MsgRecvPacketResponse{}, nil
}



// Timeout defines a rpc handler method for MsgTimeout.
func (k Keeper) Timeout(ctx sdk.Context, msg *channeltypes.MsgTimeout) (*channeltypes.MsgTimeoutResponse, error) {
	// Lookup module by channel capability
	module, cap, err := k.ChannelKeeper.LookupModuleByChannel(ctx, msg.Packet.SourcePort, msg.Packet.SourceChannel)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	// Perform TAO verification
	if err := k.ChannelKeeper.TimeoutPacket(ctx, msg.Packet, msg.ProofUnreceived, msg.ProofHeight, msg.NextSequenceRecv); err != nil {
		return nil, sdkerrors.Wrap(err, "timeout packet verification failed")
	}

	// Perform application logic callback
	_, err = cbs.OnTimeoutPacket(ctx, msg.Packet)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "timeout packet callback failed")
	}

	// Delete packet commitment
	if err = k.ChannelKeeper.TimeoutExecuted(ctx, cap, msg.Packet); err != nil {
		return nil, err
	}

	//defer func() {
	//	telemetry.IncrCounterWithLabels(
	//		[]string{"ibc", "timeout", "packet"},
	//		1,
	//		[]metrics.Label{
	//			telemetry.NewLabel("source-port", msg.Packet.SourcePort),
	//			telemetry.NewLabel("source-channel", msg.Packet.SourceChannel),
	//			telemetry.NewLabel("destination-port", msg.Packet.DestinationPort),
	//			telemetry.NewLabel("destination-channel", msg.Packet.DestinationChannel),
	//			telemetry.NewLabel("timeout-type", "height"),
	//		},
	//	)
	//}()

	return &channeltypes.MsgTimeoutResponse{}, nil
}



// Acknowledgement defines a rpc handler method for MsgAcknowledgement.
func (k Keeper) Acknowledgement(ctx sdk.Context, msg *channeltypes.MsgAcknowledgement) (*channeltypes.MsgAcknowledgementResponse, error) {

	// Lookup module by channel capability
	module, cap, err := k.ChannelKeeper.LookupModuleByChannel(ctx, msg.Packet.SourcePort, msg.Packet.SourceChannel)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "could not retrieve module from port-id")
	}

	// Retrieve callbacks from router
	cbs, ok := k.Router.GetRoute(module)
	if !ok {
		return nil, sdkerrors.Wrapf(porttypes.ErrInvalidRoute, "route not found to module: %s", module)
	}

	// Perform TAO verification
	if err := k.ChannelKeeper.AcknowledgePacket(ctx, cap, msg.Packet, msg.Acknowledgement, msg.ProofAcked, msg.ProofHeight); err != nil {
		return nil, sdkerrors.Wrap(err, "acknowledge packet verification failed")
	}

	// Perform application logic callback
	_, err = cbs.OnAcknowledgementPacket(ctx, msg.Packet, msg.Acknowledgement)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "acknowledge packet callback failed")
	}

	//defer func() {
	//	telemetry.IncrCounterWithLabels(
	//		[]string{"tx", "msg", "ibc", msg.Type()},
	//		1,
	//		[]metrics.Label{
	//			telemetry.NewLabel("source-port", msg.Packet.SourcePort),
	//			telemetry.NewLabel("source-channel", msg.Packet.SourceChannel),
	//			telemetry.NewLabel("destination-port", msg.Packet.DestinationPort),
	//			telemetry.NewLabel("destination-channel", msg.Packet.DestinationChannel),
	//		},
	//	)
	//}()

	return &channeltypes.MsgAcknowledgementResponse{}, nil
}

