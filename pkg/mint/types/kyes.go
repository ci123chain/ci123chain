package types

// the one key to use for the keeper store
var MinterKey = []byte{0x00}

const (
	ModuleName = "mint"
	DefaultParamspace = ModuleName

	QuerierRoute = ModuleName

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
)