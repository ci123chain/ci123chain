package slashing

import (
	"encoding/json"
	staking "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/slashing/keeper"
	"github.com/ci123chain/ci123chain/pkg/slashing/types"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the slashing module.
type AppModuleBasic struct {
	cdc *codec.Codec
}

var _ module.AppModuleBasic = AppModuleBasic{}

// Name returns the slashing module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// DefaultGenesis returns default genesis state as raw bytes for the slashing
// module.
func (am AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	return types.SlashingCodec.MustMarshalJSON(types.DefaultGenesisState())
}

func (AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}


//____________________________________________________________________________

// AppModule implements an application module for the slashing module.
type AppModule struct {
	AppModuleBasic

	Keeper        keeper.Keeper
	AccountKeeper types.AccountKeeper
	StakingKeeper staking.StakingKeeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc *codec.Codec, keeper keeper.Keeper, ak types.AccountKeeper, sk staking.StakingKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		Keeper:         keeper,
		AccountKeeper:  ak,
		StakingKeeper:  sk,
	}
}

// Name returns the slashing module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// InitGenesis performs genesis initialization for the slashing module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.SlashingCodec.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.Keeper, am.StakingKeeper, &genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the slashing
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc *codec.Codec) json.RawMessage {
	gs := ExportGenesis(ctx, am.Keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock returns the begin blocker for the slashing module.
func (am AppModule) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
	BeginBlocker(ctx, req, am.Keeper)
}

// EndBlock returns the end blocker for the slashing module. It returns no validator
// updates.
func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}