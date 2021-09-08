package keeper

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/params"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type PreStakingKeeper struct {
	storeKey            sdk.StoreKey
	cdc                 *codec.Codec
	AccountKeeper       account.AccountKeeper
	SupplyKeeper        supply.Keeper
	StakingKeeper       keeper.StakingKeeper
	paramstore          params.Subspace
	cdb				    dbm.DB
}

func NewPreStakingKeeper(cdc *codec.Codec, key sdk.StoreKey, ak account.AccountKeeper, sk supply.Keeper, stakingKeeper keeper.StakingKeeper,
	ps params.Subspace, cdb dbm.DB) PreStakingKeeper{
		return PreStakingKeeper{
			storeKey:key,
			cdc: cdc,
			AccountKeeper: ak,
			SupplyKeeper: sk,
			StakingKeeper:stakingKeeper,
			paramstore: ps,
			cdb:  cdb,
		}
}


func (ps PreStakingKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (ps PreStakingKeeper) SetAccountPreStaking(ctx sdk.Context, delegator sdk.AccAddress, amount sdk.Int) {
	store := ctx.KVStore(ps.storeKey)
	bz := ps.cdc.MustMarshalBinaryLengthPrefixed(sdk.IntProto{Int:amount})
	store.Set(types.GetPreStakingKey(delegator), bz)
}



func (ps PreStakingKeeper) GetAccountPreStaking(ctx sdk.Context, delegator sdk.AccAddress) sdk.Int {
	store := ctx.KVStore(ps.storeKey)
	bz := store.Get(types.GetPreStakingKey(delegator))
	if bz == nil {
		return sdk.ZeroInt()
	}

	ip := sdk.IntProto{}
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &ip)

	return ip.Int
}