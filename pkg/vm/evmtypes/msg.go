package evmtypes

import (
	"errors"
	"math/big"
	"sync/atomic"

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

var (
	_ sdk.Msg = MsgEvmTx{}
)

// message types and route constants
const (
	// TypeMsgEvmTx defines the types string of an Ethereum tranasction
	TypeMsgEvmTx = "ethereum"
	RouteKey = "vm"
)

// MsgEthereumTx encapsulates an Ethereum transaction as an SDK message.
type MsgEvmTx struct {
	From sdk.AccAddress
	Data TxData

	// caches
	size atomic.Value
	from atomic.Value
}

// sigCache is used to cache the derived sender and contains the signer used
// to derive it.
type sigCache struct {
	signer ethtypes.Signer
	from   sdk.AccAddress
}

// NewMsgEvmTx returns a reference to a new Ethereum transaction message.
func NewMsgEvmTx(
	from sdk.AccAddress, nonce uint64, to *ethcmn.Address, amount *big.Int,
	gasLimit uint64, gasPrice *big.Int, payload []byte,
) MsgEvmTx {
	return newMsgEvmTx(from, nonce, to, amount, gasLimit, gasPrice, payload)
}

// NewMsgEthereumTxContract returns a reference to a new Ethereum transaction
// message designated for contract creation.
func NewMsgEthereumTxContract(
	from sdk.AccAddress, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, payload []byte,
) MsgEvmTx {
	return newMsgEvmTx(from, nonce, nil, amount, gasLimit, gasPrice, payload)
}

func newMsgEvmTx(
	from sdk.AccAddress, nonce uint64, to *ethcmn.Address, amount *big.Int,
	gasLimit uint64, gasPrice *big.Int, payload []byte,
) MsgEvmTx {
	if len(payload) > 0 {
		payload = ethcmn.CopyBytes(payload)
	}

	txData := TxData{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      payload,
		GasLimit:     gasLimit,
		Amount:       new(big.Int),
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}

	if amount != nil {
		txData.Amount.Set(amount)
	}
	if gasPrice != nil {
		txData.Price.Set(gasPrice)
	}

	return MsgEvmTx{From: from, Data: txData}
}

func (msg MsgEvmTx) String() string {
	return msg.Data.String()
}

// Route returns the route value of an MsgEthereumTx.
func (msg MsgEvmTx) Route() string { return RouteKey }

// Type returns the types value of an MsgEthereumTx.
func (msg MsgEvmTx) MsgType() string { return TypeMsgEvmTx }

func (msg MsgEvmTx) Bytes() []byte {
	bytes, err := ModuleCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg MsgEvmTx) GetFromAddress() sdk.AccAddress {
	//sc := msg.from.Load()
	//if sc == nil {
	//	return sdk.AccAddress{}
	//}
	//
	//sigCache := sc.(sigCache)
	//
	//if len(sigCache.from.Bytes()) == 0 {
	//	return sdk.AccAddress{}
	//}
	//
	//return sdk.ToAccAddress(sigCache.from.Bytes())
	return msg.From
}

// ValidateBasic implements the sdk.Msg interface. It performs basic validation
// checks of a Transaction. If returns an error if validation fails.
func (msg MsgEvmTx) ValidateBasic() sdk.Error {
	if msg.Data.Price.Cmp(big.NewInt(0)) == 0 {
		return ErrInvalidMsg(DefaultCodespace, errors.New("price is invalid"))
	}

	if msg.Data.Price.Sign() == -1 {
		return ErrInvalidMsg(DefaultCodespace, errors.New("price sign is invalid"))
	}

	// Amount can be 0
	if msg.Data.Amount.Sign() == -1 {
		return ErrInvalidMsg(DefaultCodespace, errors.New("amount is invalid"))
	}

	return nil
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg MsgEvmTx) To() *ethcmn.Address {
	return msg.Data.Recipient
}