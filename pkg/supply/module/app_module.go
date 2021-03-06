package module

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/supply"
	"github.com/ci123chain/ci123chain/pkg/supply/module/basic"
	abci "github.com/tendermint/tendermint/abci/types"
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
		panic(err)
	}
	supply.InitGenesis(ctx, am.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) BeginBlocker(ctx types.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(_ types.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
