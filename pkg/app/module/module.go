package module

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	account "github.com/ci123chain/ci123chain/pkg/account/module"
	"github.com/ci123chain/ci123chain/pkg/auth"
	dis_basic "github.com/ci123chain/ci123chain/pkg/distribution/module/basic"
	"github.com/ci123chain/ci123chain/pkg/gravity"
	ibctransfer "github.com/ci123chain/ci123chain/pkg/ibc/application/transfer"
	ibc "github.com/ci123chain/ci123chain/pkg/ibc/core"
	mint_basic "github.com/ci123chain/ci123chain/pkg/mint/module/basic"
	order_basic "github.com/ci123chain/ci123chain/pkg/order/module/basic"
	prestaking_basic "github.com/ci123chain/ci123chain/pkg/pre_staking/module/basic"
	"github.com/ci123chain/ci123chain/pkg/slashing"
	staking_basic "github.com/ci123chain/ci123chain/pkg/staking/module/basic"
	supply_basic "github.com/ci123chain/ci123chain/pkg/supply/module/basic"
	wasm_basic "github.com/ci123chain/ci123chain/pkg/vm/module/basic"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

var ModuleBasics = module.NewBasicManager(
	account.AppModuleBasic{},
	auth.AppModuleBasic{},
	supply_basic.AppModuleBasic{},
	order_basic.AppModuleBasic{},
	staking_basic.AppModuleBasic{},
	prestaking_basic.AppModuleBasic{},
	slashing.AppModuleBasic{},
	gravity.AppModuleBasic{},
	mint_basic.AppModuleBasic{},
	wasm_basic.AppModuleBasic{},
	dis_basic.AppModuleBasic{},
	ibctransfer.AppModuleBasic{},
	ibc.AppModuleBasic{},
)

func AppGetValidator(pk crypto.PubKey, name string) types.GenesisValidator {
	validator := types.GenesisValidator{
		PubKey: pk,
		Power:  1,
		Name:   name,
	}
	return validator
}
