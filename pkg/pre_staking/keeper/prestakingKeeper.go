package keeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/params"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
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

	UpgradeKeeper 		upgrade.Keeper
}

func NewPreStakingKeeper(cdc *codec.Codec, key sdk.StoreKey, ak account.AccountKeeper, sk supply.Keeper, stakingKeeper keeper.StakingKeeper, upgradeKeeper upgrade.Keeper,
	ps params.Subspace, cdb dbm.DB) PreStakingKeeper{
	p := PreStakingKeeper{
		storeKey:key,
		cdc: cdc,
		AccountKeeper: ak,
		SupplyKeeper: sk,
		StakingKeeper: stakingKeeper,
		UpgradeKeeper: upgradeKeeper,
		paramstore: ps,
		cdb:  cdb,
	}
	return p
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

func (ps PreStakingKeeper) DelAccountPreStaking(ctx sdk.Context, delegator sdk.AccAddress, vaults types.VaultRecord) {
	store := ctx.KVStore(ps.storeKey)
	store.Delete(types.GetPreStakingKey(delegator))
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

func (ps PreStakingKeeper) UpdateAccountPreStaking(ctx sdk.Context, vaults map[string]types.Vault, delegator sdk.AccAddress) error {
	store := ctx.KVStore(ps.storeKey)
	bz := store.Get(types.GetPreStakingKey(delegator))
	if bz == nil {
		return errors.New("no vaults")
	}

	var vats types.VaultRecord
	err := json.Unmarshal(bz, &vats)
	if err != nil {
		return err
	}
	vats.Vaults = vaults
	bz, err = json.Marshal(vaults)
	if err != nil {
		return err
	}
	store.Set(types.GetPreStakingKey(delegator), bz)
	return nil
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
	iterator := store.RemoteIterator(prefix, prefix)
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
		del, val := getValDelFromKey(k)
		v := iterator.Value()
		if v != nil {
			var records types.StakingRecords
			//ps.cdc.MustUnmarshalBinaryLengthPrefixed(v, &records)
			_ = json.Unmarshal(v, &records)
			sort.Sort(records)
			var newrecords = make(types.StakingRecords, 0)
			old := ps.GetAccountPreStaking(ctx, del)
			for _, value := range records {
				if value.EndTime.Before(ctx.BlockTime()) {
					err := ps.RemoveDeadlineDelegationAndWithdraw(ctx, val, del, value.Amount)
					if err != nil {
						panic(err)
					}
					o := old.Vaults[value.VaultID.String()]
					n := types.NewVault(o.StartTime, o.EndTime, o.StorageTime, value.Amount)
					old.Vaults[value.VaultID.String()] = n
				}else {
					newrecords = append(newrecords, value)
				}
			}
			if len(newrecords) > 0 {
				ps.SetAccountStakingRecords(ctx, del, val, newrecords)
			}else {
				ps.SetAccountStakingRecords(ctx, del, val, nil)
			}
			if len(old.Vaults) > 0 {
				ps.SetAccountPreStaking(ctx, del, old)
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
	//_, err := ps.StakingKeeper.Undelegate(ctx, del, val, amount.Amount.ToDec())
	//if err != nil {
	//	return err
	//}

	validator, found := ps.StakingKeeper.GetValidator(ctx, val)
	if !found {
		return nil
	}

	if ps.StakingKeeper.HasMaxUnbondingDelegationEntries(ctx, del, val) {
		return nil
	}

	returnAmount, err := ps.StakingKeeper.Unbond(ctx, del, val, amount.Amount.ToDec())
	if err != nil {
		return err
	}

	if validator.IsBonded() {
		err := ps.StakingKeeper.BondedTokensToMoudleAccount(ctx, returnAmount, types.ModuleName,)
		if err != nil {
			return err
		}
	}else {
		err := ps.StakingKeeper.NotBondedTokensToModuleAccount(ctx, returnAmount, types.ModuleName,)
		if err != nil {
			return err
		}
	}
	return nil
}

//func (ps PreStakingKeeper) GetWeeLinkDao(ctx sdk.Context) string {
//	store := ctx.KVStore(ps.storeKey)
//
//	bz := store.Get(types.WeeLinkDAO)
//	if bz == nil {
//		return ""
//	}
//	return sdk.ToAccAddress(bz).String()
//}
//
//func (ps PreStakingKeeper) SetWeeLinkDao(ctx sdk.Context, addr sdk.AccAddress) {
//	store := ctx.KVStore(ps.storeKey)
//	store.Set(types.WeeLinkDAO, addr.Bytes())
//}

func (ps PreStakingKeeper) GetTokenManager(ctx sdk.Context) string {
	store := ctx.KVStore(ps.storeKey)

	bz := store.Get(types.TokenManager)
	if bz == nil {
		return ""
	}
	return sdk.ToAccAddress(bz).String()
}


func (ps PreStakingKeeper) SetTokenManager(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(ps.storeKey)
	store.Set(types.TokenManager, addr.Bytes())
}

func (ps PreStakingKeeper) GetTokenManagerOwner(ctx sdk.Context) string {
	store := ctx.KVStore(ps.storeKey)

	bz := store.Get(types.TokenManagerOwner)
	if bz == nil {
		return ""
	}
	return sdk.ToAccAddress(bz).String()
}


func (ps PreStakingKeeper) SetTokenManagerOwner(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(ps.storeKey)
	store.Set(types.TokenManagerOwner, addr.Bytes())
}
