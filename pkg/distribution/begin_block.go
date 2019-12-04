package distribution


import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


func BeginBlock(ctx types.Context, req abci.RequestBeginBlock, distr DistrKeeper) {

	str := Ca(req.Header.ProposerAddress).String()
	res, _ := sdk.ConsAddressFromBech32(str)

	address := "0x" + fmt.Sprintf("%X", []byte(res))
	pAddress := types.HexToAddress(address)

	if ctx.BlockHeight() > 1 {
		fee := distr.feeCollectionKeeper.GetCollectedFees(ctx)
		//分配完奖励金之后清空奖金池
		distr.feeCollectionKeeper.ClearCollectedFees(ctx)
		//distr.DistributeRewardsToValidators(ctx, pAddress, fee)

		distr.SetProposerCurrentRewards(ctx, pAddress, fee)
	}
}
