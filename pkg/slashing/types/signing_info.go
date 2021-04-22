package types

import (
	"fmt"
	"time"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

// ValidatorSigningInfo defines a validator's signing info for monitoring their
// liveness activity.
type ValidatorSigningInfo struct {
	Address string `json:"address,omitempty"`
	// Height at which validator was first a candidate OR was unjailed
	StartHeight int64 `json:"start_height,omitempty" yaml:"start_height"`
	// Index which is incremented each time the validator was a bonded
	// in a block and may have signed a precommit or not. This in conjunction with the
	// `SignedBlocksWindow` param determines the index in the `MissedBlocksBitArray`.
	IndexOffset int64 `json:"index_offset,omitempty" yaml:"index_offset"`
	// Timestamp until which the validator is jailed due to liveness downtime.
	JailedUntil time.Time `json:"jailed_until" yaml:"jailed_until"`
	// Whether or not a validator has been tombstoned (killed out of validator set). It is set
	// once the validator commits an equivocation or for any other configured misbehiavor.
	Tombstoned bool `json:"tombstoned,omitempty"`
	// A counter kept to avoid unnecessary array reads.
	// Note that `Sum(MissedBlocksBitArray)` always equals `MissedBlocksCounter`.
	MissedBlocksCounter int64 `json:"missed_blocks_counter,omitempty" yaml:"missed_blocks_counter"`
}

// NewValidatorSigningInfo creates a new ValidatorSigningInfo instance
//nolint:interfacer
func NewValidatorSigningInfo(
	condAddr sdk.AccAddress, startHeight, indexOffset int64,
	jailedUntil time.Time, tombstoned bool, missedBlocksCounter int64,
) ValidatorSigningInfo {

	return ValidatorSigningInfo{
		Address:             condAddr.String(),
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
	}
}

// String implements the stringer interface for ValidatorSigningInfo
func (i ValidatorSigningInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Index Offset:          %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d`,
		i.Address, i.StartHeight, i.IndexOffset, i.JailedUntil,
		i.Tombstoned, i.MissedBlocksCounter)
}

// unmarshal a validator signing info from a store value
func UnmarshalValSigningInfo(cdc *codec.Codec, value []byte) (signingInfo ValidatorSigningInfo, err error) {
	err = cdc.UnmarshalBinaryBare(value, &signingInfo)
	return signingInfo, err
}
