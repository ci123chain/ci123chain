package fc
/*
import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
)

var ValidatorCurrentRewardsPrefix = []byte{0x06} // key for current validator rewards

type FcKeeper struct {
	cdc 		*codec.Codec
	storeKey 	types.StoreKey
	ak 			account.AccountKeeper
}
var (
	FcStoreKey = "fc"
)


func NewFcKeeper(cdc *codec.Codec, key sdk.StoreKey, ak account.AccountKeeper) FcKeeper {
	return FcKeeper{
		cdc:      cdc,
		storeKey: key,
		ak:       ak,
	}
}

func GetValidatorCurrentRewardsKey(v types.AccAddress) []byte{
	return append(ValidatorCurrentRewardsPrefix, v.Bytes()...)
}
var (
	collectedFeesKey = []byte("collectedFees")
)

func (fck *FcKeeper) AddCollectedFees(ctx types.Context, coins types.Coin) types.Coin {
	newCoins := fck.GetCollectedFees(ctx).Add(coins)
	fck.SetCollectedFees(ctx, newCoins)

	return newCoins
}

func (fck *FcKeeper) GetCollectedFees(ctx types.Context) (fee types.Coin) {
	store := ctx.KVStore(fck.storeKey)
	b := store.Get(collectedFeesKey)
	if b == nil {
		return types.NewCoin(sdk.NewInt(0))
	}
	fck.cdc.MustUnmarshalBinaryLengthPrefixed(b, &fee)
	return
}

func (fck *FcKeeper) SetCollectedFees(ctx types.Context, coin types.Coin) {
	bz := fck.cdc.MustMarshalBinaryLengthPrefixed(coin)
	store := ctx.KVStore(fck.storeKey)
	store.Set(collectedFeesKey, bz)
}

func (fck *FcKeeper) ClearCollectedFees(ctx types.Context) {
	fck.SetCollectedFees(ctx, types.NewCoin(sdk.NewInt(0)))
}
*/