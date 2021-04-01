package types

import (
	"github.com/tendermint/tendermint/proto/tendermint/types"
	"sort"
)

// HistoricalInfo contains the historical information that gets stored at
// each height.
type HistoricalInfo struct {
	Header    types.Header   `json:"header"`
	Valset    []Validator	 `json:"valset"`
}


// NewHistoricalInfo will create a historical information struct from header and valset
// it will first sort valset before inclusion into historical info
func NewHistoricalInfo(header types.Header, valSet Validators) HistoricalInfo {
	sort.Sort(valSet)
	return HistoricalInfo{
		Header: header,
		Valset: valSet,
	}
}