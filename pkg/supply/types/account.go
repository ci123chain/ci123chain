package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	types2 "github.com/tanhuiya/ci123chain/pkg/account/types"
	"github.com/tendermint/tendermint/crypto"
)

type ModuleAccount struct {
	*account.BaseAccount

	Name 	string	`json:"name" yaml:"name"`
	Permissions 	[]string `json:"permissions" yaml:"permissions"`
}

func NewModuleAddress(name string) types.AccAddress {
	return types.ToAccAddress(crypto.AddressHash([]byte(name)))
}

func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := types2.NewBaseAccountWithAddress(moduleAddress)

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: 	&baseAcc,
		Name: 			name,
		Permissions:	permissions,
	}
}


