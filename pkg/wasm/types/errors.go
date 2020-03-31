package types

import (
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
)


type CodeType = sdk.CodeType
const (
	DefaultCodespace sdk.CodespaceType = "wasm"
)

const (
	CodeCheckParamsError	CodeType = 50
	CodeInvalidMsgError     CodeType = 51
	CodeHandleMsgFailedError  CodeType = 52
	CodeSetSequenceFailedError CodeType = 53
)

func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid", str)
}

func ErrInvalidMsg(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeInvalidMsgError, "msg invalid", err)
}

func ErrCreateFailed(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeHandleMsgFailedError, "create failed", err)
}

/*func ErrCheckWasmCode(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeHandleMsgFailedError, "uncompress code failed", err)
}*/

func ErrInstantiateFailed(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeHandleMsgFailedError, "instantiate failed", err)
}


func ErrExecuteFailed(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeHandleMsgFailedError, "execute failed", err)
}

func ErrSetNewAccountSequence(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeSetSequenceFailedError, "set sequence of account failed", err)
}
