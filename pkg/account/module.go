package account

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/keeper"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
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


func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator, accAddress []string) json.RawMessage {
	var genesisValidatorAccounts GenesisState
	for _, addr := range accAddress {
		account := acc_types.NewGenesisAccountRaw(types.HexToAddress(addr), types.NewCoin(types.NewInt(1000000)))
		genesisValidatorAccounts = append(genesisValidatorAccounts, account)
	}
	//return ModuleCdc.MustMarshalJSON(GenesisState{})
	return ModuleCdc.MustMarshalJSON(genesisValidatorAccounts)
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}

func (am AppModule) Committer(ctx types.Context) {

}