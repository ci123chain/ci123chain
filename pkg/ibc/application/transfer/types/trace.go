package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
)

// ParseHexHash parses a hex hash in string format to bytes and validates its correctness.
func ParseHexHash(hexHash string) (cmn.HexBytes, error) {
	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return nil, err
	}

	if err := tmtypes.ValidateHash(hash); err != nil {
		return nil, err
	}

	return hash, nil
}
// GetFullDenomPath returns the full denomination according to the ICS20 specification:
// tracePath + "/" + baseDenom
// If there exists no trace then the base denomination is returned.
func (dt DenomTrace) GetFullDenomPath() string {
	if dt.Path == "" {
		return dt.BaseDenom
	}
	return dt.GetPrefix() + dt.BaseDenom
}


// Hash returns the hex bytes of the SHA256 hash of the DenomTrace fields using the following formula:
//
// hash = sha256(tracePath + "/" + baseDenom)
func (dt DenomTrace) Hash() cmn.HexBytes {
	hash := sha256.Sum256([]byte(dt.GetFullDenomPath()))
	return hash[:]
}

// GetPrefix returns the receiving denomination prefix composed by the trace info and a separator.
func (dt DenomTrace) GetPrefix() string {
	return dt.Path + "/"
}

// IBCDenom a coin denomination for an ICS20 fungible token in the format
// 'ibc/{hash(tracePath + baseDenom)}'. If the trace is empty, it will return the base denomination.
func (dt DenomTrace) IBCDenom() string {
	if dt.Path != "" {
		return fmt.Sprintf("%s/%s", DenomPrefix, dt.Hash())
	}
	return dt.BaseDenom
}



// ParseDenomTrace parses a string with the ibc prefix (denom trace) and the base denomination
// into a DenomTrace type.
//
// Examples:
//
// 	- "portidone/channelidone/uatom" => DenomTrace{Path: "portidone/channelidone", BaseDenom: "uatom"}
// 	- "uatom" => DenomTrace{Path: "", BaseDenom: "uatom"}
func ParseDenomTrace(rawDenom string) DenomTrace {
	denomSplit := strings.Split(rawDenom, "/")

	if denomSplit[0] == rawDenom {
		return DenomTrace{
			Path:      "",
			BaseDenom: rawDenom,
		}
	}

	return DenomTrace{
		Path:      strings.Join(denomSplit[:len(denomSplit)-1], "/"),
		BaseDenom: denomSplit[len(denomSplit)-1],
	}
}