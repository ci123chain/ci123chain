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
}


// DefaultTxDecoder logic for standard transaction decoding
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

