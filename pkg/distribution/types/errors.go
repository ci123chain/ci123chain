package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace  sdk.CodespaceType = "distribution"

	CodeInvalidHeight CodeType = 300

	CodeBadMarshal  CodeType = 301
)


func ErrBadHeight(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHeight, "param Height invalid: %s", err.Error())
}

func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
}