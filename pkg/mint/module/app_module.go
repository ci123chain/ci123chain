package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/mint"
	"github.com/ci123chain/ci123chain/pkg/mint/keeper"
	"github.com/ci123chain/ci123chain/pkg/mint/module/basic"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic

	Keeper keeper.MinterKeeper
}

/*func NewAppModule(keeper keeper.MinterKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		Keeper:         keeper,
	}
}*/

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//
	var genesisState mint.GenesisState
	err := mint.ModuleCdc.UnmarshalJSON(data, &genesisState)
	if err != nil {
		panic(err)
	}
	mint.InitGenesis(ctx, am.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	mint.BeginBlocker(ctx, am.Keeper)
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	mint.EndBlocker(ctx, am.Keeper)
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}