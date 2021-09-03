package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	i "github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/keeper"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/module/basic"
	abci "github.com/tendermint/tendermint/abci/types"
)


type AppModule struct {
	basic.AppModuleBasic

	Keeper keeper.InfrastructureKeeper
}


func (am AppModule) Committer(ctx sdk.Context) {
	//do nothing
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var res i.GenesisState
	_ = json.Unmarshal(data, &res)
	i.InitGenesis(ctx, am.Keeper, res)
	return nil
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	//do nothing
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	//do nothing
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	res, _ := json.Marshal(i.ExportGenesis(ctx, am.Keeper))
	return res
}