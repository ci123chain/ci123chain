package types

import (
	"fmt"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"regexp"
)

const SubModuleName = "connection"
const KeyNextConnectionSequence = "nextConnectionSequence"
// ConnectionPrefix is the prefix used when creating a connection identifier
const ConnectionPrefix = "connection-"


// FormatConnectionIdentifier returns the connection identifier with the sequence appended.
// This is a SDK specific format not enforced by IBC protocol.
func FormatConnectionIdentifier(sequence uint64) string {
	return fmt.Sprintf("%s%d", ConnectionPrefix, sequence)
}


// IsConnectionIDFormat checks if a connectionID is in the format required on the SDK for
// parsing connection identifiers. The connection identifier must be in the form: `connection-{N}
var IsConnectionIDFormat = regexp.MustCompile(`^connection-[0-9]{1,20}$`).MatchString

// IsValidConnectionID checks if the connection identifier is valid and can be parsed to
// the connection identifier format.
func IsValidConnectionID(connectionID string) bool {
	_, err := ParseConnectionSequence(connectionID)
	return err == nil
}

// ParseConnectionSequence parses the connection sequence from the connection identifier.
func ParseConnectionSequence(connectionID string) (uint64, error) {
	if !IsConnectionIDFormat(connectionID) {
		return 0, sdkerrors.Wrap(host.ErrInvalidID, "connection identifier is not in the format: `connection-{N}`")
	}

	sequence, err := host.ParseIdentifier(connectionID, ConnectionPrefix)
	if err != nil {
		return 0, sdkerrors.Wrap(err, "invalid connection identifier")
	}

	return sequence, nil
}