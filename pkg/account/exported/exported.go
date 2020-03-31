package exported

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"time"
)

type Account interface {
	GetAddress() types.AccAddress
	SetAddress(types.AccAddress) error // errors if already set.

	GetPubKey() crypto.PubKey // can return nil.
	SetPubKey(crypto.PubKey) error

	GetAccountNumber() uint64
	SetAccountNumber(uint64) error

	GetSequence() uint64
	SetSequence(uint64) error

	GetCoin() types.Coin
	SetCoin(types.Coin) error


	// Calculates the amount of coins that can be sent to other accounts given
	// the current time.
	SpendableCoins(blockTime time.Time) types.Coin
	// Ensure that account implements stringer
	String() string

	AddContract(contractAddress types.AccAddress)
	GetContractList() []string
}