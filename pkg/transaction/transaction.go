package transaction

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types"
)

// Transaction code
const (
	TRANSFER uint8 = 1 + iota
	CONTRACT_DEPLOY
	CONTRACT_CALL
)

type Transaction interface {
	GetSignBytes() []byte
	GetSignature() []byte
	SetSignature([]byte)
	Bytes() []byte
	SetPubKey([]byte)
	GetGas() uint64
	GetNonce() uint64
	GetFromAddress() types.AccAddress
}

