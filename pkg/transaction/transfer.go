package transaction

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/rlp"
)

var emptyAddr types.AccAddress

const RouteKey = "Transfer"

func NewTransferTx(from, to types.AccAddress, gas, nonce uint64, amount types.Coin, isFabric bool ) Transaction {
	tx := &TransferTx{
		Common: CommonTx{
			Code: TRANSFER,
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
	Common CommonTx
	To     types.AccAddress
	Amount types.Coin
	FabricMode bool
}

func DecodeTransferTx(b []byte) (*TransferTx, error) {
	tx := new(TransferTx)
	return tx, rlp.DecodeBytes(b, tx)
}

func isEmptyAddr(addr types.AccAddress) bool {
	return addr == emptyAddr
}

func (tx *TransferTx) SetPubKey(pub []byte) {
	tx.Common.PubKey = pub
}

func (tx *TransferTx) SetSignature(sig []byte) {
	tx.Common.SetSignature(sig)
}


func (tx *TransferTx) ValidateBasic() types.Error {
	if err := tx.Common.ValidateBasic(); err != nil {
		return err
	}
	if tx.Amount == 0 {
		return ErrInvalidTransfer(DefaultCodespace, "tx.Amount == 0")
	}
	if isEmptyAddr(tx.To) {
		return ErrInvalidTransfer(DefaultCodespace, "tx.To == empty")
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
	b, err := rlp.EncodeToBytes(tx)
	if err != nil {
		panic(err)
	}
	return b
}
