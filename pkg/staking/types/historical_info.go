package types

import (
	"github.com/tendermint/tendermint/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"sort"
)

type HistoricalInfo struct {
	Header    types.Header   `json:"header"`
	Valset    []Validator	 `json:"valset"`
}


// NewHistoricalInfo will create a historical information struct from header and valset
// it will first sort valset before inclusion into historical info
func NewHistoricalInfo(header abci.Header, valSet Validators) HistoricalInfo {
	sort.Sort(valSet)
	return HistoricalInfo{
		Header: header,
		Valset: valSet,
	}
}