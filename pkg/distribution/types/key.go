package types


import (
	"encoding/binary"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

const (
	ModuleName  = "distribution"
	RouteKey = ModuleName
	DefaultParamspace = ModuleName
)

var (
	FeePoolKey                        = []byte{0x00}
	ProposerKey                       = []byte{001}
	ValidatorOutstandingRewardsPrefix = []byte{0x02} // key for outstanding rewards
	DelegatorWithdrawAddrPrefix          = []byte{0x03} // key for delegator withdraw address
	DelegatorStartingInfoPrefix          = []byte{0x04} // key for delegator starting info
	ValidatorHistoricalRewardsPrefix     = []byte{0x05} // key for historical validators rewards / stake
	ValidatorCurrentRewardsPrefix        = []byte{0x06} // key for current validator rewards
	ValidatorAccumulatedCommissionPrefix = []byte{0x07} // key for accumulated validator commission
	ValidatorSlashEventPrefix            = []byte{0x08} // key for validator slash fraction
)


func GetValidatorAccumulatedCommissionKey(v sdk.AccAddress) []byte {
	return append(ValidatorAccumulatedCommissionPrefix, v.Bytes()...)
}

func GetValidatorCurrentRewardsKey(v sdk.AccAddress) []byte {
	return append(ValidatorCurrentRewardsPrefix, v.Bytes()...)
}

// gets the outstanding rewards key for a validator
func GetValidatorOutstandingRewardsKey(valAddr sdk.AccAddress) []byte {
	return append(ValidatorOutstandingRewardsPrefix, valAddr.Bytes()...)
}
// gets the key for a delegator's withdraw addr
func GetDelegatorWithdrawAddrKey(delAddr sdk.AccAddress) []byte {
	return append(DelegatorWithdrawAddrPrefix, delAddr.Bytes()...)
}

// gets the prefix key for a validator's historical rewards
func GetValidatorHistoricalRewardsPrefix(v sdk.AccAddress) []byte {
	return append(ValidatorHistoricalRewardsPrefix, v.Bytes()...)
}

// gets the key for a validator's historical rewards
func GetValidatorHistoricalRewardsKey(v sdk.AccAddress, k uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, k)
	return append(append(ValidatorHistoricalRewardsPrefix, v.Bytes()...), b...)
}

// gets the key for a delegator's starting info
func GetDelegatorStartingInfoKey(v sdk.AccAddress, d sdk.AccAddress) []byte {
	return append(append(DelegatorStartingInfoPrefix, v.Bytes()...), d.Bytes()...)
}


// gets the prefix key for a validator's slash fraction (ValidatorSlashEventPrefix + height)
func GetValidatorSlashEventKeyPrefix(v sdk.AccAddress, height uint64) []byte {
	heightBz := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBz, height)
	return append(
		ValidatorSlashEventPrefix,
		append(v.Bytes(), heightBz...)...,
	)
}

// gets the key for a validator's slash fraction
func GetValidatorSlashEventKey(v sdk.AccAddress, height, period uint64) []byte {
	periodBz := make([]byte, 8)
	binary.BigEndian.PutUint64(periodBz, period)
	prefix := GetValidatorSlashEventKeyPrefix(v, height)
	return append(prefix, periodBz...)
}

// gets the prefix key for a validator's slash fractions
func GetValidatorSlashEventPrefix(v sdk.AccAddress) []byte {
	return append(ValidatorSlashEventPrefix, v.Bytes()...)
}