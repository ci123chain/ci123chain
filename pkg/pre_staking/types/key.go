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
	PreStakingDaoQuery = "queryPreStakingDao"
)


var (
	PreStakingKey = []byte{0x50}
	StakingRecordKey = []byte{0x51}

	WeeLinkDAO = []byte("weeLinkDAO")
)



func GetPreStakingKey(delegator sdk.AccAddress) []byte {
	return append(PreStakingKey, delegator.Bytes()...)
}

func GetStakingRecordKey(delegator, validator sdk.AccAddress) []byte {
	return append(StakingRecordKey, append(delegator.Bytes(), validator.Bytes()...)...)
}