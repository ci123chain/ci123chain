package gravity

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"encoding/json"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/ci123chain/ci123chain/pkg/gravity/keeper"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic object for module implementation
type AppModuleBasic struct{}

// Name implements app module basic
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}


// DefaultGenesis implements app module basic
func (AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	bz, _ := json.Marshal(types.DefaultGenesisState())
	return bz
}

//____________________________________________________________________________

// AppModule object for module implementation
type AppModule struct {
	AppModuleBasic
	keeper     keeper.Keeper
	accKeeper account.AccountKeeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k keeper.Keeper, accKeeper account.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		accKeeper:     	accKeeper,
	}
}

// Name implements app module
func (AppModule) Name() string {
	return types.ModuleName
}

// InitGenesis initializes the genesis state for this module and implements app module.
func (am AppModule) InitGenesis(ctx sdk.Context, bz json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	err := json.Unmarshal(bz, &genesisState)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", types.ModuleName, err))
	}
	keeper.InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

//// ExportGenesis exports the current genesis state to a json.RawMessage
//func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
//	gs := keeper.ExportGenesis(ctx, am.keeper)
//	return cdc.MustMarshalJSON(&gs)
//}

// EndBlock implements app module
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.keeper)
	// this begin blocker is only for testing purposes, don't import into your
	// own chain running gravity
	TestingEndBlocker(ctx, am.keeper)
	return []abci.ValidatorUpdate{}
}

func (m AppModule) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
}

func (m AppModule) Committer(ctx sdk.Context) {
}