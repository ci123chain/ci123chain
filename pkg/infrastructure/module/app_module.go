package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
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
	//do nothing
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