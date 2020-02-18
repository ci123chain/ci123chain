package staking

import (
	"github.com/tanhuiya/ci123chain/pkg/staking/client/rest"
	h "github.com/tanhuiya/ci123chain/pkg/staking/handler"
	k "github.com/tanhuiya/ci123chain/pkg/staking/keeper"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
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
	RegisterRoutes = rest.RegisterTxRoutes
	KeyBondDenom = types.KeyBondDenom
	DefaultGenesisState                = types.DefaultGenesisState
)