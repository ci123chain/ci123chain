package types

import (
	"fmt"
	"time"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type GenesisState struct {
	// params defines all the paramaters of related to deposit.
	Params Params `json:"params"`
	// signing_infos represents a map between validator addresses and their
	// signing infos.
	SigningInfos []SigningInfo `json:"signing_infos" yaml:"signing_infos"`
	// signing_infos represents a map between validator addresses and their
	// missed blocks.
	MissedBlocks []ValidatorMissedBlocks `json:"missed_blocks" yaml:"missed_blocks"`
}

// SigningInfo stores validator signing info of corresponding address.
type SigningInfo struct {
	// address is the validator address.
	Address string `json:"address,omitempty"`
	// validator_signing_info represents the signing info of this validator.
	ValidatorSigningInfo ValidatorSigningInfo `json:"validator_signing_info" yaml:"validator_signing_info"`
}

// ValidatorMissedBlocks contains array of missed blocks of corresponding
// address.
type ValidatorMissedBlocks struct {
	// address is the validator address.
	Address string `json:"address,omitempty"`
	// missed_blocks is an array of missed blocks by the validator.
	MissedBlocks []MissedBlock `json:"missed_blocks" yaml:"missed_blocks"`
}

// MissedBlock contains height and missed status as boolean.
type MissedBlock struct {
	// index is the height at which the block was missed.
	Index int64 `json:"index,omitempty"`
	// missed is the missed status.
	Missed bool `json:"missed,omitempty"`
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, signingInfos []SigningInfo, missedBlocks []ValidatorMissedBlocks,
) *GenesisState {

	return &GenesisState{
		Params:       params,
		SigningInfos: signingInfos,
		MissedBlocks: missedBlocks,
	}
}

// NewMissedBlock creates a new MissedBlock instance
func NewMissedBlock(index int64, missed bool) MissedBlock {
	return MissedBlock{
		Index:  index,
		Missed: missed,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		SigningInfos: []SigningInfo{},
		MissedBlocks: []ValidatorMissedBlocks{},
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	downtime := data.Params.SlashFractionDowntime
	if downtime.IsNegative() || downtime.GT(sdk.OneDec()) {
		return fmt.Errorf("slashing fraction downtime should be less than or equal to one and greater than zero, is %s", downtime.String())
	}

	dblSign := data.Params.SlashFractionDoubleSign
	if dblSign.IsNegative() || dblSign.GT(sdk.OneDec()) {
		return fmt.Errorf("slashing fraction double sign should be less than or equal to one and greater than zero, is %s", dblSign.String())
	}

	minSign := data.Params.MinSignedPerWindow
	if minSign.IsNegative() || minSign.GT(sdk.OneDec()) {
		return fmt.Errorf("min signed per window should be less than or equal to one and greater than zero, is %s", minSign.String())
	}

	downtimeJail := data.Params.DowntimeJailDuration
	if downtimeJail < 1*time.Minute {
		return fmt.Errorf("downtime unblond duration must be at least 1 minute, is %s", downtimeJail.String())
	}

	signedWindow := data.Params.SignedBlocksWindow
	if signedWindow < 10 {
		return fmt.Errorf("signed blocks window must be at least 10, is %d", signedWindow)
	}

	return nil
}
