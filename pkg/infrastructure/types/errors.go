package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace  sdk.CodespaceType = "infrastructure"

	CodeCheckParamsError	CodeType = 101
)

func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", str)
}

func ErrCdcUnmarshalFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, err.Error())
}

func ErrMarshalFailed(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, err.Error())
}

func ErrGetInvalidResponse(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, err.Error())
}