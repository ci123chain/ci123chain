package staking

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	ak "github.com/tanhuiya/ci123chain/pkg/account/keeper"
	k "github.com/tanhuiya/ci123chain/pkg/staking/keeper"
	"github.com/tanhuiya/ci123chain/pkg/staking/types"
	sk "github.com/tanhuiya/ci123chain/pkg/supply/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	AppModuleBasic
	StakingKeeper k.StakingKeeper
	AccountKeeper ak.AccountKeeper
	SupplyKeeper  sk.Keeper
}

func (a AppModule) Name() string {
	//panic("implement me")
	return ModuleName
}

func (a AppModule) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (a AppModule) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.StakingKeeper, am.AccountKeeper, am.SupplyKeeper, genesisState)
}

func (a AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlock(ctx, a.StakingKeeper)
}

func (a AppModule) Committer(ctx sdk.Context) {
	//
}

func (a AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlock(ctx, a.StakingKeeper)
}

type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	p := types.DefaultGenesisState()
	b,err := json.Marshal(p)
	if err != nil{
		panic(err)
	}
	return b
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}