package host

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"regexp"
	"strings"
)

// DefaultMaxCharacterLength defines the default maximum character length used
// in validation of identifiers including the client, connection, port and
// channel identifiers.
//
// NOTE: this restriction is specific to this golang implementation of IBC. If
// your use case demands a higher limit, please open an issue and we will consider
// adjusting this restriction.
const DefaultMaxCharacterLength = 64

var IsValidID = regexp.MustCompile(`^[a-zA-Z0-9\.\_\+\-\#\[\]\<\>]+$`).MatchString


func defaultIdentifierValidator(id string, min, max int) error { //nolint:unparam
	if strings.TrimSpace(id) == "" {
		return sdkerrors.Wrap(ErrInvalidID, "identifier cannot be blank")
	}
	// valid id MUST NOT contain "/" separator
	if strings.Contains(id, "/") {
		return sdkerrors.Wrapf(ErrInvalidID, "identifier %s cannot contain separator '/'", id)
	}
	// valid id must fit the length requirements
	if len(id) < min || len(id) > max {
		return sdkerrors.Wrapf(ErrInvalidID, "identifier %s has invalid length: %d, must be between %d-%d characters", id, len(id), min, max)
	}
	// valid id must contain only lower alphabetic characters
	if !IsValidID(id) {
		return sdkerrors.Wrapf(
			ErrInvalidID,
			"identifier %s must contain only alphanumeric or the following characters: '.', '_', '+', '-', '#', '[', ']', '<', '>'",
			id,
		)	}
	return nil
}

// ClientIdentifierValidator is the default validator function for Client identifiers.
// A valid Identifier must be between 9-64 characters and only contain alphanumeric and some allowed
// special characters (see IsValidID).
func ClientIdentifierValidator(id string) error {
	return defaultIdentifierValidator(id, 9, DefaultMaxCharacterLength)
}


// ConnectionIdentifierValidator is the default validator function for Connection identifiers.
// A valid Identifier must be between 10-64 characters and only contain alphanumeric and some allowed
// special characters (see IsValidID).
func ConnectionIdentifierValidator(id string) error {
	return defaultIdentifierValidator(id, 10, DefaultMaxCharacterLength)
}

// PortIdentifierValidator is the default validator function for Port identifiers.
// A valid Identifier must be between 2-64 characters and only contain alphanumeric and some allowed
// special characters (see IsValidID).
func PortIdentifierValidator(id string) error {
	return defaultIdentifierValidator(id, 2, DefaultMaxCharacterLength)
}

// ChannelIdentifierValidator is the default validator function for Channel identifiers.
// A valid Identifier must be between 8-64 characters and only contain alphanumeric and some allowed
// special characters (see IsValidID).
func ChannelIdentifierValidator(id string) error {
	return defaultIdentifierValidator(id, 8, DefaultMaxCharacterLength)
}
