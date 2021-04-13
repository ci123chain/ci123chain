package types

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"sort"
)

// GenesisState defines the ibc client submodule's genesis state.
type GenesisState struct {
	// client states with their corresponding identifiers
	Clients IdentifiedClientStates `protobuf:"bytes,1,rep,name=clients,proto3,castrepeated=IdentifiedClientStates" json:"clients"`
	// consensus states from each client
	ClientsConsensus ClientsConsensusStates `protobuf:"bytes,2,rep,name=clients_consensus,json=clientsConsensus,proto3,castrepeated=ClientsConsensusStates" json:"clients_consensus" yaml:"clients_consensus"`
	// metadata from each client
	//ClientsMetadata []IdentifiedGenesisMetadata `protobuf:"bytes,3,rep,name=clients_metadata,json=clientsMetadata,proto3" json:"clients_metadata" yaml:"clients_metadata"`
	Params          Params                      `protobuf:"bytes,4,opt,name=params,proto3" json:"params"`
	// create localhost on initialization
	CreateLocalhost bool `protobuf:"varint,5,opt,name=create_localhost,json=createLocalhost,proto3" json:"create_localhost,omitempty" yaml:"create_localhost"`
	// the sequence for the next generated client identifier
	NextClientSequence uint64 `protobuf:"varint,6,opt,name=next_client_sequence,json=nextClientSequence,proto3" json:"next_client_sequence,omitempty" yaml:"next_client_sequence"`
}


// ConsensusStateWithHeight defines a consensus state with an additional height field.
type ConsensusStateWithHeight struct {
	// consensus state height
	Height Height `protobuf:"bytes,1,opt,name=height,proto3" json:"height"`
	// consensus state
	ConsensusState exported.ConsensusState `protobuf:"bytes,2,opt,name=consensus_state,json=consensusState,proto3" json:"consensus_state,omitempty" yaml"consensus_state"`
}

// ClientConsensusStates defines all the stored consensus states for a given
// client.
type ClientConsensusStates struct {
	// client identifier
	ClientId string `protobuf:"bytes,1,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty" yaml:"client_id"`
	// consensus states and their heights associated with the client
	ConsensusStates []ConsensusStateWithHeight `protobuf:"bytes,2,rep,name=consensus_states,json=consensusStates,proto3" json:"consensus_states" yaml:"consensus_states"`
}
// ClientsConsensusStates defines a slice of ClientConsensusStates that supports the sort interface
type ClientsConsensusStates []ClientConsensusStates

// Len implements sort.Interface
func (ccs ClientsConsensusStates) Len() int { return len(ccs) }

// Less implements sort.Interface
func (ccs ClientsConsensusStates) Less(i, j int) bool { return ccs[i].ClientId < ccs[j].ClientId }

// Swap implements sort.Interface
func (ccs ClientsConsensusStates) Swap(i, j int) { ccs[i], ccs[j] = ccs[j], ccs[i] }

// Sort is a helper function to sort the set of ClientsConsensusStates in place
func (ccs ClientsConsensusStates) Sort() ClientsConsensusStates {
	sort.Sort(ccs)
	return ccs
}


func DefaultGenesisState() GenesisState {
	return GenesisState{
		Clients:            []IdentifiedClientState{},
		ClientsConsensus:   ClientsConsensusStates{},
		Params:             DefaultParams(),
		CreateLocalhost:    false,
		NextClientSequence: 0,
	}
}