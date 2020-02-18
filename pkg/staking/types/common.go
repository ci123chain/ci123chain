package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"time"
)

const (
	// TODO: Why can't we just have one string description which can be JSON by convention
	MaxMonikerLength         = 70
	MaxIdentityLength        = 3000
	MaxWebsiteLength         = 140
	MaxSecurityContactLength = 140
	MaxDetailsLength         = 280
)

type CommissionRates struct {
	Rate           sdk.Dec
	MaxRate        sdk.Dec
	MaxChangeRate  sdk.Dec
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
	Moniker         string
	Identity        string
	Website         string
	SecurityContact string
	Details         string
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
	CommissionRates   CommissionRates
	UpdateTime        time.Time
}

func (m *Commission) Reset()      { *m = Commission{} }

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