package types

import (
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
)

// GenesisState defines the ibc module's genesis state.
type GenesisState struct {
	// ICS002 - Clients genesis state
	ClientGenesis clienttypes.GenesisState `protobuf:"bytes,1,opt,name=client_genesis,json=clientGenesis,proto3" json:"client_genesis" yaml:"client_genesis"`
	// ICS003 - Connections genesis state
	ConnectionGenesis connectiontypes.GenesisState `protobuf:"bytes,2,opt,name=connection_genesis,json=connectionGenesis,proto3" json:"connection_genesis" yaml:"connection_genesis"`
	// ICS004 - Channel genesis state
	ChannelGenesis channeltypes.GenesisState `protobuf:"bytes,3,opt,name=channel_genesis,json=channelGenesis,proto3" json:"channel_genesis" yaml:"channel_genesis"`
}

// DefaultGenesisState returns the ibc module's default genesis state.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		ClientGenesis:     clienttypes.DefaultGenesisState(),
		ConnectionGenesis: connectiontypes.DefaultGenesisState(),
		ChannelGenesis:    channeltypes.DefaultGenesisState(),
	}
}
