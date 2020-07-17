package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"strings"
	"time"
)

//CreationHeight:  创建的高度
//CompletionTime:  完成时间
//InitialBalance:   二次抵押的Token数量
//SharesDst:   二次抵押的shares数量；由Token数量计算得来；
type RedelegationEntry struct {
	CreationHeight   int64       `json:"creation_height"`
	CompletionTime   time.Time	 `json:"completion_time"`
	InitialBalance   sdk.Int	 `json:"initial_balance"`
	SharesDst        sdk.Dec	 `json:"shares_dst"`
}

//表示一个 特定的委托者 从SrcVal 更换绑定到 DstVal 的时间列表；
type Redelegation struct {
	DelegatorAddress sdk.AccAddress		`json:"delegator_address"`
	ValidatorSrcAddress sdk.AccAddress	`json:"validator_src_address"`
	ValidatorDstAddress sdk.AccAddress	`json:"validator_dst_address"`
	Entries           []RedelegationEntry `json:"entries"`
}

func NewRedelegation(delegatorAddr sdk.AccAddress, validatorSrcAddr, validatorDstAddr sdk.AccAddress,
	creationHeight int64, minTime time.Time, balance sdk.Int, sharesDst sdk.Dec) Redelegation {
	return Redelegation{
		DelegatorAddress:    delegatorAddr,
		ValidatorSrcAddress: validatorSrcAddr,
		ValidatorDstAddress: validatorDstAddr,
		Entries: []RedelegationEntry{
			NewRedelegationEntry(creationHeight, minTime, balance, sharesDst),
		},
	}
}

// AddEntry - append entry to the unbonding delegation
func (red *Redelegation) AddEntry(creationHeight int64, minTime time.Time, balance sdk.Int, sharesDst sdk.Dec) {
	entry := NewRedelegationEntry(creationHeight, minTime, balance, sharesDst)
	red.Entries = append(red.Entries, entry)
}

func NewRedelegationEntry(creationHeight int64, completionTime time.Time, balance sdk.Int, sharesDst sdk.Dec) RedelegationEntry {
	return RedelegationEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		SharesDst:      sharesDst,
	}
}

// String returns a human readable string representation of a Redelegation.
func (red Redelegation) String() string {
	out := fmt.Sprintf(`Redelegations between:
  Delegator:                 %s
  Source Validator:          %s
  Destination Validator:     %s
  Entries:
`,
		red.DelegatorAddress, red.ValidatorSrcAddress, red.ValidatorDstAddress,
	)

	for i, entry := range red.Entries {
		out += fmt.Sprintf(`    Redelegation Entry #%d:
      Creation height:           %v
      Min time to unbond (unix): %v
      Dest Shares:               %s
`,
			i, entry.CreationHeight, entry.CompletionTime, entry.SharesDst,
		)
	}

	return strings.TrimRight(out, "\n")
}


// Redelegations are a collection of Redelegation
type Redelegations []Redelegation

func (d Redelegations) String() (out string) {
	for _, red := range d {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}

// RedelegationEntryResponse is equivalent to a RedelegationEntry except that it
// contains a balance in addition to shares which is more suitable for client
// responses.
type RedelegationEntryResponse struct {
	RedelegationEntry
	Balance sdk.Int `json:"balance"`
}

// NewRedelegationEntryResponse creates a new RedelegationEntryResponse instance.
func NewRedelegationEntryResponse(
	creationHeight int64, completionTime time.Time, sharesDst sdk.Dec, initialBalance, balance sdk.Int) RedelegationEntryResponse {
	return RedelegationEntryResponse{
		RedelegationEntry: NewRedelegationEntry(creationHeight, completionTime, initialBalance, sharesDst),
		Balance:           balance,
	}
}

// RedelegationResponse is equivalent to a Redelegation except that its entries
// contain a balance in addition to shares which is more suitable for client
// responses.
type RedelegationResponse struct {
	Redelegation
	Entries []RedelegationEntryResponse `json:"entries" yaml:"entries"`
}

// NewRedelegationResponse crates a new RedelegationEntryResponse instance.
func NewRedelegationResponse(
	delegatorAddr sdk.AccAddress, validatorSrc, validatorDst sdk.AccAddress, entries []RedelegationEntryResponse,
) RedelegationResponse {
	return RedelegationResponse{
		Redelegation: Redelegation{
			DelegatorAddress:    delegatorAddr,
			ValidatorSrcAddress: validatorSrc,
			ValidatorDstAddress: validatorDst,
		},
		Entries: entries,
	}
}

// RedelegationResponses are a collection of RedelegationResp
type RedelegationResponses []RedelegationResponse

func (r RedelegationResponses) String() (out string) {
	for _, red := range r {
		out += red.String() + "\n"
	}
	return strings.TrimSpace(out)
}