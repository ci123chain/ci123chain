package transfer

import (
	"errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/transfer/types"
	"github.com/tanhuiya/ci123chain/pkg/util"
)

const RouteKey = "Transfer"

func NewTransferTx(from, to sdk.AccAddress, gas, nonce uint64, amount sdk.Coin, isFabric bool ) transaction.Transaction {
	tx := &TransferTx{
		Common: transaction.CommonTx{
			From: from,
			Gas:  gas,
			Nonce:nonce,
		},
		To: 		to,
		Amount: 	amount,
		FabricMode: isFabric,
	}
	return tx
}

type TransferTx struct {
	Common transaction.CommonTx
	To     sdk.AccAddress
	Amount sdk.Coin
	FabricMode bool
}

//func DecodeTransferTx(b []byte) (*TransferTx, error) {
//	var transfer TransferTx
//	err := transferCdc.UnmarshalBinaryLengthPrefixed(b, &transfer)
//	if err != nil {
//		return nil, err
//	}
//	return &transfer, nil
//
//	//tx := new(TransferTx)
//	//return tx, rlp.DecodeBytes(b, tx)
//}

func (tx *TransferTx) SetPubKey(pub []byte) {
	tx.Common.PubKey = pub
}

func (tx *TransferTx) SetSignature(sig []byte) {
	tx.Common.SetSignature(sig)
}


func (tx *TransferTx) ValidateBasic() sdk.Error {
	if err := tx.Common.ValidateBasic(); err != nil {
		return err
	}
	if tx.Amount == 0 {
		return types.ErrBadAmount(types.DefaultCodespace, errors.New("amount = 0"))
	}
	if transaction.EmptyAddr(tx.To) {
		return types.ErrBadReceiver(types.DefaultCodespace, errors.New("empty to address"))
	}
	return tx.Common.VerifySignature(tx.GetSignBytes(), tx.FabricMode)
}

func (tx *TransferTx) Route() string {
	return RouteKey
}

func (tx *TransferTx) GetSignBytes() []byte {
	ntx := *tx
	ntx.SetSignature(nil)
	return util.TxHash(ntx.Bytes())
}

func (tx *TransferTx) Bytes() []byte {

	bytes, err := transferCdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *TransferTx) GetGas() uint64 {
	return tx.Common.Gas
}

func (tx *TransferTx) GetNonce() uint64 {
	return tx.Common.Nonce
}

func (tx *TransferTx) GetFromAddress() sdk.AccAddress {
	return tx.Common.From
}
