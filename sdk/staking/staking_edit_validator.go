package staking

import (
	"encoding/hex"
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
)

func SignEditValidator(from string, gas, nonce uint64, priv,
	moniker, identity, website, secu, details string, minSelfDelegation, newRate int64) ([]byte, error) {

	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	fromAddr := sdk.HexToAddress(from)
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

	tx := staking.NewEditValidatorMsg(fromAddr, gas, nonce, desc, nrArg, minArg)

	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}
