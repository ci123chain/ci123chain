package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace  						sdk.CodespaceType = "mint"
	CodeBadMarshal  CodeType    =     600
)
//
//func ErrFailedMarshal(codespace sdk.CodespaceType, detailStr string) sdk.Error {
//	return sdk.NewError(codespace, CodeBadMarshal, "Marshal Error: %s", detailStr)
//}
