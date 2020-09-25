package staking

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

func SignEditValidator(from sdk.AccAddress, gas, nonce uint64, priv, moniker, identity, website, secu, details string,
	minSelfDelegation, newRate int64) ([]byte, error) {
	var nrArg *sdk.Dec
	var minArg *sdk.Int
	if newRate < 0 {
		if newRate == -1 {
			nrArg = nil
		}else {
			return nil, errors.New("invalid newRate")
		}
	}else {
		nr := sdk.NewDecWithPrec(newRate, 2)
		nrArg = &nr
	}
	if minSelfDelegation < 0 {
		if minSelfDelegation == -1 {
			minArg = nil
		}else {
			return nil, errors.New("invalid minSelfDelegation")
		}
	}else {
		min := sdk.NewInt(minSelfDelegation)
		minArg = &min
	}
	desc := types.Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: secu,
		Details:         details,
	}

	msg := staking.NewEditValidatorMsg(from, desc, nrArg, minArg)

	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}

	return txByte, nil
}

func NewEditValidatorMsg(from sdk.AccAddress, desc types.Description, nrArg *sdk.Dec, minArg *sdk.Int) []byte {
	msg := staking.NewEditValidatorMsg(from, desc, nrArg, minArg)
	return msg.Bytes()
}