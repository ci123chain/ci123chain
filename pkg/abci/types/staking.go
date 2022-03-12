package types

import (
	"math/big"
)

type BondStatus int32

// PowerReduction is the amount of staking tokens required for 1 unit of consensus-engine power
var PowerReduction = NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))

func TokensToConsensusPower(tokens Int) int64 {
	return (tokens.Quo(PowerReduction)).Int64()
}

// TokensFromConsensusPower - convert input power to tokens
func TokensFromConsensusPower(power int64) Int {
	return NewInt(power).Mul(PowerReduction)
}

const (

	// default bond denomination
	DefaultBondDenom = "WLK"

	ValidatorUpdateDelay int64 = 1

	Unbonded  BondStatus = 1
	Unbonding BondStatus = 2
	Bonded    BondStatus = 3

	BondStatusUnbonded = "Unbonded"
	BondStatusUnbonding = "Unbonding"
	BondStatusBonded = "Bonded"
	BondStatusAll = "all"
)

var (
	ChainCoinDenom = DefaultBondDenom
)

func SetCoinDenom(denom string) {
	ChainCoinDenom = denom
}

// Equal compares two BondStatus instances
func (b BondStatus) Equal(b2 BondStatus) bool {
	return byte(b) == byte(b2)
}

// String implements the Stringer interface for BondStatus.
func (b BondStatus) String() string {
	switch b {
	case Unbonded:
		return BondStatusUnbonded

	case Unbonding:
		return BondStatusUnbonding

	case Bonded:
		return BondStatusBonded

	default:
		panic("invalid bond status")
	}
}