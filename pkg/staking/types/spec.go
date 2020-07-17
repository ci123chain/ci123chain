package types

import sdk "github.com/ci123chain/ci123chain/pkg/abci/types"


//与解除绑定相关， 用于存储和获取 到达解绑时间的 Delegator和Validator对；
type DVPair struct {
	DelegatorAddress   sdk.AccAddress   `json:"delegator_address"`
	ValidatorAddress  sdk.AccAddress     `json:"validator_address"`
}

type DVPairs struct {
	Pairs []DVPair    `json:"pairs"`
}

//与二次绑定有关；
type DVVTriplet struct {
	DelegatorAddress sdk.AccAddress     `json:"delegator_address"`
	ValidatorSrcAddress sdk.AccAddress  `json:"validator_src_address"`
	ValidatorDstAddress sdk.AccAddress  `json:"validator_dst_address"`
}

type DVVTriplets struct {
	Triplets []DVVTriplet  `json:"triplets"`
}