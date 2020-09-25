package module

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/auth"
	distr "github.com/ci123chain/ci123chain/pkg/distribution"
	"github.com/ci123chain/ci123chain/pkg/mint"
	"github.com/ci123chain/ci123chain/pkg/order"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"github.com/ci123chain/ci123chain/pkg/supply"
	wasm_module "github.com/ci123chain/ci123chain/pkg/wasm/module"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

var ModuleBasics = module.NewBasicManager(
	account.AppModuleBasic{},
	auth.AppModuleBasic{},
	supply.AppModuleBasic{},
	order.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	wasm_module.AppModuleBasic{},
	distr.AppModuleBasic{},
)

func AppGetValidator(pk crypto.PubKey, name string) types.GenesisValidator {
	validator := types.GenesisValidator{
		PubKey: pk,
		Power:  1,
		Name:   name,
	}
	return validator
}
