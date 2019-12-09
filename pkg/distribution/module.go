package distribution

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	k "github.com/tanhuiya/ci123chain/pkg/distribution/keeper"
	acc_types "github.com/tanhuiya/ci123chain/pkg/distribution/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	AppModuleBasic

	DistributionKeeper  k.DistrKeeper
}

func (am AppModule) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) {
	BeginBlock(ctx, req, am.DistributionKeeper)
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage)  {
	//do something
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