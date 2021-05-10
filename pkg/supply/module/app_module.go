package module

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/supply/module/basic"
	abci "github.com/tendermint/tendermint/abci/types"
	"os"
)

type AppModule struct {
	basic.AppModuleBasic
	Keeper supply.Keeper
}

func (am AppModule) Committer(ctx types.Context) {}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState supply.GenesisState
	err := supply.ModuleCdc.UnmarshalJSON(data, &genesisState)
	if err != nil {
		ctx.Logger().Error("init genesis failed in supply module", "err", err.Error())
		os.Exit(1)
	}
	supply.InitGenesis(ctx, am.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) BeginBlocker(ctx types.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(_ types.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}