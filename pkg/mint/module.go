package mint

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/mint/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator, _ []string) json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}


type AppModule struct {
	AppModuleBasic

	Keeper keeper.MinterKeeper
}

/*func NewAppModule(keeper keeper.MinterKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		Keeper:         keeper,
	}
}*/

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//
	var genesisState GenesisState
	err := ModuleCdc.UnmarshalJSON(data, &genesisState)
	if err != nil {
		panic(err)
	}
	InitGenesis(ctx, am.Keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	BeginBlocker(ctx, am.Keeper)
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	//
	return []abci.ValidatorUpdate{}
}