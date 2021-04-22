package keeper

import (
	"container/list"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/params"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

var ModuleCdc *codec.Codec
const aminoCacheSize = 500

type StakingKeeper struct {
	storeKey            sdk.StoreKey
	cdc                 *codec.Codec
	AccountKeeper       account.AccountKeeper
	SupplyKeeper        supply.Keeper
	hooks               types.StakingHooks
	paramstore          params.Subspace
	validatorCache      map[string]cachedValidator
	validatorCacheList  *list.List
	cdb				    dbm.DB
}

const (
	RouteKey = "staking"
)


func NewStakingKeeper(cdc *codec.Codec, key sdk.StoreKey, ak account.AccountKeeper, sk supply.Keeper, ps params.Subspace, cdb dbm.DB) StakingKeeper {
	return StakingKeeper{
		storeKey: key,
		cdc:      cdc,
		AccountKeeper:       ak,
		SupplyKeeper:sk,
		paramstore:ps.WithKeyTable(types.ParamKeyTable()),
		hooks:nil,
		validatorCache: make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
		cdb:	  cdb,
	}


}

func (k StakingKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *StakingKeeper) SetHooks(sh types.StakingHooks) *StakingKeeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

// Set the last total validator power.
func (k StakingKeeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(&sdk.IntProto{Int: power})
	store.Set(types.LastTotalPowerKey, bz)
}

// iterate through the validator set and perform the provided function
func (k StakingKeeper) IterateValidators(ctx sdk.Context, fn func(index int64, validator types.Validator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := store.RemoteIterator(types.ValidatorsKey, sdk.PrefixEndBytes(types.ValidatorsKey))
	if !iterator.Valid() {
		iterator.Close()
		iterator = sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	}

	i := int64(0)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		validator := MustUnmarshalValidator(k.cdc, iterator.Value())
		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?

		if stop {
			break
		}
		i++
	}
}

// unmarshal a redelegation from a store value
func MustUnmarshalValidator(cdc *codec.Codec, value []byte) types.Validator {
	validator, err := UnmarshalValidator(cdc, value)
	if err != nil {
		panic(err)
	}

	return validator
}

// unmarshal a redelegation from a store value
func UnmarshalValidator(cdc *codec.Codec, value []byte) (v types.Validator, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &v)
	return v, err
}