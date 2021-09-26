package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/module/basic"
	"github.com/ci123chain/ci123chain/pkg/pre_staking/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/ci123chain/ci123chain/pkg/pre_staking"
)

type AppModule struct {
	basic.AppModuleBasic
	Keeper keeper.PreStakingKeeper
}

func (am AppModule) RegisterGRPCGatewayRoutes(context.Context, *runtime.ServeMux) {
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	err := json.Unmarshal(data, &genesisState)
	if err != nil {
		panic(err)
	}
	pre_staking.InitGenesis(ctx, am.Keeper, genesisState)
	return nil
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	//
}

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	pre_staking.EndBlock(ctx, am.Keeper)
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}


func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}