package keeper

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/registry/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
	"github.com/tendermint/tendermint/libs/log"
)

type Keeper struct {
	cdc 		*codec.Codec
	storeKey 	sdk.StoreKey
	SupplyKeeper        supply.Keeper
	UpgradeKeeper 		upgrade.Keeper
}


func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk supply.Keeper, upgradeKeeper upgrade.Keeper) Keeper{
	p := Keeper{
		storeKey:key,
		cdc: cdc,
		SupplyKeeper: sk,
		UpgradeKeeper: upgradeKeeper,
	}
	p.upgrade()
	return p
}



func (k Keeper)upgrade()  {
	k.upgradeForOnlineRegister()
}

func (k Keeper)upgradeForOnlineRegister()  {
	k.UpgradeKeeper.SetUpgradeHandler(types.OnlineRegisterVersion, func(ctx sdk.Context, info []byte) {
		k.SetupRegistry(ctx)
	})
}

func (k Keeper)SetupRegistry(ctx sdk.Context)  {
	// default registry address: 0xB2076659C4ba32B72DdA4f8640BA837Ef044eC45
	registryAddr, err := k.SupplyKeeper.DeployRegistryContract(ctx, types.ModuleName, nil)
	k.Logger(ctx).Info("OnlineRegisty", "address", registryAddr.String())
	if err != nil {
		st := ctx.KVStore(k.storeKey)
		st.Set(types.RegistryKey(), registryAddr.Bytes())
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}