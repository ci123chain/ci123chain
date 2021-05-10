package types

import (
	"fmt"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/gogo/protobuf/proto"
	"math"
	"strings"
)

var (
	_ codectypes.UnpackInterfacesMessage = IdentifiedClientState{}
	_ codectypes.UnpackInterfacesMessage = ConsensusStateWithHeight{}
)

func NewIdentifiedClientState(clientID string, clientState exported.ClientState) IdentifiedClientState {
	msg, ok := clientState.(proto.Message)
	if !ok {
		panic(fmt.Errorf("cannot proto marshal %T", clientState))
	}

	anyClientState, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}
	return IdentifiedClientState{
		ClientState: anyClientState,
		ClientId: clientID,
	}
}


// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (ics IdentifiedClientState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(ics.ClientState, new(exported.ClientState))
}


// ValidateClientType validates the client types. It cannot be blank or empty. It must be a valid
// client identifier when used with '0' or the maximum uint64 as the sequence.
func ValidateClientType(clientType string) error {
	if strings.TrimSpace(clientType) == "" {
		return sdkerrors.Wrap(ErrInvalidClientType, "client type cannot be blank")
	}

	smallestPossibleClientID := FormatClientIdentifier(clientType, 0)
	largestPossibleClientID := FormatClientIdentifier(clientType, uint64(math.MaxUint64))

	// IsValidClientID will check client types format and if the sequence is a uint64
	if !IsValidClientID(smallestPossibleClientID) {
		return sdkerrors.Wrap(ErrInvalidClientType, "")
	}

	if err := host.ClientIdentifierValidator(smallestPossibleClientID); err != nil {
		return sdkerrors.Wrap(err, "client types results in smallest client identifier being invalid")
	}
	if err := host.ClientIdentifierValidator(largestPossibleClientID); err != nil {
		return sdkerrors.Wrap(err, "client types results in largest client identifier being invalid")
	}

	return nil
}

// NewConsensusStateWithHeight creates a new ConsensusStateWithHeight instance
func NewConsensusStateWithHeight(height Height, consensusState exported.ConsensusState) ConsensusStateWithHeight {
	msg, ok := consensusState.(proto.Message)
	if !ok {
		panic(fmt.Errorf("cannot proto marshal %T", consensusState))
	}

	anyConsensusState, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	return ConsensusStateWithHeight{
		Height:         height,
		ConsensusState: anyConsensusState,
	}
}

// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (cswh ConsensusStateWithHeight) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(cswh.ConsensusState, new(exported.ConsensusState))
}