package staking

import (
	h "github.com/ci123chain/ci123chain/pkg/staking/handler"
	k "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

const (
	RouteKey = types.RouteKey
	StoreKey = types.StoreKey
	ModuleName = types.ModuleName
	DefaultCodespace = types.DefaultCodespace

	/*QueryDelegation = "delegation"
	QueryAllDelegation = "allDelegation"
	QueryValidators = "validators"
	QueryValidator = "validator"
	QueryDelegatorValidators = "delegator_validators"
	QueryDelegatorValidator = "delegator_validator"
	QueryRedelegations = "redelegations"*/
)

var (
	ModuleCdc = types.StakingCodec
	NewHandler = h.NewHandler
	NewKeeper = k.NewStakingKeeper
	NewQuerier = k.NewQuerier
	KeyBondDenom = types.KeyBondDenom
	DefaultGenesisState                = types.DefaultGenesisState

	NewCreateValidatorMsg = types.NewCreateValidatorTx
	NewDelegateMsg = types.NewDelegateTx
	NewRedelegateMsg = types.NewRedelegateTx
	NewUndelegateMsg = types.NewUndelegateTx

	NewMultiStakingHooks  = types.NewMultiStakingHooks
)