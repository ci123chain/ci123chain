package auth

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	abci_types "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/auth/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	_ module.AppModule = AppModule{}
)

type AppModuleBasic struct {
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	return
}

func (am AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {}

// Name returns the auth module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

func (AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModule) InitGenesis(ctx abci_types.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.AuthKeeper, genesisState)
	return nil
}

type AppModule struct {
	AppModuleBasic

	AuthKeeper AuthKeeper
}

func (am AppModule) EndBlock(ctx abci_types.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	am.AuthKeeper.SetNumTxs(ctx)
	//panic("implement me")
	return nil
}

func (am AppModule) BeginBlocker(ctx abci_types.Context, req abci.RequestBeginBlock) {
	//do you want to do
	numTxs := am.AuthKeeper.GetNumTxs(ctx)
	gasPrice := ctx.GasPriceConfig().BaseGasPrice + float64(numTxs * numTxs)/ctx.GasPriceConfig().GasPriceTxBase

	abci_types.SetGasPrice(gasPrice)
}

func (am AppModule) Committer(ctx abci_types.Context) {

}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}