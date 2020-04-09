package account

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account/keeper"
	acc_types "github.com/tanhuiya/ci123chain/pkg/account/types"
	abci "github.com/tendermint/tendermint/abci/types"
)


type AppModule struct {
	AppModuleBasic

	AccountKeeper	keeper.AccountKeeper
}

func (am AppModule) EndBlock(ctx types.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	//panic("implement me")
	return []abci.ValidatorUpdate{}
}

func (am AppModule) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) {
	//do you want to do
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate  {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, ModuleCdc, am.AccountKeeper, genesisState)
	return nil
}


type AppModuleBasic struct {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	acc_types.RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(GenesisState{})
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}

func (am AppModule) Committer(ctx types.Context) {

}