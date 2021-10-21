package keeper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"testing"
)

func TestCodec(t *testing.T) {
	commision := &types.ValidatorAccumulatedCommission{
		Commission: sdk.NewDecCoin("stake", sdk.NewInt(100)),
	}
	bz :=types.DistributionCdc.MustMarshalBinaryLengthPrefixed(commision)
	var aa types.ValidatorAccumulatedCommission
	types.DistributionCdc.MustUnmarshalBinaryLengthPrefixed(bz, &aa)
}