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



// DecodeTx function is called by tendermint when node receives tx
func DecodeTx(bs []byte) (types.Tx, types.Error) {
	tx, err := decodeTx(bs)
	if err != nil {
		return nil, types.ErrTxDecode(err.Error())
	}
	return tx, nil
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

func decodeTx(bs []byte) (types.Tx, error) {

	var transfer Transaction
	err := transferCdc.UnmarshalBinaryLengthPrefixed(bs, &transfer)
	if err != nil {
		return nil, err
	}
	return transfer, nil


	//return DecodeTransferTx(bs)
	//code, err := FetchCodeValue(bs)
	//if err != nil {
	//	return nil, errors.New("fail to fetch tx code")
	//}
	//switch code {
	//case TRANSFER:
	//	return DecodeTransferTx(bs)
	//case CONTRACT_CALL:
	//	return DecodeContractCallTx(bs)
	//case CONTRACT_DEPLOY:
	//	return DecodeContractDeployTx(bs)
	//default:
	//	return nil, fmt.Errorf("unknown code '%v'", code)
	//}
}
