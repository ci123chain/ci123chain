package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

// MustUnmarshalClientState attempts to decode and return an ClientState object from
// raw encoded bytes. It panics on error.
func MustUnmarshalClientState(cdc codec.BinaryMarshaler, bz []byte) exported.ClientState {
	clientState, err := UnmarshalClientState(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to decode client state: %w", err))
	}

	return clientState
}


func UnmarshalClientState(cdc codec.BinaryMarshaler, bz []byte) (exported.ClientState, error) {
	var clientState exported.ClientState
	if err := cdc.UnmarshalInterface(bz, &clientState); err != nil {
		return nil, err
	}
	return clientState, nil
}


// MustMarshalClientState attempts to encode an ClientState object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalClientState(cdc codec.BinaryMarshaler, clientState exported.ClientState) []byte {
	bz, err := MarshalClientState(cdc, clientState)
	if err != nil {
		panic(fmt.Errorf("failed to encode client state: %w", err))
	}

	return bz
}

// MarshalClientState protobuf serializes an ClientState interface
func MarshalClientState(cdc codec.BinaryMarshaler, clientStateI exported.ClientState) ([]byte, error) {
	return cdc.MarshalInterface(clientStateI)
}


// MustMarshalConsensusState attempts to encode an ConsensusState object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalConsensusState(cdc codec.BinaryMarshaler, consensusState exported.ConsensusState) []byte {
	bz, err := MarshalConsensusState(cdc, consensusState)
	if err != nil {
		panic(fmt.Errorf("failed to encode consensus state: %w", err))
	}

	return bz
}
// MustUnmarshalConsensusState attempts to decode and return an ConsensusState object from
// raw encoded bytes. It panics on error.
func MustUnmarshalConsensusState(cdc codec.BinaryMarshaler, bz []byte) exported.ConsensusState {
	consensusState, err := UnmarshalConsensusState(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to decode consensus state: %w", err))
	}

	return consensusState
}

// MarshalConsensusState protobuf serializes an ConsensusState interface
func MarshalConsensusState(cdc codec.BinaryMarshaler, cs exported.ConsensusState) ([]byte, error) {
	return cdc.MarshalInterface(cs)
}

// UnmarshalConsensusState returns an ConsensusState interface from raw encoded clientState
// bytes of a Proto-based ConsensusState types. An error is returned upon decoding
// failure.
func UnmarshalConsensusState(cdc codec.BinaryMarshaler, bz []byte) (exported.ConsensusState, error) {
	var consensusState exported.ConsensusState
	if err := cdc.UnmarshalInterface(bz, &consensusState); err != nil {
		return nil, err
	}

	return consensusState, nil
}

func MustMarshalClientStateResp(cdc *codec.Codec, resp QueryClientStatesResponse) []byte {
	bz, err := MarshalClientStatesResp(cdc, resp)
	if err != nil {
		panic(fmt.Errorf("failed to decode clientstates response: %w", err))
	}
	return bz
}
func MarshalClientStatesResp(cdc *codec.Codec, resp QueryClientStatesResponse) ([]byte, error) {
	return cdc.MarshalJSON(resp)
}

func MustUnmarshalClientStateResp(cdc *codec.Codec, bz []byte) QueryClientStatesResponse {
	resp, err := UnmarshalClientStateResp(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to encode clientstates response: %w", err))
	}
	return resp
}
func UnmarshalClientStateResp(cdc *codec.Codec, bz []byte) (QueryClientStatesResponse, error) {
	var req QueryClientStatesResponse
	err := cdc.UnmarshalJSON(bz, &req)
	return req, err
}