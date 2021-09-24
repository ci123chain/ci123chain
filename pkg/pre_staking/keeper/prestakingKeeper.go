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
	"sort"
	"time"
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


func (ps PreStakingKeeper) SetAccountStakingRecord(ctx sdk.Context, val, del sdk.AccAddress, st time.Time, amount sdk.Coin) error {
	store := ctx.KVStore(ps.storeKey)
	var update = ctx.BlockTime()
	t, err := time.ParseDuration(st.String())
	if err != nil {
		return err
	}
	var end = update.Add(t)
	var record = types.NewStakingRecord(st, update, end, amount)
	var key = types.GetStakingRecordKey(del, val)

	before := ps.GetAccountStakingRecord(ctx, val, del)
	var records = make([]types.StakingRecord, 0)
	if before != nil {
		records = append(records, before...)
	}
	records = append(records, record)
	bz := ps.cdc.MustMarshalBinaryLengthPrefixed(records)
	store.Set(key, bz)
	return nil
}

func (ps PreStakingKeeper) GetAccountStakingRecord(ctx sdk.Context, val, del sdk.AccAddress) []types.StakingRecord {
	store := ctx.KVStore(ps.storeKey)
	var key = types.GetStakingRecordKey(del, val)
	bz := store.Get(key)
	var res []types.StakingRecord
	if bz == nil {
		return nil
	}
	ps.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	return res
}

func (ps PreStakingKeeper) ClearStakingRecord(ctx sdk.Context, val, del sdk.AccAddress) {
	store := ctx.KVStore(ps.storeKey)
	var key = types.GetStakingRecordKey(del, val)
	store.Set(key, nil)
}

func (ps PreStakingKeeper) UpdateStakingRecord(ctx sdk.Context, val, del sdk.AccAddress, updates types.StakingRecords) {
	store := ctx.KVStore(ps.storeKey)
	var key = types.GetStakingRecordKey(del, val)
	before := ps.GetAccountStakingRecord(ctx, val, del)
	var records = make([]types.StakingRecord, 0)
	if before != nil {
		records = append(records, before...)
	}
	records = append(records, updates...)
	bz := ps.cdc.MustMarshalBinaryLengthPrefixed(records)
	store.Set(key, bz)
}

func (ps PreStakingKeeper) StakingRecordIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(ps.storeKey)
	prefix := types.StakingRecordKey
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if iterator.Valid() {
		return iterator
	} else {
		iterator.Close()
		store := ctx.KVStore(ps.storeKey)
		return sdk.KVStorePrefixIterator(store, prefix)
	}
}

func (ps PreStakingKeeper) UpdateDeadlineRecord(ctx sdk.Context) {
	iterator := ps.StakingRecordIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		k := iterator.Key()
		val, del := getValDelFromKey(k)
		v := iterator.Value()
		if v != nil {
			var records types.StakingRecords
			ps.cdc.MustUnmarshalBinaryLengthPrefixed(v, &records)
			sort.Sort(records)
			for _, v := range records {
				if v.EndTime.Before(ctx.BlockTime()) {
					err := ps.RemoveDeadlineDelegationAndWithdraw(ctx, val, del, v.Amount)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

func getValDelFromKey(key []byte) (sdk.AccAddress, sdk.AccAddress) {
	val := sdk.ToAccAddress(key[1:21])
	del := sdk.ToAccAddress(key[21:])
	return val, del
}

func (ps PreStakingKeeper) RemoveDeadlineDelegationAndWithdraw(ctx sdk.Context, val, del sdk.AccAddress, amount sdk.Coin) error {
	_, err := ps.StakingKeeper.Undelegate(ctx, del, val, amount.Amount.ToDec())
	if err != nil {
		return err
	}
	return nil
}