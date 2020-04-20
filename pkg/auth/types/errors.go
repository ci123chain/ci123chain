package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace 				sdk.CodespaceType = "auth"
	CodeTxValidateBasicError       	CodeType = 501
)

func ErrTxValidateBasic(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeTxValidateBasicError,"Validate basic Error: %s", err.Error())
}
