package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
)

// ParseHexHash parses a hex hash in string format to bytes and validates its correctness.
func ParseHexHash(hexHash string) (tmbytes.HexBytes, error) {
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
func (dt DenomTrace) Hash() tmbytes.HexBytes {
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


// ValidateIBCDenom validates that the given denomination is either:
//
//  - A valid base denomination (eg: 'uatom')
//  - A valid fungible token representation (i.e 'ibc/{hash}') per ADR 001 https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-001-coin-source-tracing.md
func ValidateIBCDenom(denom string) error {
	if err := sdk.ValidateDenom(denom); err != nil {
		return err
	}

	denomSplit := strings.SplitN(denom, "/", 2)

	switch {
	case strings.TrimSpace(denom) == "",
		len(denomSplit) == 1 && denomSplit[0] == DenomPrefix,
		len(denomSplit) == 2 && (denomSplit[0] != DenomPrefix || strings.TrimSpace(denomSplit[1]) == ""):
		return sdkerrors.Wrapf(ErrInvalidDenomForTransfer, "denomination should be prefixed with the format 'ibc/{hash(trace + \"/\" + %s)}'", denom)

	case denomSplit[0] == denom && strings.TrimSpace(denom) != "":
		return nil
	}

	if _, err := ParseHexHash(denomSplit[1]); err != nil {
		return sdkerrors.Wrapf(err, "invalid denom trace hash %s", denomSplit[1])
	}

	return nil
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