package types
//
import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType
//
//// transfer errors reserve 200 ~ 299.
const (
	DefaultCodespace 		sdk.CodespaceType = "transfer"
	CodeInvalidAmount       CodeType = 201
	CodeInvalidReceiver     CodeType = 202
	CodeCheckParamsError	CodeType = 203
	CodeQueryTxError		CodeType = 204
	CodeGetNodeFailed       CodeType = 205
	CodeGetBlockFailed      CodeType = 206
)

var (
	ErrQueryTx = sdkerrors.Register(string(DefaultCodespace), uint32(CodeQueryTxError), "query Tx failed")
	ErrInvalidTxHash = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), "invalid tx hash")
	ErrGetNodeFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetNodeFailed), "get node failed")
	ErrGetBlockFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetBlockFailed), "get block with height failed")
)
//
////----------------------------------------
//// Error constructors
//
//func ErrBadAmount(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeInvalidAmount, "param Amount invalid: %s", err.Error())
//}
//
//func ErrBadReceiver(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeInvalidReceiver, "param To invalid: %s", err.Error())
//}
//
//func ErrCheckParams(codespace sdk.CodespaceType, str string) sdk.Error {
//	return sdk.NewError(codespace, CodeCheckParamsError, "param invalid: %s", str)
//}
//
//func ErrQueryTx(codespace sdk.CodespaceType, str string) error {
//	return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("query Tx failed: %v", str))
//}
//
//
//
//
