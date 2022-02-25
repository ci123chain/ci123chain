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

	StakingRecordQuery = "queryStakingRecord"
	PreStakingTokenQuery = "queryPreStakingToken"
)

var MinPreStakingTime = time.Minute * 1

var (
	PreStakingIDKey = []byte{0x50}
	StakingRecordKey = []byte{0x51}

	TokenManager = []byte("tokenManager")
	TokenManagerOwner = []byte("tokenManagerOwner")
)



func GetStakingRecordKeyByID(id uint64) []byte {
	bz := sdk.Uint64ToBigEndian(id)
	return append(StakingRecordKey, bz...)
}


func GetStakingRecordID() []byte {
	return PreStakingIDKey
}