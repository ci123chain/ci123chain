package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 		sdk.CodespaceType = "upgrade"
	CodeCheckParamsError	CodeType = 203
)


func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid", str)
}