package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 			sdk.CodespaceType = "accounts"
	//CodeSetAccountError       	CodeType = 701
	//CodeGetAccountError       	CodeType = 701
)

//func ErrSetAccount(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeSetAccountError,"Set Account Error:%s", err.Error())
//}
//
//func ErrGetAccount(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeGetAccountError,"Get Account Error:%s", err.Error())
//}
