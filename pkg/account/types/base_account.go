package types

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

var _ exported.Account = (*BaseAccount)(nil)

const (
	EvmContractType string = "evm"
	WasmContractType string = "wasm"
)

func ProtoBaseAccount() exported.Account  {
	return &BaseAccount{}
}

type BaseAccount struct {
	Address       	types.AccAddress `json:"address" yaml:"address"`
	Coin        	types.Coin       `json:"coin" yaml:"coin"`
	Sequence      	uint64         `json:"sequence_number" yaml:"sequence_number"`
	AccountNumber 	uint64         `json:"account_number" yaml:"account_number"`
	PubKey 			crypto.PubKey  `json:"pub_key" yaml:"pub_key"`
	ContractList    []string        `json:"contract_list" yaml:"contract_list"`
	ContractType    string         `json:"contract_type" yaml:"contract_type"`
	CodeHash		[]byte 			`json:"code_hash" yaml:"code_hash"`
}

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address types.AccAddress, coin types.Coin,
	pubKey crypto.PubKey, accountNumber uint64, sequence uint64) *BaseAccount {

	return &BaseAccount{
		Address:       address,
		Coin:          coin,
		PubKey:        pubKey,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
}

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr types.AccAddress) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

//func (acc *BaseAccount) AddContract(contractAddress types.AccAddress) {
//	contractAddrStr := contractAddress.String()
//	for _,v := range acc.ContractList {
//		if v == contractAddrStr {
//			return
//		}
//	}
//	acc.ContractList = append(acc.ContractList, contractAddrStr)
//}
//
//func (acc BaseAccount) GetContractList() []string {
//	return acc.ContractList
//}

// GetAddress - Implements sdk.Account.
func (acc *BaseAccount) GetAddress() types.AccAddress {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr types.AccAddress) error {
	if !acc.Address.Empty(){
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}


// GetPubKey - Implements sdk.Account.
func (acc *BaseAccount) GetPubKey() crypto.PubKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PubKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoin() types.Coin {
	return acc.Coin
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoin(coin types.Coin) error {
	acc.Coin = coin
	return nil
}

// GetAccountNumber - Implements Account
func (acc *BaseAccount) GetAccountNumber() uint64 {
	return acc.AccountNumber
}

// SetAccountNumber - Implements Account
func (acc *BaseAccount) SetAccountNumber(accNumber uint64) error {
	acc.AccountNumber = accNumber
	return nil
}

// GetSequence - Implements sdk.Account.
func (acc *BaseAccount) GetSequence() uint64 {
	return acc.Sequence
}

// SetSequence - Implements sdk.Account.
func (acc *BaseAccount) SetSequence(seq uint64) error {
	acc.Sequence = seq
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) types.Coin{
	return acc.GetCoin()
}

func (acc *BaseAccount) String() string {
	return fmt.Sprintf(`Vesting Account:
  Address:          %s
  Pubkey:           %s
  Coins:            %s
  AccountNumber:    %d
  Sequence:         %d`,
		acc.Address, acc.PubKey, acc.Coin, acc.AccountNumber, acc.Sequence,
	)
}

// EthAddress returns the account address ethereum format.
func (acc *BaseAccount) EthAddress() ethcmn.Address {
	return ethcmn.BytesToAddress(acc.Address.Bytes())
}

// Balance returns the balance of an account.
func (acc *BaseAccount) Balance(denom string) types.Int {
	return acc.GetCoin().AmountOf(denom)
}

// SetBalance sets an account's balance of the given coin denomination.
//
// CONTRACT: assumes the denomination is valid.
func (acc *BaseAccount) SetBalance(denom string, amt types.Int) {
	newCoin := types.NewCoin(amt)
	if err := acc.SetCoin(newCoin); err != nil {
		panic(fmt.Errorf("could not set %s coins for address %s: %w", denom, acc.EthAddress().String(), err))
	}
}

func (acc *BaseAccount) SetContractType(contractType string) error {
	if contractType != EvmContractType && contractType != WasmContractType {
		return errors.New("error contractType")
	}
	acc.ContractType = contractType
	return nil
}

func (acc *BaseAccount) GetContractType() string {
	return acc.ContractType
}