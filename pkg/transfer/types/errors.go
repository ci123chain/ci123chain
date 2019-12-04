package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType

// transfer errors reserve 200 ~ 299.
const (
	DefaultCodespace 		sdk.CodespaceType = "transfer"
	CodeInvalidAmount       CodeType = 201
	CodeInvalidReceiver     CodeType = 202
	CodeCheckParamsError	CodeType = 203
	CodeQueryTxError		CodeType = 204
)

//----------------------------------------
// Error constructors

func ErrBadAmount(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAmount, "param Amount invalid", err)
}

func ErrBadReceiver(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidReceiver, "param To invalid", err)
}

func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid", str)
}

func ErrQueryTx(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeQueryTxError, "query error", str)
}




