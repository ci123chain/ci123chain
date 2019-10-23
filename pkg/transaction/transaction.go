package transaction

import (
	"CI123Chain/pkg/abci/types"
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pkg/errors"
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
}



// DecodeTx function is called by tendermint when node receives tx
func DecodeTx(bs []byte) (types.Tx, types.Error) {
	tx, err := decodeTx(bs)
	if err != nil {
		return nil, types.ErrTxDecode(err.Error())
	}
	return tx, nil
}


func decodeTx(bs []byte) (types.Tx, error) {
	code, err := FetchCodeValue(bs)
	if err != nil {
		return nil, errors.New("fail to fetch tx code")
	}
	switch code {
	case TRANSFER:
		return DecodeTransferTx(bs)
	//case CONTRACT_CALL:
	//	return DecodeContractCallTx(bs)
	//case CONTRACT_DEPLOY:
	//	return DecodeContractDeployTx(bs)
	default:
		return nil, fmt.Errorf("unknown code '%v'", code)
	}
}

// FetchCodeValue returns code from rlp bytes
func FetchCodeValue(bs []byte) (byte, error) {
	r := bytes.NewReader(bs)
	s := rlp.NewStream(r, uint64(len(bs)))
	if _, err := s.List(); err != nil {
		return 0, err
	}
	var code byte
	return code, s.Decode(&code)
}
