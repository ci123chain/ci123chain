package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"time"
)

const (
	DefaultCodespace = "preStaking"
	RouteKey = DefaultCodespace
	ModuleName = DefaultCodespace
	StoreKey = "preStaking"

	PreStakingRecordQuery = "queryPreStakingRecord"
	StakingRecordQuery = "queryStakingRecord"
	PreStakingTokenQuery = "queryPreStakingToken"
)

var MinPreStakingTime = time.Minute * 1

var (
	PreStakingKey = []byte{0x50}
	StakingRecordKey = []byte{0x51}

	TokenManager = []byte("tokenManager")
	TokenManagerOwner = []byte("tokenManagerOwner")
)



func GetPreStakingKey(delegator sdk.AccAddress) []byte {
	return append(PreStakingKey, delegator.Bytes()...)
}

func GetStakingRecordKey(delegator, validator sdk.AccAddress) []byte {
	return append(StakingRecordKey, append(delegator.Bytes(), validator.Bytes()...)...)
}