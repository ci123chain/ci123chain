package staking

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
)


func CommonParseArgs(from string, amount int64, validatorAddr, delegatorAddr string) (sdk.AccAddress, sdk.Coin, sdk.AccAddress, sdk.AccAddress, error) {

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{}, errors.New("unexpected from address")
	}
	amt := sdk.NewCoin(sdk.NewInt(amount))

	validatorAddress, err := helper.StrToAddress(validatorAddr)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{}, errors.New("unexpected validator address")
	}
	delegatorAddress, err := helper.StrToAddress(delegatorAddr)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{}, errors.New("unexpected delegator address")
	}

	return fromAddr, amt, validatorAddress, delegatorAddress, nil

}

func RedelegateParseArgs(from string, amount int64, validatorSrcAddr, validatorDstAddr, delegatorAddr string) (sdk.AccAddress, sdk.Coin, sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, error) {

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{},sdk.AccAddress{}, errors.New("unexpected from address")
	}
	amt := sdk.NewCoin(sdk.NewInt(amount))

	validatorSrcAddress, err := helper.StrToAddress(validatorSrcAddr)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{},sdk.AccAddress{}, errors.New("unexpected validator address")
	}
	validatorDstAddress, err := helper.StrToAddress(validatorDstAddr)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{},sdk.AccAddress{}, errors.New("unexpected validator address")
	}
	delegatorAddress, err := helper.StrToAddress(delegatorAddr)
	if err != nil {
		return sdk.AccAddress{},sdk.Coin{}, sdk.AccAddress{}, sdk.AccAddress{},sdk.AccAddress{}, errors.New("unexpected delegator address")
	}

	return fromAddr, amt, validatorSrcAddress, validatorDstAddress, delegatorAddress, nil

}

func CreateParseArgs(selfDelegation int64, rate, maxRate, maxChangeRate int64) (sdk.Int, sdk.Dec, sdk.Dec, sdk.Dec) {
	minSelfDelegation := sdk.NewInt(selfDelegation)
	r := sdk.NewDec(rate)
	mr := sdk.NewDec(maxRate)
	mxr := sdk.NewDec(maxChangeRate)

	return minSelfDelegation, r, mr, mxr
}