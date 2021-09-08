package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

const (
	DefaultCodespace = "preStaking"
	RouteKey = DefaultCodespace
	ModuleName = DefaultCodespace
	StoreKey = "preStaking"

	PreStakingRecordQuery = "queryPreStakingRecord"
	StakingRecordQuery = "queryStakingRecord"
)


var (
	PreStakingKey = []byte{0x50}
)



func GetPreStakingKey(delegator sdk.AccAddress) []byte {
	return append(PreStakingKey, delegator.Bytes()...)
}