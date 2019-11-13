package transaction

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)


type CodeType = sdk.CodeType

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace 	sdk.CodespaceType = "transaction"
	CodeInvalidTx       CodeType = 101
	CodeInvalidTransfer CodeType = 102
	CodeInvalidSignature CodeType = 103
	CodeInvalidDeploy   CodeType = 104
	CodeInvalidCall     CodeType = 105
)

//----------------------------------------
// Error constructors

func ErrInvalidTx(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTx, "tx invalid")
}

func ErrInvalidTransfer(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTransfer, "transfer parameter error")
}

func ErrInvalidSignature(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSignature, "signature error")
}
