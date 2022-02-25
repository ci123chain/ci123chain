package keeper

import (
	"bytes"
	"encoding/binary"
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
	"time"
)

type PreStakingKeeper struct {
	storeKey            sdk.StoreKey
	Cdc                 *codec.Codec
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
		Cdc: cdc,
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

func (ps PreStakingKeeper) SetStakingVault(ctx sdk.Context, val, del sdk.AccAddress, endtime time.Time, duration time.Duration, coin sdk.Coin)  {
	id := ps.GetStakingRecordID(ctx)
	id++
	vat := types.NewStakingVault(id, ctx.BlockTime(), endtime, duration, coin, val, del)
	key := types.GetStakingRecordKeyByID(id)
	vatbz := ps.Cdc.MustMarshalBinaryBare(vat)
	store := ctx.KVStore(ps.storeKey)
	store.Set(key, vatbz)
	ps.SetStakingRecordID(ctx, id)
}

func (ps PreStakingKeeper) GetAllStakingVault(ctx sdk.Context) (res []types.StakingVault) {
	store := ctx.KVStore(ps.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakingRecordKey)
	for ; iterator.Valid(); iterator.Next() {
		var sv types.StakingVault
		ps.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &sv)
		res = append(res, sv)
	}
	return res
}

func (ps PreStakingKeeper) DeleteStakingVault(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(ps.storeKey)
	store.Delete(key)
}


func (ps PreStakingKeeper) GetStakingRecordID(ctx sdk.Context) uint64 {
	key := types.GetStakingRecordID()
	store := ctx.KVStore(ps.storeKey)
	idbz := store.Get(key)
	if idbz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(idbz)
}

func (ps PreStakingKeeper) SetStakingRecordID(ctx sdk.Context, n uint64) {
	key := types.GetStakingRecordID()
	store := ctx.KVStore(ps.storeKey)
	store.Set(key, sdk.Uint64ToBigEndian(n))
}

func (ps PreStakingKeeper) UpdateStakingRecordProcessed(ctx sdk.Context, key []byte) {
	sv := ps.getStakingRecord(ctx, key)
	sv.Processed = true
	ps.saveStakingRecord(ctx, key, sv)
}

func (ps PreStakingKeeper) ChangeStakingRecordToNewValidator(ctx sdk.Context, recordID uint64, srcValidator, dstValidator sdk.AccAddress) error {
	key := types.GetStakingRecordKeyByID(recordID)
	store := ctx.KVStore(ps.storeKey)
	bz := store.Get(key)
	var sv types.StakingVault
	ps.Cdc.MustUnmarshalBinaryBare(bz, &sv)

	if !bytes.Equal(sv.Validator.Bytes(), srcValidator.Bytes()) {
		return types.ErrInvalidDelegatorAddress
	}
	sv.TransLogs = append(sv.TransLogs, sv.Validator)
	sv.Validator = dstValidator

	bz = ps.Cdc.MustMarshalBinaryBare(sv)
	store.Set(key, bz)
	return nil
}



func (ps PreStakingKeeper) RemoveDeadlineDelegationAndWithdraw(ctx sdk.Context, val, del sdk.AccAddress, amount sdk.Coin) (sdk.Int, error) {
	//_, err := ps.StakingKeeper.Undelegate(ctx, del, val, amount.Amount.ToDec())
	//if err != nil {
	//	return err
	//}

	validator, found := ps.StakingKeeper.GetValidator(ctx, val)
	if !found {
		return sdk.NewInt(0), nil
	}

	if ps.StakingKeeper.HasMaxUnbondingDelegationEntries(ctx, del, val) {
		return sdk.NewInt(0), nil
	}

	returnAmount, err := ps.StakingKeeper.Unbond(ctx, del, val, amount.Amount.ToDec())
	if err != nil {
		return sdk.NewInt(0), err
	}

	if validator.IsBonded() {
		err := ps.StakingKeeper.BondedTokensToMoudleAccount(ctx, returnAmount, types.ModuleName,)
		if err != nil {
			return sdk.NewInt(0), err
		}
	}else {
		err := ps.StakingKeeper.NotBondedTokensToModuleAccount(ctx, returnAmount, types.ModuleName,)
		if err != nil {
			return sdk.NewInt(0), err
		}
	}
	return returnAmount, nil
}


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

func (ps PreStakingKeeper) getStakingRecord(ctx sdk.Context, key []byte) types.StakingVault {
	store := ctx.KVStore(ps.storeKey)
	bz := store.Get(key)
	var sv types.StakingVault
	ps.Cdc.MustUnmarshalBinaryBare(bz, &sv)
	return sv
}
func (ps PreStakingKeeper) saveStakingRecord(ctx sdk.Context, key []byte, vault types.StakingVault)  {
	store := ctx.KVStore(ps.storeKey)
	bz := ps.Cdc.MustMarshalBinaryBare(vault)
	store.Set(key, bz)
}

func (ps PreStakingKeeper)Iter(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(ps.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.StakingRecordKey)
	return iterator
}