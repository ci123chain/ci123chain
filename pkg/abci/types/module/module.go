package module

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	client "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/order"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleGenesis interface {
	AppModuleBasic

	// 根据 genesis 配置初始化
	InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate
	BeginBlocker(ctx types.Context, req abci.RequestBeginBlock)
	EndBlock(ctx types.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate
	ExportGenesis(types.Context) json.RawMessage
}

type AppModuleBasic interface {
	Name() string
	RegisterCodec(codec *codec.Codec)
	RegisterInterfaces(codectypes.InterfaceRegistry)

	RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux)

	// 默认的 genesis 配置
	DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage
}


// BasicManager
type BasicManager map[string]AppModuleBasic

func NewBasicManager(modules ...AppModuleBasic) BasicManager {
	moduleMap := make(map[string]AppModuleBasic)
	for _, module := range modules {
		moduleMap[module.Name()] = module
	}
	return moduleMap
}

func (bm BasicManager) RegisterCodec(cdc *codec.Codec)  {
	for _, b := range bm {
		b.RegisterCodec(cdc)
	}
}

// RegisterInterfaces registers all module interface types
func (bm BasicManager) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	for _, m := range bm {
		m.RegisterInterfaces(registry)
	}
}

// RegisterGRPCGatewayRoutes registers all module rest routes
func (bm BasicManager) RegisterGRPCGatewayRoutes(clientCtx client.Context, rtr *runtime.ServeMux) {
	for _, b := range bm {
		b.RegisterGRPCGatewayRoutes(clientCtx, rtr)
	}
}

func (bm BasicManager) DefaultGenesis(validators []tmtypes.GenesisValidator) map[string]json.RawMessage {
	genesis := make(map[string]json.RawMessage)
	for _, b := range bm {
		if b.DefaultGenesis(validators) != nil {
			genesis[b.Name()] = b.DefaultGenesis(validators)
		}
	}
	return genesis
}



type AppModule interface {
	AppModuleGenesis

	// RegisterServices allows a module to register services
	RegisterServices(Configurator)
}

type AppManager struct {
	Modules 	map[string]AppModule
	Orders     []string
}

func NewManager(orders []string,modules ...AppModule) *AppManager {
	moduleMap := make(map[string]AppModule)
	//var orders []string
	for _, module := range modules {
		moduleMap[module.Name()] = module
		//orders = append(orders, module.Name())
	}
	return &AppManager{
		Modules: moduleMap,
		Orders:  orders,
	}
}

func (am AppManager) InitGenesis(ctx types.Context, data map[string]json.RawMessage) abci.ResponseInitChain {
	var validatorUpdates []abci.ValidatorUpdate

	for _, name := range am.Orders {
		m := am.Modules[name]
		moduleValUpdates := m.InitGenesis(ctx, data[m.Name()])
		// use these validator updates if provided, the module manager assumes
		// only one module will update the validator set
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator InitGenesis updates already set by a previous module")
			}
			validatorUpdates = moduleValUpdates
		}
	}
	return abci.ResponseInitChain{
		Validators: validatorUpdates,
	}
}

func (am AppManager) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	o, ok := am.Modules[order.ModuleName]
	if ok && o != nil{
		o.BeginBlocker(ctx, req)
	}

	for _, name := range am.Orders {
		if name == order.ModuleName {
			continue
		}
		m := am.Modules[name]
		m.BeginBlocker(ctx, req)
	}
	return abci.ResponseBeginBlock{}
}

func (am AppManager) EndBlocker(ctx types.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {

	ctx = ctx.WithEventManager(types.NewEventManager())
	validatorUpdates := make([]abci.ValidatorUpdate, 0)
	//all_events := make([]abci.Event, 0)
	for _, name := range am.Orders {
		m := am.Modules[name]
		if m == am.Modules[order.ModuleName] {
			continue
		}
		moduleValUpdates := m.EndBlock(ctx, req)
		if len(moduleValUpdates) > 0 {
			if len(validatorUpdates) > 0 {
				panic("validator EndBlock updates already set by a previous module")
			}
			validatorUpdates = append(validatorUpdates, moduleValUpdates...)
		}
	}
	return abci.ResponseEndBlock{ValidatorUpdates:validatorUpdates, Events: ctx.EventManager().ABCIEvents()}
}

// RegisterServices registers all module services
func (am *AppManager) RegisterServices(cfg Configurator) {
	for _, module := range am.Modules {
		module.RegisterServices(cfg)
	}
}

func (am *AppManager) ExportGenesis(ctx types.Context) map[string]json.RawMessage {
	genesisData := make(map[string]json.RawMessage)
	for _, moduleName := range am.Orders {
		genesisData[moduleName] = am.Modules[moduleName].ExportGenesis(ctx)
	}

	return genesisData
}
