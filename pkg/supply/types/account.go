package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	types2 "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
	"github.com/tendermint/tendermint/crypto"
)

// 0x2f8833FCe544807E6F2b030c758aFe1e0a16Eb29

var _ exported.ModuleAccountI = (*ModuleAccount)(nil)

type ModuleAccount struct {
	*account.BaseAccount

	Name 	string	`json:"name" yaml:"name"`
	Permissions 	[]string `json:"permissions" yaml:"permissions"`
}

func (macc ModuleAccount) GetName() string {
	return macc.Name
}

func (macc ModuleAccount) GetPermissions() []string {
	return macc.Permissions
}

func (macc ModuleAccount) HasPermission(perm string) bool {
	return true
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


