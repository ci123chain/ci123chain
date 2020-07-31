package types

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"time"
)

const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280

	// constant used in flags to indicate that description field should not be updated
	DoNotModifyDesc = "[do-not-modify]"
)

//Rate:   佣金收取率；0-100;
//maxRate:  佣金最大收取率；
//maxChangeRate: 每日最大变动百分比
type CommissionRates struct {
	Rate           sdk.Dec     `json:"rate"`
	MaxRate        sdk.Dec     `json:"max_rate"`
	MaxChangeRate  sdk.Dec     `json:"max_change_rate"`
}

func NewCommissionRates(rate, maxRate, maxChangeRate sdk.Dec) CommissionRates {
	return CommissionRates{
		Rate:          rate,
		MaxRate:       maxRate,
		MaxChangeRate: maxChangeRate,
	}
}
func (cr *CommissionRates) Reset()      { *cr = CommissionRates{} }


func (cr CommissionRates) Validate() error {
	switch {
	case cr.MaxRate.IsNegative():
		return ErrCommissionNegative
	case cr.MaxRate.GT(sdk.OneDec()):
		// max rate cannot be greater than 1
		return ErrCommissionHuge

	case cr.Rate.IsNegative():
		// rate cannot be negative
		return ErrCommissionNegative

	case cr.Rate.GT(cr.MaxRate):
		// rate cannot be greater than the max rate
		return ErrCommissionGTMaxRate

	case cr.MaxChangeRate.IsNegative():
		// change rate cannot be negative
		return ErrCommissionChangeRateNegative

	case cr.MaxChangeRate.GT(cr.MaxRate):
		// change rate cannot be greater than the max rate
		return ErrCommissionChangeRateGTMaxRate
	}
	return nil
}

type Description struct {
	Moniker         string    `json:"moniker"`
	Identity        string	  `json:"identity"`
	Website         string 	  `json:"website"`
	SecurityContact string	  `json:"security_contact"`
	Details         string	  `json:"details"`
}

func NewDescription(moniker, identity, website, securityContact, details string) Description {
	return Description{
		Moniker:         moniker,
		Identity:        identity,
		Website:         website,
		SecurityContact: securityContact,
		Details:         details,
	}
}

func (m *Description) Reset() { *m = Description{}}

func (m *Description) GetMoniker() string {
	if m != nil {
		return m.Moniker
	}
	return ""
}

func (m Description) UpdateDescription(d2 Description) (Description, error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = m.Moniker
	}
	if d2.Identity == DoNotModifyDesc {
		d2.Identity = m.Identity
	}
	if d2.Website == DoNotModifyDesc {
		d2.Website = m.Website
	}
	if d2.SecurityContact == DoNotModifyDesc {
		d2.SecurityContact = m.SecurityContact
	}
	if d2.Details == DoNotModifyDesc {
		d2.Details = m.Details
	}

	return NewDescription(
		d2.Moniker,
		d2.Identity,
		d2.Website,
		d2.SecurityContact,
		d2.Details,
	).EnsureLength()
}

func (m *Description) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

func (m *Description) GetWebsite() string {
	if m != nil {
		return m.Website
	}
	return ""
}

func (m *Description) GetSecurityContact() string {
	if m != nil {
		return m.SecurityContact
	}
	return ""
}

func (m *Description) GetDetails() string {
	if m != nil {
		return m.Details
	}
	return ""
}

func (m Description) EnsureLength() (Description, error) {
	if len(m.Moniker) > MaxMonikerLength {
		return m, Wrapf(ErrInvalidRequest, "invalid moniker length; got: %d, max: %d", len(m.Moniker), MaxMonikerLength)
	}
	if len(m.Identity) > MaxIdentityLength {
		return m, Wrapf(ErrInvalidRequest, "invalid identity length; got: %d, max: %d", len(m.Identity), MaxIdentityLength)
	}
	if len(m.Website) > MaxWebsiteLength {
		return m, Wrapf(ErrInvalidRequest, "invalid website length; got: %d, max: %d", len(m.Website), MaxWebsiteLength)
	}
	if len(m.SecurityContact) > MaxSecurityContactLength {
		return m, Wrapf(ErrInvalidRequest, "invalid security contact length; got: %d, max: %d", len(m.SecurityContact), MaxSecurityContactLength)
	}
	if len(m.Details) > MaxDetailsLength {
		return m, Wrapf(ErrInvalidRequest, "invalid details length; got: %d, max: %d", len(m.Details), MaxDetailsLength)
	}
	return m, nil
}

type Commission struct {
	CommissionRates   CommissionRates   `json:"commission_rates"`
	UpdateTime        time.Time         `json:"update_time"`
}

func (m *Commission) Reset()      { *m = Commission{} }

func (m Commission) String() string {
	out, _ := json.Marshal(m)
	return string(out)
}

func (m Commission) ValidateNewRate(newRate sdk.Dec, blockTime time.Time) error {
	switch {
	case blockTime.Sub(m.UpdateTime).Hours() < 24:
		// new rate cannot be changed more than once within 24 hours
		return ErrCommissionUpdateTime
	case newRate.IsNegative():
		// new rate cannot be negative
		return ErrCommissionNegative
	case newRate.GT(m.CommissionRates.MaxRate):
		// new rate cannot be greater than the max rate
		return ErrCommissionGTMaxRate
	case newRate.Sub(m.CommissionRates.Rate).GT(m.CommissionRates.MaxChangeRate):
		// new rate % points change cannot be greater than the max change rate
		return ErrCommissionChangeRateGTMaxRate
	}
	return nil
}


func NewCommission(rate, maxRate, maxChangeRate sdk.Dec) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:       time.Unix(0, 0).UTC(),
	}
}

// NewCommissionWithTime returns an initialized validator commission with a specified
// update time which should be the current block BFT time.
func NewCommissionWithTime(rate, maxRate, maxChangeRate sdk.Dec, updatedAt time.Time) Commission {
	return Commission{
		CommissionRates: NewCommissionRates(rate, maxRate, maxChangeRate),
		UpdateTime:      updatedAt,
	}
}