package core

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/keeper"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


var _ module.AppModule = AppModule{}
var _ module.AppModuleBasic = AppModule{}

// AppModuleBasic is the IBC Transfer AppModuleBasic
type AppModuleBasic struct{}

// NewAppModule creates a new 20-transfer module
func NewAppModule(keeper *keeper.Keeper) AppModule {
	return AppModule{
		keeper: keeper,
	}
}

func (b AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

// Name returns the ibc module's name.
func (AppModuleBasic) Name() string {
	return host.ModuleName
}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// module.
func (m AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	bz, _ := json.Marshal(types.DefaultGenesisState())
	return bz
}

// InitGenesis performs genesis initialization for the ibc module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, bz json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	err := json.Unmarshal(bz, &gs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", host.ModuleName, err))
	}
	InitGenesis(ctx, *am.keeper, am.createLocalhost, &gs)
	return []abci.ValidatorUpdate{}
}
// AppModule implements an application module for the ibc module.
type AppModule struct {
	AppModuleBasic
	keeper *keeper.Keeper
	cdc *codec.Codec
	// create localhost by default
	createLocalhost bool
}


func (m AppModule) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
}

func (m AppModule) Committer(ctx sdk.Context) {
}

func (m AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
