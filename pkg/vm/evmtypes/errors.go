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

	CodeCheckParamsError	CodeType = 1701

	CodeComputeGas          CodeType = 1702
	CodeInvalidState        CodeType = 1703
)

var (
	ErrComputeGas = sdkerrors.Register(string(DefaultCodespace), uint32(CodeComputeGas), "invalid intrinsic gas for transaction")
	ErrInvalidState = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidState), "invalid vm state")
	ErrInvalidParams = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), "invalid params")
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
//func ErrComputeGas(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeComputeGas), desc)
//}
////
//func ErrInvalidState(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidState), desc)
//}
//
//func ErrInvalidParams(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), desc)
//}
//
//func ErrChainConfigNotFound(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeErrChainConfigNotFound, msg)
//}
//
//func ErrTransitionDb(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeErrTransitionDb, msg)
//}
//
