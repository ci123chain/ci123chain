package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 		sdk.CodespaceType = "order"
	CodeCheckParamsError	CodeType = 401
	CodeBadMarshal  CodeType = 402
	CodeQueryTxError		CodeType = 403
)


func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", str)
}

func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
}

func ErrQueryTx(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeQueryTxError, "query error: %s", str)
}

