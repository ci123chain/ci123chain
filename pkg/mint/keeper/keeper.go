package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mint/types"
	"github.com/ci123chain/ci123chain/pkg/params"
	sk "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	supply "github.com/ci123chain/ci123chain/pkg/supply/keeper"
)

type MinterKeeper struct {
	cdc     		   *codec.Codec
	storeKey  		   sdk.StoreKey
	paramSpace  	   params.Subspace
	sk      		   sk.StakingKeeper
	supplyKeeper 	   supply.Keeper
	feeCollectorName   string
}


func NewMinterKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	sk sk.StakingKeeper, supplyKeeper supply.Keeper, feeCollectorName string,
	) MinterKeeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr.Empty() {
		panic("the mint module account has not been set")
	}

	return MinterKeeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		sk:               sk,
		supplyKeeper:     supplyKeeper,
		feeCollectorName: feeCollectorName,
	}
}

// get the minter
func (k MinterKeeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &minter)
	return
}

// set the minter
func (k MinterKeeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(minter)
	store.Set(types.MinterKey, b)
}

// GetParams returns the total set of minting parameters.
func (k MinterKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k MinterKeeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k MinterKeeper) StakingTokenSupply(ctx sdk.Context) sdk.Int {
	return k.sk.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k MinterKeeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return k.sk.BondedRatio(ctx)
}

func (k MinterKeeper) AllBonded(ctx sdk.Context) sdk.Coin {
	return k.sk.GetBondedPool(ctx).GetCoin()
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k MinterKeeper) MintCoins(ctx sdk.Context, newCoins sdk.Coin) error {
	/*if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}*/

	return k.supplyKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k MinterKeeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coin) error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

func (k MinterKeeper) SetLatestMintedCoin(ctx sdk.Context, fees sdk.Coin) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(fees)
	store.Set(types.LatestMintedKey, b)
}

func (k MinterKeeper) GetLatestMintedCoin(ctx sdk.Context) sdk.Coin {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LatestMintedKey)
	var fees sdk.Coin
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &fees)
	return fees
}