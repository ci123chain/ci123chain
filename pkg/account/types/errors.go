package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 			sdk.CodespaceType = "accounts"
	CodeSetAccountError       	CodeType = 701
)

func ErrSetAccount(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeSetAccountError,"Set Account Error:%s", err.Error())
}
