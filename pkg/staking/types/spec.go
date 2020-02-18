package types

import sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"

type DVPair struct {
	DelegatorAddress   sdk.AccAddress
	ValidatorAddress  sdk.AccAddress
}

type DVPairs struct {
	Pairs []DVPair
}

type DVVTriplet struct {
	DelegatorAddress sdk.AccAddress
	ValidatorSrcAddress sdk.AccAddress
	ValidatorDstAddress sdk.AccAddress
}

type DVVTriplets struct {
	Triplets []DVVTriplet
}