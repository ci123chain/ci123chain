package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	types2 "github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

// 0x2f8833FCe544807E6F2b030c758aFe1e0a16Eb29

var _ exported.ModuleAccountI = (*ModuleAccount)(nil)

type ModuleAccount struct {
	*account.BaseAccount

	Name 	string	`json:"name" yaml:"name"`
	Permissions 	[]string `json:"permissions" yaml:"permissions"`
}

func (macc ModuleAccount) SetIsModule(flag bool) error {
	return macc.BaseAccount.SetIsModule(flag)
}

func (macc ModuleAccount) GetName() string {
	return macc.Name
}

func (macc ModuleAccount) GetPermissions() []string {
	return macc.Permissions
}

func (macc ModuleAccount) HasPermission(perm string) bool {
	for _, v := range macc.Permissions {
		if v == perm {
			return true
		}
	}
	return false
}

func (macc *ModuleAccount) SetAddress(acc types.AccAddress) error {
	macc.Address = acc
	return nil
}

func (macc ModuleAccount) GetAddress() types.AccAddress {
	return macc.Address
}

func (macc ModuleAccount) GetPubKey() crypto.PubKey {
	return macc.PubKey
}

func (macc *ModuleAccount) SetPubKey(key crypto.PubKey) error {
	macc.PubKey = key
	return nil
}

func (macc *ModuleAccount) SetAccountNumber(n uint64) error {
	macc.AccountNumber = n
	return nil
}

func (macc ModuleAccount) GetAccountNumber() uint64{
	return macc.AccountNumber
}

func (macc ModuleAccount) GetSequence() uint64 {
	return macc.Sequence
}

func (macc *ModuleAccount) SetSequence(s uint64) error {
	macc.Sequence = s
	return nil
}

func (macc *ModuleAccount) SetCoins(c types.Coins) error {
	macc.Coins = c
	return nil
}

func (macc ModuleAccount) GetCoins() types.Coins {
	return macc.Coins
}

func (macc ModuleAccount) GetCodeHash() []byte {
	return macc.CodeHash
}

func (macc *ModuleAccount) SetCodeHash(hash []byte) {
	macc.CodeHash = hash
}

func (macc *ModuleAccount) SetContractType(c string) error {
	macc.ContractType = c
	return nil
}

func (macc ModuleAccount) GetContractType() string {
	return macc.ContractType
}

func (macc ModuleAccount) String() string {
	return fmt.Sprintf(`Vesting Account:
  Address:          %s
  Pubkey:           %s
  Coins:            %v
  AccountNumber:    %d
  Sequence:         %d`,
		macc.Address, macc.PubKey, macc.Coins, macc.AccountNumber, macc.Sequence,
	)
}

func (macc ModuleAccount) GetIsModule() bool {
	return macc.IsModule
}

func (macc ModuleAccount) SpendableCoins(bt time.Time) types.Coins {
	return macc.Coins
}

func NewModuleAddress(name string) types.AccAddress {
	return types.ToAccAddress(crypto.AddressHash([]byte(name)))
}

func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := types2.NewBaseAccountWithAddress(moduleAddress)
	baseAcc.SetIsModule(true)

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		BaseAccount: 	&baseAcc,
		Name: 			name,
		Permissions:	permissions,
	}
}

func NewModuleAccountFromBaseAccount(acc *types2.BaseAccount, name string, permissions ...string) *ModuleAccount {
	return &ModuleAccount{
		BaseAccount: acc,
		Name:        name,
		Permissions: permissions,
	}
}

