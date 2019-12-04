package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)


type CodeType = sdk.CodeType

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace 		sdk.CodespaceType = "transaction"
	CodeInvalidTx       	CodeType = 101
	CodeInvalidTransfer		CodeType = 102
	CodeInvalidSignature 	CodeType = 103
	CodeBadPubkey			CodeType = 104
	CodeBadPrivkey			CodeType = 105
	CodeSetSequenceError	CodeType = 106
	CodeSendCoinError		CodeType = 107
)

//----------------------------------------
// Error constructors

func ErrInvalidTx(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTx, "tx invalid", str)
}

func ErrInvalidTransfer(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTransfer, "transfer parameter error", err)
}

func ErrSignature(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSignature, "signature error", err)
}

func ErrDecodePubkey(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBadPubkey, "Pubkey error", err)
}

func ErrDecodePrivkey(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeBadPrivkey, "Privkey error", err)
}

func ErrSetSequence(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeSetSequenceError, "Set sequence error", str)
}

func ErrSendCoin(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeSendCoinError, "Send coin to module error", err)
}