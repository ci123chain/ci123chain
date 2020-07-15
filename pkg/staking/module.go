package staking

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	ak "github.com/ci123chain/ci123chain/pkg/account/keeper"
	k "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	sk "github.com/ci123chain/ci123chain/pkg/supply/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModule struct {
	AppModuleBasic
	StakingKeeper k.StakingKeeper
	AccountKeeper ak.AccountKeeper
	SupplyKeeper  sk.Keeper
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//
	var genesisState types.GenesisState
	err := types.StakingCodec.UnmarshalJSON(data, &genesisState)
	if err != nil {
		panic(err)
	}
	return InitGenesis(ctx, am.StakingKeeper, am.AccountKeeper, am.SupplyKeeper, genesisState)
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlock(ctx, am.StakingKeeper)
}

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return EndBlock(ctx, am.StakingKeeper)
}

type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	p := types.DefaultGenesisState(validators)
	b, err := types.StakingCodec.MarshalJSONIndent(p, "", "")
	if err != nil{
		panic(err)
	}
	return b
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}