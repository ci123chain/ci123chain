package transaction

import "github.com/ci123chain/ci123chain/pkg/transaction/types"

var (
	ErrInvalidTx				= types.ErrInvalidTx
	ErrInvalidTransfer			= types.ErrInvalidTransfer
	ErrSignature				= types.ErrSignature
	ErrBadPubkey				= types.ErrDecodePubkey
	ErrBadPrivkey				= types.ErrDecodePrivkey
	ErrSetSequence				= types.ErrSetSequence
	ErrSendCoin					= types.ErrSendCoin
	ErrAmount					= types.ErrAmount
)
