package types

import (
	"bytes"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	"sort"
	"time"
)

type Validator struct {
	OperatorAddress    sdk.AccAddress	`json:"operator_address"`
	Address            cmn.HexBytes     `json:"address"`
	ConsensusKey       crypto.PubKey    `json:"pub_key"`
	Jailed             bool             `json:"jailed"`
	Status             sdk.BondStatus   `json:"status"`
	Tokens             sdk.Int          `json:"tokens"`
	DelegatorShares    sdk.Dec          `json:"delegator_shares"`
	Description        Description      `json:"description"`
	UnbondingHeight    int64      		`json:"unbonding_height"`
	UnbondingTime      time.Time        `json:"unbonding_time"`
	Commission         Commission       `json:"commission"`
	MinSelfDelegation   sdk.Int  		`json:"min_self_delegation"`
}

//crypto pubKey
func NewValidator(operator sdk.AccAddress, pubKey crypto.PubKey, description Description) Validator {
	/*var pkStr string
	if pubKey != nil {
		pkStr = string(pubKey)
	}*/

	return Validator{
		OperatorAddress: operator,
		ConsensusKey: pubKey,
		Address: pubKey.Address(),
		Jailed:          false,
		Status:          sdk.Unbonded,
		Tokens:          sdk.ZeroInt(),
		DelegatorShares:  sdk.ZeroDec(),
		Description:      description,
		UnbondingHeight:  int64(0),
		UnbondingTime:    time.Unix(0,0).UTC(),
		Commission:        NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		MinSelfDelegation:  sdk.OneInt(),
	}
}

func (v Validator) SetInitialCommission(commission Commission) (Validator, error) {
	if err := commission.CommissionRates.Validate(); err != nil {
		return v, err
	}

	v.Commission = commission
	return v, nil
}

// IsBonded checks if the validator status equals Bonded
func (v Validator) IsBonded() bool {
	return v.GetStatus().Equal(sdk.Bonded)
}

// IsUnbonded checks if the validator status equals Unbonded
func (v Validator) IsUnbonded() bool {
	return v.GetStatus().Equal(sdk.Unbonded)
}

// IsUnbonding checks if the validator status equals Unbonding
func (v Validator) IsUnbonding() bool {
	return v.GetStatus().Equal(sdk.Unbonding)
}

// In some situations, the exchange rate becomes invalid, e.g. if
// Validator loses all tokens due to slashing. In this case,
// make all future delegations invalid.
func (v Validator) InvalidExRate() bool {
	return v.Tokens.IsZero() && v.DelegatorShares.IsPositive()
}

// AddTokensFromDel adds tokens to a validator
func (v Validator) AddTokensFromDel(amount sdk.Int) (Validator, sdk.Dec) {

	// calculate the shares to issue
	var issuedShares sdk.Dec
	if v.DelegatorShares.IsZero() {
		// the first delegation to a validator sets the exchange rate to one
		issuedShares = amount.ToDec()
	} else {
		shares, err := v.SharesFromTokens(amount)
		if err != nil {
			panic(err)
		}

		issuedShares = shares
	}

	//issuedShares = amount.ToDec()
	v.Tokens = v.Tokens.Add(amount)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, issuedShares
}

// RemoveDelShares removes delegator shares from a validator.
// NOTE: because token fractions are left in the valiadator,
//       the exchange rate of future shares of this validator can increase.
func (v Validator) RemoveDelShares(delShares sdk.Dec) (Validator, sdk.Int) {
	remainingShares := v.DelegatorShares.Sub(delShares)

	var issuedTokens sdk.Int
	if remainingShares.IsZero() {
		// last delegation share gets any trimmings
		issuedTokens = v.Tokens
		v.Tokens = sdk.ZeroInt()
	} else {
		// leave excess tokens in the validator
		// however fully use all the delegator shares
		issuedTokens = v.TokensFromShares(delShares).TruncateInt()
		v.Tokens = v.Tokens.Sub(issuedTokens)

		if v.Tokens.IsNegative() {
			panic("attempting to remove more tokens than available in validator")
		}
	}

	v.DelegatorShares = remainingShares
	return v, issuedTokens
}

// calculate the token worth of provided shares
func (v Validator) TokensFromShares(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).Quo(v.DelegatorShares)
}

// SharesFromTokens returns the shares of a delegation given a bond amount. It
// returns an error if the validator has no tokens.
func (v Validator) SharesFromTokens(amt sdk.Int) (sdk.Dec, error) {
	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoInt(v.GetTokens()), nil
}

func (v Validator) SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) {

	if v.Tokens.IsZero() {
		return sdk.ZeroDec(), ErrInsufficientShares
	}

	return v.GetDelegatorShares().MulInt(amt).QuoTruncate(v.GetTokens().ToDec()), nil
}

// get the bonded tokens which the validator holds
func (v Validator) BondedTokens() sdk.Int {
	if v.IsBonded() {
		return v.Tokens
	}
	return sdk.ZeroInt()
}

// get the consensus-engine power
// a reduction of 10^6 from validator tokens is applied
func (v Validator) ConsensusPower() int64 {
	if v.IsBonded() {
		return v.PotentialConsensusPower()
	}
	return 0
}

// UpdateStatus updates the location of the shares within a validator
// to reflect the new status
func (v Validator) UpdateStatus(newStatus sdk.BondStatus) Validator {
	v.Status = newStatus
	return v
}

// ABCIValidatorUpdateZero returns an abci.ValidatorUpdate from a staking validator type
// with zero power used for validator updates.
func (v Validator) ABCIValidatorUpdateZero() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.GetConsPubKey()),
		Power:  0,
	}
}

// ABCIValidatorUpdate returns an abci.ValidatorUpdate from a staking validator type
// with the full validator power
func (v Validator) ABCIValidatorUpdate() abci.ValidatorUpdate {
	return abci.ValidatorUpdate{
		PubKey: tmtypes.TM2PB.PubKey(v.GetConsPubKey()),
		Power:  v.ConsensusPower(),
	}
}

// potential consensus-engine power
func (v Validator) PotentialConsensusPower() int64 {
	return sdk.TokensToConsensusPower(v.Tokens)
}

// nolint - for ValidatorI
func (v Validator) IsJailed() bool              { return v.Jailed }
func (v Validator) GetMoniker() string          { return v.Description.Moniker }
func (v Validator) GetStatus() sdk.BondStatus   { return v.Status }
func (v Validator) GetOperator() sdk.AccAddress  { return v.OperatorAddress }

func (v Validator) GetConsAddress() sdk.AccAddress  { return sdk.ToAccAddress(v.GetConsPubKey().Address())}//公钥
func (v Validator) GetTokens() sdk.Int            { return v.Tokens }
func (v Validator) GetBondedTokens() sdk.Int      { return v.BondedTokens() }
func (v Validator) GetConsensusPower() int64      { return v.ConsensusPower() }
func (v Validator) GetCommission() sdk.Dec        { return v.Commission.CommissionRates.Rate }
func (v Validator) GetMinSelfDelegation() sdk.Int { return v.MinSelfDelegation }
func (v Validator) GetDelegatorShares() sdk.Dec   { return v.DelegatorShares }

func (v Validator) GetConsAddr() sdk.AccAddress {return sdk.ToAccAddress(v.GetConsPubKey().Address())}

func (v Validator) GetConsPubKey() crypto.PubKey {
	return v.ConsensusKey
}

// calculate the token worth of provided shares, truncated
func (v Validator) TokensFromSharesTruncated(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoTruncate(v.DelegatorShares)
}

// TokensFromSharesRoundUp returns the token worth of provided shares, rounded
// up.
func (v Validator) TokensFromSharesRoundUp(shares sdk.Dec) sdk.Dec {
	return (shares.MulInt(v.Tokens)).QuoRoundUp(v.DelegatorShares)
}


// Validators is a collection of Validator
type Validators []Validator

// Sort Validators sorts validator array in ascending operator address order
func (v Validators) Sort() {
	sort.Sort(v)
}

// Implements sort interface
func (v Validators) Len() int {
	return len(v)
}

// Implements sort interface
func (v Validators) Less(i, j int) bool {
	return bytes.Compare(v[i].OperatorAddress.Bytes(), v[j].OperatorAddress.Bytes()) == -1
}

// Implements sort interface
func (v Validators) Swap(i, j int) {
	it := v[i]
	v[i] = v[j]
	v[j] = it
}

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}