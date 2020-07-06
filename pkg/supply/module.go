package supply

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	types2 "github.com/ci123chain/ci123chain/pkg/supply/types"
)

type AppModule struct {
	AppModuleBasic
	Keeper  Keeper
}

func (am AppModule) Committer(ctx types.Context) {}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	err := ModuleCdc.UnmarshalJSON(data, &genesisState)
	if err != nil {
		panic(err)
	}
	InitGenesis(ctx, am.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) BeginBlocker(ctx types.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(_ types.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}


type AppModuleBasic struct {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator, _ []string) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types2.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}