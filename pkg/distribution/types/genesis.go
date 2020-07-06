package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Params struct {
	CommunityTax         sdk.Dec    `json:"community_tax"`
	BaseProposerReward   sdk.Dec    `json:"base_proposer_reward"`
	BonusProposerReward  sdk.Dec    `json:"bounse_proposer_reward"`
	WithdrawAddrEnabled   bool       `json:"withdraw_addr_enable"`
}

type DelegatorWithdrawInfo struct {
	DelegatorAddress     sdk.AccAddress    `json:"delegator_address"`
	WithdrawAddress      sdk.AccAddress    `json:"withdraw_address"`
}

type ValidatorOutstandingRewardsRecord struct {
	ValidatorAddress     sdk.AccAddress    `json:"validator_address"`
	OutstandingRewards   sdk.DecCoin       `json:"outstanding_rewards"`
}

type ValidatorAccumulatedCommissionRecord struct {
	ValidatorAddress     sdk.AccAddress    `json:"validator_address"`
	Accumulated          ValidatorAccumulatedCommission   `json:"accumulated"`
}


type ValidatorHistoricalRewardsRecord struct {
	ValidatorAddress    sdk.AccAddress      `json:"validator_address"`
	Period              uint64              `json:"period"`
	Rewards             ValidatorHistoricalRewards  `json:"rewards"`
}

type ValidatorCurrentRewardsRecord struct {
	ValidatorAddress    sdk.AccAddress      `json:"validator_address"`
	Rewards             ValidatorCurrentRewards  `json:"rewards"`
}

type DelegatorStartingInfoRecord struct {
	DelegatorAddress    sdk.AccAddress      `json:"delegator_address"`
	ValidatorAddress    sdk.AccAddress      `json:"validator_address"`
	StartingInfo        DelegatorStartingInfo  `json:"starting_info"`
}

type ValidatorSlashEventRecord struct {
	ValidatorAddress    sdk.AccAddress       `json:"validator_address"`
	Height              uint64               `json:"height"`
	Period              uint64               `json:"period"`
	Event               ValidatorSlashEvent  `json:"event"`
}


type GenesisState struct {
	Params        Params										   `json:"params"`
	FeePool       FeePool           							   `json:"fee_pool"`
	DelegatorWithdrawInfos []DelegatorWithdrawInfo    			   `json:"delegator_withdraw_infos"`
	PreviousProposer       sdk.AccAddr             			   	   `json:"previous_proposer"`
	OutstandingRewards     []ValidatorOutstandingRewardsRecord     `json:"outstanding_rewards"`
	ValidatorAccumulatedCommissions  []ValidatorAccumulatedCommissionRecord   `json:"validator_accumulated_commissions"`
	ValidatorHistoricalRewards    []ValidatorHistoricalRewardsRecord          `json:"validator_historical_rewards"`
	ValidatorCurrentRewards       []ValidatorCurrentRewardsRecord             `json:"validator_current_rewards"`
	DelegatorStartingInfos        []DelegatorStartingInfoRecord               `json:"delegator_starting_infos"`
	ValidatorSlashEvents          []ValidatorSlashEventRecord                 `json:"validator_slash_events"`
}

func NewGenesisState(
	params Params, fp FeePool, dwis []DelegatorWithdrawInfo, pp sdk.AccAddr, r []ValidatorOutstandingRewardsRecord,
	acc []ValidatorAccumulatedCommissionRecord, historical []ValidatorHistoricalRewardsRecord,
	cur []ValidatorCurrentRewardsRecord, dels []DelegatorStartingInfoRecord, slashes []ValidatorSlashEventRecord,
) GenesisState {

	return GenesisState{
		Params:                          params,
		FeePool:                         fp,
		DelegatorWithdrawInfos:          dwis,
		PreviousProposer:                pp,
		OutstandingRewards:              r,
		ValidatorAccumulatedCommissions: acc,
		ValidatorHistoricalRewards:      historical,
		ValidatorCurrentRewards:         cur,
		DelegatorStartingInfos:          dels,
		ValidatorSlashEvents:            slashes,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState(validators []tmtypes.GenesisValidator, accAddresses []string) GenesisState {
	var OutstandingRewards []ValidatorOutstandingRewardsRecord
	var CurrentRewards []ValidatorCurrentRewardsRecord
	for i, val := range validators {
		if val.PubKey != nil {
			addValOutstandingRewards := ValidatorOutstandingRewardsRecord{
				ValidatorAddress:  sdk.HexToAddress(accAddresses[i]),
				OutstandingRewards: sdk.NewDecCoin(sdk.DefaultCoinDenom, sdk.NewInt(0)),
			}
			currentReward := ValidatorCurrentRewardsRecord{
				ValidatorAddress:  sdk.HexToAddress(accAddresses[i]),
				Rewards:ValidatorCurrentRewards{
					Rewards: sdk.NewDecCoin(sdk.DefaultCoinDenom, sdk.NewInt(0)),
					Period:  0,
				},
			}
			OutstandingRewards = append(OutstandingRewards, addValOutstandingRewards)
			CurrentRewards = append(CurrentRewards, currentReward)
		}
	}

	return GenesisState{
		FeePool:                         InitialFeePool(),
		Params:                          DefaultParams(),
		DelegatorWithdrawInfos:          []DelegatorWithdrawInfo{},
		PreviousProposer:                nil,
		OutstandingRewards:              OutstandingRewards,//[]ValidatorOutstandingRewardsRecord{},
		ValidatorAccumulatedCommissions: []ValidatorAccumulatedCommissionRecord{},
		ValidatorHistoricalRewards:      []ValidatorHistoricalRewardsRecord{},
		ValidatorCurrentRewards:         CurrentRewards,//[]ValidatorCurrentRewardsRecord{},
		DelegatorStartingInfos:          []DelegatorStartingInfoRecord{},
		ValidatorSlashEvents:            []ValidatorSlashEventRecord{},
	}
}

func ValidateGenesis(gs GenesisState) error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}
	return gs.FeePool.ValidateGenesis()
}