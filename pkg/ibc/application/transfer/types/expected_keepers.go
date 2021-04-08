package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	capabilitytypes "github.com/ci123chain/ci123chain/pkg/capability/types"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	ibcexported "github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	supplytypes "github.com/ci123chain/ci123chain/pkg/supply/exported"
)

// AccountKeeper defines the contract required for account APIs.
type AccountKeeper interface {

}


// BankKeeper defines the expected bank keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplytypes.ModuleAccountI
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coin) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coin) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coin) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coin) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coin) error
}


// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
	//ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error
}

// ClientKeeper defines the expected IBC client keeper
type ClientKeeper interface {
	GetClientConsensusState(ctx sdk.Context, clientID string) (connection ibcexported.ConsensusState, found bool)
}

// ConnectionKeeper defines the expected IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection connectiontypes.ConnectionEnd, found bool)
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
}