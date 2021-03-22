package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

// MustUnmarshalClientState attempts to decode and return an ClientState object from
// raw encoded bytes. It panics on error.
func (k Keeper) MustUnmarshalClientState(bz []byte) exported.ClientState {
	return types.MustUnmarshalClientState(k.cdc, bz)
}

// MustMarshalClientState attempts to encode an ClientState object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalClientState(clientState exported.ClientState) []byte {
	return types.MustMarshalClientState(k.cdc, clientState)
}


func (k Keeper) MustMarshalConsensusState(consensusState exported.ConsensusState) []byte {
	return types.MustMarshalConsensusState(k.cdc, consensusState)
}