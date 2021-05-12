package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)


type CodeType = sdk.CodeType
const (
	DefaultCodespace sdk.CodespaceType = "wasm"
)

const (
	CodeCheckParamsError	CodeType = 1750
	CodeInvalidEndPoint     CodeType = 1751
	//CodeInvalidMsgError     CodeType = 1751
	//CodeHandleMsgFailedError  CodeType = 1752
	//CodeSetSequenceFailedError CodeType = 1753
	//CodeInvalidAddress        CodeType  = 1754
	//CodeQueryCodeInfoFailed   CodeType = 1755
	CodeCdcUnMarshalFailed      CodeType = 1756
	CodeQueryFailed             CodeType = 1757
	CodeJsonUnmarshalFailed     CodeType = 1758
	CodeCdcMarshalFailed        CodeType = 1759
	CodeGetBlockBloomFailed     CodeType = 1760
)

func ErrInvalidParams(desc string) error {
	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), desc)
}

func ErrInvalidEndPoint(desc string) error {
	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), desc)
}

//func ErrQueryCodeInfo(desc string) error {
//	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeQueryCodeInfoFailed), desc)
//}
//
func ErrCdcUnMarshalFailed(desc string) error {
	return  sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcUnMarshalFailed), desc)
}

//func ErrUploadFailed(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeHandleMsgFailedError, "upload failed: %s", err.Error())
//}

//func ErrUninstallFailed(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeHandleMsgFailedError, "uninstall failed: %s", err.Error())
//}

/*func ErrCheckWasmCode(codespce sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespce, CodeHandleMsgFailedError, "uncompress code failed", err)
}*/
//
//func ErrMigrateFailed(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeHandleMsgFailedError, "migrate failed: %s", err.Error())
//}

//func ErrInstantiateFailed(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeHandleMsgFailedError, "instantiate failed: %s", err.Error())
//}

//func ErrExecuteFailed(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeHandleMsgFailedError, "execute failed: %s", err.Error())
//}

func ErrQueryFailed(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeQueryFailed), desc)
}

func ErrJsonUnmarshalFailed(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeJsonUnmarshalFailed), desc)
}

func ErrCdcMarshalFailed(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), desc)
}

func ErrGetBlockBloomFailed(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetBlockBloomFailed), desc)
}

//func ErrSetNewAccountSequence(codespce sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespce, CodeSetSequenceFailedError, "set sequence of account failed: %s", err.Error())
//}

//func ErrInvalidAddress(codespce sdk.CodespaceType, msg string) sdk.Error {
//	return sdk.NewError(codespce, CodeInvalidAddress, msg)
//}