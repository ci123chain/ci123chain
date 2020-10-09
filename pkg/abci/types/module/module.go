package module

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/order"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleGenesis interface {
	AppModuleBasic

	// 根据 genesis 配置初始化
	InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate
	BeginBlocker(ctx types.Context, req abci.RequestBeginBlock)
	Committer(ctx types.Context)
	EndBlock(ctx types.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate
}

type AppModuleBasic interface {
	Name() string
	RegisterCodec(codec *codec.Codec)

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

func (bm BasicManager)RegisterCodec(cdc *codec.Codec)  {
	for _, b := range bm {
		b.RegisterCodec(cdc)
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

func (am AppManager) Committer(ctx types.Context) abci.ResponseCommit {

	for _, m := range am.Modules {
		m.Committer(ctx)
	}
	return abci.ResponseCommit{}
}

func (am AppManager) EndBlocker(ctx types.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {

	validatorUpdates := []abci.ValidatorUpdate{}
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

			validatorUpdates = moduleValUpdates
		}
	}
	return abci.ResponseEndBlock{ValidatorUpdates:validatorUpdates}
}