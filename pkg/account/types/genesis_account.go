package types

import (
	"bytes"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
)

// GenesisAccount is a struct for account initialization used exclusively during genesis
//type GenesisAccount struct {
//	*BaseAccount
//
//	Name 	string	`json:"name" yaml:"name"`
//	Permissions 	[]string `json:"permissions" yaml:"permissions"`
//}
//
//
//func (macc GenesisAccount) SetIsModule(flag bool) error {
//	return macc.BaseAccount.SetIsModule(flag)
//}
//
//func (macc GenesisAccount) GetName() string {
//	return macc.Name
//}
//
//func (macc GenesisAccount) GetPermissions() []string {
//	return macc.Permissions
//}
//
//func (macc GenesisAccount) HasPermission(perm string) bool {
//	for _, v := range macc.Permissions {
//		if v == perm {
//			return true
//		}
//	}
//	return false
//}
//
//func (macc *GenesisAccount) SetAddress(acc types.AccAddress) error {
//	macc.Address = acc
//	return nil
//}
//
//func (macc GenesisAccount) GetAddress() types.AccAddress {
//	return macc.Address
//}
//
//func (macc GenesisAccount) GetPubKey() crypto.PubKey {
//	return macc.PubKey
//}
//
//func (macc *GenesisAccount) SetPubKey(key crypto.PubKey) error {
//	macc.PubKey = key
//	return nil
//}
//
//func (macc *GenesisAccount) SetAccountNumber(n uint64) error {
//	macc.AccountNumber = n
//	return nil
//}
//
//func (macc GenesisAccount) GetAccountNumber() uint64{
//	return macc.AccountNumber
//}
//
//func (macc GenesisAccount) GetSequence() uint64 {
//	return macc.Sequence
//}
//
//func (macc *GenesisAccount) SetSequence(s uint64) error {
//	macc.Sequence = s
//	return nil
//}
//
//func (macc *GenesisAccount) SetCoins(c types.Coins) error {
//	macc.Coins = c
//	return nil
//}
//
//func (macc GenesisAccount) GetCoins() types.Coins {
//	return macc.Coins
//}
//
//func (macc GenesisAccount) GetCodeHash() []byte {
//	return macc.CodeHash
//}
//
//func (macc *GenesisAccount) SetCodeHash(hash []byte) {
//	macc.CodeHash = hash
//}
//
//func (macc *GenesisAccount) SetContractType(c string) error {
//	macc.ContractType = c
//	return nil
//}
//
//func (macc GenesisAccount) GetContractType() string {
//	return macc.ContractType
//}
//
//func (macc GenesisAccount) String() string {
//	return fmt.Sprintf(`Vesting Account:
//  Address:          %s
//  Pubkey:           %s
//  Coins:            %v
//  AccountNumber:    %d
//  Sequence:         %d`,
//		macc.Address, macc.PubKey, macc.Coins, macc.AccountNumber, macc.Sequence,
//	)
//}
//
//func (macc GenesisAccount) GetIsModule() bool {
//	return macc.IsModule
//}
//
//func (macc GenesisAccount) SpendableCoins(bt time.Time) types.Coins {
//	return macc.Coins
//}

// GenesisAccounts defines a set of genesis account
type GenesisAccounts []exported.Account

//func NewGenesisAccountRaw(acc *BaseAccount, name string, permissions ...string) exported.Account {
//	return GenesisAccount{
//		BaseAccount: acc, Name: name, Permissions: permissions,
//	}
//}

//func (ga GenesisAccount) Validate() error {
//	return nil
//}
//
//func (ga GenesisAccount) ToAccount() exported.Account {
//	bacc := ga.BaseAccount
//	return bacc
//}

func (gaccs GenesisAccounts) Contains(acc types.AccAddress) bool {
	for _, gacc := range gaccs {
		if bytes.Equal(gacc.GetAddress().Bytes(), acc.Bytes()) {
			return true
		}
	}
	return false
}