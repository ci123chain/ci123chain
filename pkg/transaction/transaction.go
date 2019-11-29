package transaction

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
)

// Transaction code
const (
	TRANSFER uint8 = 1 + iota
	CONTRACT_DEPLOY
	CONTRACT_CALL
)

type Transaction interface {
	types.Tx
	GetSignBytes() []byte
	SetSignature([]byte)
	Bytes() []byte
	SetPubKey([]byte)
	GetGas() uint64
	GetNonce() uint64
	GetFromAddress() types.AccAddress
}


// DefaultTxDecoder logic for standard transfer decoding
func DefaultTxDecoder(cdc *codec.Codec) types.TxDecoder {
	return func(txBytes []byte) (types.Tx, types.Error) {
		var transfer Transaction
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &transfer)
		if err != nil {
			return nil, types.ErrTxDecode("decode msg failed").TraceSDK(err.Error())
		}
		return transfer, nil
	}
}

