package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/params"
)


// Parameter keys
var (
	ParamStoreKeyCommunityTax        = []byte("communitytax")
	ParamStoreKeyBaseProposerReward  = []byte("baseproposerreward")
	ParamStoreKeyBonusProposerReward = []byte("bonusproposerreward")
	ParamStoreKeyWithdrawAddrEnabled = []byte("withdrawaddrenabled")
)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func DefaultParams() Params {
	return Params{
		CommunityTax:        sdk.NewDecWithPrec(2, 2), // 2%
		BaseProposerReward:  sdk.NewDecWithPrec(1, 2), // 1%
		BonusProposerReward: sdk.NewDecWithPrec(4, 2), // 4%
		WithdrawAddrEnabled: true,
	}
}

func (p Params) ValidateBasic() error {
	if p.CommunityTax.IsNegative() || p.CommunityTax.GT(sdk.OneDec()) {
		return fmt.Errorf(
			"community tax should non-negative and less than one: %s", p.CommunityTax,
		)
	}
	if p.BaseProposerReward.IsNegative() {
		return fmt.Errorf(
			"base proposer reward should be positive: %s", p.BaseProposerReward,
		)
	}
	if p.BonusProposerReward.IsNegative() {
		return fmt.Errorf(
			"bonus proposer reward should be positive: %s", p.BonusProposerReward,
		)
	}
	if v := p.BaseProposerReward.Add(p.BonusProposerReward); v.GT(sdk.OneDec()) {
		return fmt.Errorf(
			"sum of base and bonus proposer reward cannot greater than one: %s", v,
		)
	}

	return nil
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamStoreKeyCommunityTax, &p.CommunityTax, validateCommunityTax),
		params.NewParamSetPair(ParamStoreKeyBaseProposerReward, &p.BaseProposerReward, validateBaseProposerReward),
		params.NewParamSetPair(ParamStoreKeyBonusProposerReward, &p.BonusProposerReward, validateBonusProposerReward),
		params.NewParamSetPair(ParamStoreKeyWithdrawAddrEnabled, &p.WithdrawAddrEnabled, validateWithdrawAddrEnabled),
	}
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter types: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("community tax must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("community tax must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("community tax too large: %s", v)
	}

	return nil
}

func validateBaseProposerReward(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter types: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("base proposer reward must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("base proposer reward must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("base proposer reward too large: %s", v)
	}

	return nil
}

func validateBonusProposerReward(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter types: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("bonus proposer reward must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("bonus proposer reward must be positive: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("bonus proposer reward too large: %s", v)
	}

	return nil
}

func validateWithdrawAddrEnabled(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter types: %T", i)
	}

	return nil
}