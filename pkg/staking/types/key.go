package types

import (
	"encoding/binary"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"strconv"
	"time"
)

const (
	RouteKey = "staking"
	StoreKey = "staking"
	ModuleName = "staking"

	QueryDelegation = "delegation"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidators = "validators"
	QueryValidator = "validator"
	QueryDelegatorValidators = "delegatorValidators"
	QueryDelegatorValidator = "delegatorValidator"
	QueryRedelegations = "redelegations"
	QueryDelegatorDelegations          = "delegatorDelegations"

	QueryOperatorAddressSet = "operatorAddresses"
)

var (
	LastTotalPowerKey = []byte{0x11}
	LastValidatorPowerKey = []byte{0x12}

	ValidatorsKey = []byte{0x21}
	ValidatorsByConsAddrKey = []byte{0x22}
	ValidatorsByPowerIndexKey = []byte{0x23}

	DelegationKey = []byte{0x31}
	RedelegationKey = []byte{0x32}
	RedelegationByValSrcIndexKey = []byte{0x33}
	RedelegationByValDstIndexKey = []byte{0x34}
	UnbondingDelegationKey  = []byte{0x35}
	UnbondingDelegationByValIndexKey = []byte{0x36}

	UnbondingQueueKey = []byte{0x41}
	RedelegationQueueKey = []byte{0x42}
	ValidatorQueueKey = []byte{0x43}

	HistoricalInfoKey = []byte{0x50}
)


func GetValidatorKey(operatorAddr sdk.AccAddress) []byte {
	return append(ValidatorsKey, operatorAddr.Bytes()...)
}

func GetValidatorByConsAddrKey(addr sdk.AccAddress) []byte {
	return append(ValidatorsByConsAddrKey, addr.Bytes()...)
}

func GetValidatorsByPowerIndexKey(validator Validator) []byte {
	return getValidatorPowerRank(validator)
}

func getValidatorPowerRank(validator Validator) []byte {

	consensusPower := sdk.TokensToConsensusPower(validator.Tokens)
	consensusPowerBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(consensusPowerBytes, uint64(consensusPower))

	powerBytes := consensusPowerBytes
	powerBytesLen := len(powerBytes) // 8

	// key is of format prefix || powerbytes || addrBytes
	key := make([]byte, 1+powerBytesLen+sdk.AddrLen)

	key[0] = ValidatorsByPowerIndexKey[0]
	copy(key[1:powerBytesLen+1], powerBytes)
	operAddrInvr := sdk.CopyBytes(validator.OperatorAddress.Bytes())
	for i, b := range operAddrInvr {
		operAddrInvr[i] = ^b
	}
	copy(key[powerBytesLen+1:], operAddrInvr)

	return key
}

func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.AccAddress) []byte {
	return append(GetDelegationsKey(delAddr), valAddr.Bytes()...)
}

func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, delAddr.Bytes()...)
}

// gets the prefix keyspace for redelegations from a delegator
func GetREDsKey(delAddr sdk.AccAddress) []byte {
	return append(RedelegationKey, delAddr.Bytes()...)
}

func GetREDKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.AccAddress) []byte {
	key := make([]byte, 1+sdk.AddrLen*3)
	copy(key[0:sdk.AddrLen+1], GetREDsKey(delAddr))
	copy(key[sdk.AddrLen+1:2*sdk.AddrLen+1], valSrcAddr.Bytes())
	copy(key[2*sdk.AddrLen+1:3*sdk.AddrLen+1], valDstAddr.Bytes())

	return key
}

// gets the index-key for a redelegation, stored by source-validator-index
// VALUE: none (key rearrangement used)
func GetREDByValSrcIndexKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.AccAddress) []byte {
	REDSFromValsSrcKey := GetREDsFromValSrcIndexKey(valSrcAddr)
	offset := len(REDSFromValsSrcKey)

	// key is of the form REDSFromValsSrcKey || delAddr || valDstAddr
	key := make([]byte, len(REDSFromValsSrcKey)+2*sdk.AddrLen)
	copy(key[0:offset], REDSFromValsSrcKey)
	copy(key[offset:offset+sdk.AddrLen], delAddr.Bytes())
	copy(key[offset+sdk.AddrLen:offset+2*sdk.AddrLen], valDstAddr.Bytes())
	return key
}

// gets the index-key for a redelegation, stored by destination-validator-index
// VALUE: none (key rearrangement used)
func GetREDByValDstIndexKey(delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.AccAddress) []byte {
	REDSToValsDstKey := GetREDsToValDstIndexKey(valDstAddr)
	offset := len(REDSToValsDstKey)

	// key is of the form REDSToValsDstKey || delAddr || valSrcAddr
	key := make([]byte, len(REDSToValsDstKey)+2*sdk.AddrLen)
	copy(key[0:offset], REDSToValsDstKey)
	copy(key[offset:offset+sdk.AddrLen], delAddr.Bytes())
	copy(key[offset+sdk.AddrLen:offset+2*sdk.AddrLen], valSrcAddr.Bytes())

	return key
}

// gets the prefix keyspace for all redelegations redelegating away from a source validator
func GetREDsFromValSrcIndexKey(valSrcAddr sdk.AccAddress) []byte {
	return append(RedelegationByValSrcIndexKey, valSrcAddr.Bytes()...)
}

// gets the prefix keyspace for all redelegations redelegating towards a destination validator
func GetREDsToValDstIndexKey(valDstAddr sdk.AccAddress) []byte {
	return append(RedelegationByValDstIndexKey, valDstAddr.Bytes()...)
}

// gets the prefix keyspace for all redelegations redelegating towards a destination validator
// from a particular delegator
func GetREDsByDelToValDstIndexKey(delAddr sdk.AccAddress, valDstAddr sdk.AccAddress) []byte {
	return append(
		GetREDsToValDstIndexKey(valDstAddr),
		delAddr.Bytes()...)
}

// gets the prefix for all unbonding delegations from a delegator
func GetUBDsKey(delAddr sdk.AccAddress) []byte {
	return append(UnbondingDelegationKey, delAddr.Bytes()...)
}

// gets the key for an unbonding delegation by delegator and validator addr
// VALUE: staking/UnbondingDelegation
func GetUBDKey(delAddr sdk.AccAddress, valAddr sdk.AccAddress) []byte {
	return append(
		GetUBDsKey(delAddr),
		valAddr.Bytes()...)
}

// gets the prefix keyspace for the indexes of unbonding delegations for a validator
func GetUBDsByValIndexKey(valAddr sdk.AccAddress) []byte {
	return append(UnbondingDelegationByValIndexKey, valAddr.Bytes()...)
}

// gets the index-key for an unbonding delegation, stored by validator-index
// VALUE: none (key rearrangement used)
func GetUBDByValIndexKey(delAddr sdk.AccAddress, valAddr sdk.AccAddress) []byte {
	return append(GetUBDsByValIndexKey(valAddr), delAddr.Bytes()...)
}

// gets the prefix for all unbonding delegations from a delegator
func GetUnbondingDelegationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}

// gets the prefix for all unbonding delegations from a delegator
func GetRedelegationTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(RedelegationQueueKey, bz...)
}

// get the bonded validator index key for an operator address
func GetLastValidatorPowerKey(operator sdk.AccAddress) []byte {
	return append(LastValidatorPowerKey, operator.Bytes()...)
}

// gets the prefix for all unbonding delegations from a delegator
func GetValidatorQueueTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(ValidatorQueueKey, bz...)
}

//________________________________________________________________________________

// GetHistoricalInfoKey gets the key for the historical info
func GetHistoricalInfoKey(height int64) []byte {
	return append(HistoricalInfoKey, []byte(strconv.FormatInt(height, 10))...)
}


// Get the validator operator address from LastValidatorPowerKey
func AddressFromLastValidatorPowerKey(key []byte) []byte {
	return key[1:] // remove prefix bytes
}

// GetREDKeyFromValSrcIndexKey rearranges the ValSrcIndexKey to get the REDKey
func GetREDKeyFromValSrcIndexKey(indexKey []byte) []byte {
	// note that first byte is prefix byte
	if len(indexKey) != 3*sdk.AddrLen+1 {
		panic("unexpected key length")
	}
	valSrcAddr := indexKey[1 : sdk.AddrLen+1]
	delAddr := indexKey[sdk.AddrLen+1 : 2*sdk.AddrLen+1]
	valDstAddr := indexKey[2*sdk.AddrLen+1 : 3*sdk.AddrLen+1]

	return GetREDKey(sdk.ToAccAddress(delAddr), sdk.ToAccAddress(valSrcAddr), sdk.ToAccAddress(valDstAddr))
}