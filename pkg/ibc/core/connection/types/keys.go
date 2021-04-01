package types

import "fmt"

const SubModuleName = "connection"
const KeyNextConnectionSequence = "nextConnectionSequence"
// ConnectionPrefix is the prefix used when creating a connection identifier
const ConnectionPrefix = "connection-"


// FormatConnectionIdentifier returns the connection identifier with the sequence appended.
// This is a SDK specific format not enforced by IBC protocol.
func FormatConnectionIdentifier(sequence uint64) string {
	return fmt.Sprintf("%s%d", ConnectionPrefix, sequence)
}