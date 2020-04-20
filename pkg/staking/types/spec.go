package types

import sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

type DVPair struct {
	DelegatorAddress   sdk.AccAddress   `json:"delegator_address"`
	ValidatorAddress  sdk.AccAddress     `json:"validator_address"`
}

type DVPairs struct {
	Pairs []DVPair    `json:"pairs"`
}

type DVVTriplet struct {
	DelegatorAddress sdk.AccAddress     `json:"delegator_address"`
	ValidatorSrcAddress sdk.AccAddress  `json:"validator_src_address"`
	ValidatorDstAddress sdk.AccAddress  `json:"validator_dst_address"`
}

type DVVTriplets struct {
	Triplets []DVVTriplet  `json:"triplets"`
}