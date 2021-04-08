package evmtypes
//
import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
//
type CodeType = sdk.CodeType
const (
	DefaultCodespace sdk.CodespaceType = "evm"
)
//
//const (
//	CodeCheckParamsError	CodeType = 60
//	CodeInvalidMsgError     CodeType = 61
//	CodeHandleMsgFailedError  CodeType = 62
//	CodeSetSequenceFailedError CodeType = 63
//	CodeInvalidAddress        CodeType  = 64
//	CodeComputeGasFailedError  CodeType  = 65
//	CodeErrInvalidState			CodeType  = 66
//	CodeErrChainConfigNotFound CodeType = 67
//	CodeErrTransitionDb		CodeType = 68
//)
//
//func ErrCheckParams(codespace sdk.CodespaceType, keyname string) sdk.Error {
//	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", keyname)
//}
//
//func ErrInvalidMsg(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeInvalidMsgError, "msg invalid: %s", err.Error())
//}
//
//func ErrSetNewAccountSequence(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeSetSequenceFailedError, "set sequence of account failed: %s", err.Error())
//}
//
//func ErrInvalidAddress(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeInvalidAddress, msg)
//}
//
func ErrComputeGas(codespce sdk.CodespaceType, msg string) error {
	return sdkerrors.Wrap(sdkerrors.ErrInternal, msg)
}
//
func ErrInvalidState(codespce sdk.CodespaceType, msg string) error {
	return sdkerrors.Wrap(sdkerrors.ErrInternal, msg)
}
//
//func ErrChainConfigNotFound(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeErrChainConfigNotFound, msg)
//}
//
//func ErrTransitionDb(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeErrTransitionDb, msg)
//}
//
