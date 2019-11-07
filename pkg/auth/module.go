package auth

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/types/module"
	"github.com/tanhuiya/ci123chain/pkg/auth/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	abci_types "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

var (
	_ module.AppModule = AppModule{}
)

type AppModuleBasic struct {
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

// Name returns the auth module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModule) InitGenesis(ctx abci_types.Context, data json.RawMessage) {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.AuthKeeper, genesisState)
}

type AppModule struct {
	AppModuleBasic

	AuthKeeper AuthKeeper
}
