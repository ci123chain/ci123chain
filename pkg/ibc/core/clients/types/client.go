package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"math"
	"strings"
)
type IdentifiedClientState struct {
	// client identifier
	ClientId string `json:"client_id,omitempty" yaml:"client_id"`
	// client state
	ClientState exported.ClientState `json:"client_state,omitempty" yaml:"client_state"`
}

func NewIdentifiedClientState(clientID string, clientState exported.ClientState) IdentifiedClientState {
	return IdentifiedClientState{
		ClientState: clientState,
		ClientId: clientID,
	}
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
