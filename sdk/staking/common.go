package staking

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

func CreateParseArgs(selfDelegation int64, rate, maxRate, maxChangeRate int64) (sdk.Int, sdk.Dec, sdk.Dec, sdk.Dec) {
	minSelfDelegation := sdk.NewInt(selfDelegation)
	r := sdk.NewDecWithPrec(rate, 2)
	mr := sdk.NewDecWithPrec(maxRate, 2)
	mxr := sdk.NewDecWithPrec(maxChangeRate, 2)

	return minSelfDelegation, r, mr, mxr
}