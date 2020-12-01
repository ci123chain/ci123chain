package mint

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mint/keeper"
	"github.com/ci123chain/ci123chain/pkg/mint/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.MinterKeeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	/*fmt.Printf("bonded ratio = %s\n", bondedRatio.String())
	fmt.Printf("minter.Inflation = %s\n", minter.Inflation.String())*/
	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
	//fmt.Printf("minter.AnnualProvisions = %s\n", minter.AnnualProvisions.String())
	k.SetMinter(ctx, minter)

	// mint coins, update supply
	mintedCoin := minter.BlockProvision(params)
	//mintedCoins := sdk.NewCoins(mintedCoin)
	//fmt.Printf("mintedCoin = %v\n", mintedCoin.Amount.Uint64())

	err := k.MintCoins(ctx, mintedCoin)
	if err != nil {
		panic(err)
	}

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoin)
	if err != nil {
		panic(err)
	}

	//ctx.EventManager().EmitEvent(
	//	sdk.NewEvent(
	//		types.EventTypeMint,
	//		sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
	//		sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
	//		sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
	//		sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
	//	),
	//)
}

func EndBlocker(ctx sdk.Context, k keeper.MinterKeeper) []abci.Event {
	events := make([]abci.Event, 0)
	totalStakingSupply := k.StakingTokenSupply(ctx)
	event := abci.Event(sdk.NewEvent(
		types.EventTypeMint,
		sdk.NewAttribute(sdk.AttributeKeyTotalSupply, totalStakingSupply.String()),
	))
	events = append(events, event)
	return events
}