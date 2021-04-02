package types


import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type FeePool struct {
	//
	CommunityPool  sdk.DecCoin   `json:"community_pool"`
}

func InitialFeePool() FeePool {
	return FeePool{
		CommunityPool: sdk.NewDecCoin(sdk.ChainCoinDenom, sdk.NewInt(0)),
	}
}

func (f FeePool) ValidateGenesis() error {
	if f.CommunityPool.IsNegative() {
		return fmt.Errorf("negative CommunityPool in distribution fee pool, is %v",
			f.CommunityPool)
	}
	return nil
}