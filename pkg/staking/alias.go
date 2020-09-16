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

)

var (
	ModuleCdc = types.StakingCodec
	NewHandler = h.NewHandler
	NewKeeper = k.NewStakingKeeper
	NewQuerier = k.NewQuerier
	KeyBondDenom = types.KeyBondDenom
	DefaultGenesisState                = types.DefaultGenesisState

	NewCreateValidatorMsg = types.NewMsgCreateValidator
	NewEditValidatorMsg = types.NewMsgEditValidator
	NewDelegateMsg = types.NewMsgDelegate
	NewRedelegateMsg = types.NewMsgRedelegate
	NewUndelegateMsg = types.NewMsgUndelegate

	NewMultiStakingHooks  = types.NewMultiStakingHooks
)