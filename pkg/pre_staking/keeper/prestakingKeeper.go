package keeper

import (
	"encoding/json"
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
	"math/big"
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

func (ps PreStakingKeeper) SetAccountPreStaking(ctx sdk.Context, delegator sdk.AccAddress, vaults types.VaultRecord) {
	store := ctx.KVStore(ps.storeKey)
	bz, err := json.Marshal(vaults)
	if err != nil {
		panic(err)
	}
	store.Set(types.GetPreStakingKey(delegator), bz)
}



func (ps PreStakingKeeper) GetAccountPreStaking(ctx sdk.Context, delegator sdk.AccAddress) types.VaultRecord {
	store := ctx.KVStore(ps.storeKey)
	bz := store.Get(types.GetPreStakingKey(delegator))
	if bz == nil {
		return types.NewEmptyVaultRecord()
	}

	var vats types.VaultRecord
	err := json.Unmarshal(bz, &vats)
	if err != nil {
		panic(err)
	}

	return vats
}

func (ps PreStakingKeeper) AccountPreStakingIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(ps.storeKey)
	prefix := types.PreStakingKey
	iterator := store.RemoteIterator(prefix, sdk.PrefixEndBytes(prefix))
	if iterator.Valid() {
		return iterator
	} else {
		iterator.Close()
		store := ctx.KVStore(ps.storeKey)
		return sdk.KVStorePrefixIterator(store, prefix)
	}
}

func (ps PreStakingKeeper) GetAllAccountPreStaking(ctx sdk.Context) []types.InitPrestaking {
	var res = make([]types.InitPrestaking, 0)
	iterator := ps.AccountPreStakingIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		del := sdk.ToAccAddress(key[1:])
		var vr types.VaultRecord
		_ = json.Unmarshal(iterator.Value(), &vr)
		r := types.InitPrestaking{
			Delegator: del,
			Staking:   vr,
		}
		res = append(res, r)
	}
	return res
}


func (ps PreStakingKeeper) SetAccountStakingRecord(ctx sdk.Context, val, del sdk.AccAddress, id *big.Int, et time.Time, amount sdk.Coin) error {
	store := ctx.KVStore(ps.storeKey)
	var record = types.NewStakingRecord(id, et, amount)
	var key = types.GetStakingRecordKey(del, val)

	before := ps.GetAccountStakingRecord(ctx, val, del)
	var records = make([]types.StakingRecord, 0)
	if before != nil {
		records = append(records, before...)
	}
	records = append(records, record)
	//bz := ps.cdc.MustMarshalBinaryLengthPrefixed(records)
	bz, err := json.Marshal(records)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (ps PreStakingKeeper) SetAccountStakingRecords(ctx sdk.Context, del, val sdk.AccAddress, records []types.StakingRecord) {
	store := ctx.KVStore(ps.storeKey)
	var key = types.GetStakingRecordKey(del, val)
	//bz := ps.cdc.MustMarshalBinaryLengthPrefixed(records)
	bz, err := json.Marshal(records)
	if err != nil {
		panic(err)
	}
	store.Set(key, bz)
}

func (ps PreStakingKeeper) GetAccountStakingRecord(ctx sdk.Context, val, del sdk.AccAddress) []types.StakingRecord {
	store := ctx.KVStore(ps.storeKey)
	var key = types.GetStakingRecordKey(del, val)
	bz := store.Get(key)
	var res []types.StakingRecord
	if bz == nil {
		return nil
	}
	//ps.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &res)
	err := json.Unmarshal(bz, &res)
	if err != nil {
		return nil
	}
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
	//bz := ps.cdc.MustMarshalBinaryLengthPrefixed(records)
	bz, err := json.Marshal(records)
	if err != nil {
		panic(err)
	}
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

func (ps PreStakingKeeper) GetAllStakingRecords(ctx sdk.Context) []types.InitStakingRecords{
	res := make([]types.InitStakingRecords, 0)
	iterator := ps.StakingRecordIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		val, del := getValDelFromKey(iterator.Key())
		var sr []types.StakingRecord
		_ = json.Unmarshal(iterator.Value(), &sr)
		r := types.InitStakingRecords{
			Delegator: del,
			Validator: val,
			Records:   sr,
		}
		res = append(res, r)
	}
	return res
}

func (ps PreStakingKeeper) UpdateDeadlineRecord(ctx sdk.Context) {
	iterator := ps.StakingRecordIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		k := iterator.Key()
		val, del := getValDelFromKey(k)
		v := iterator.Value()
		if v != nil {
			var records types.StakingRecords
			//ps.cdc.MustUnmarshalBinaryLengthPrefixed(v, &records)
			_ = json.Unmarshal(v, &records)
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