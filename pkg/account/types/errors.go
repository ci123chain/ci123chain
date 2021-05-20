package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 			sdk.CodespaceType = "accounts"

	CodeAccountNotExisted      CodeType = 1101
	CodeInsufficientFunds      CodeType = 1102
	CodeInvalidParam           CodeType = 1103
	CodeInvalidEndPoint        CodeType = 1104
	CodeCdcMarshalFailed       CodeType = 1105
	CodeCdcUnmarshalFailed     CodeType = 1106
)

var (
	ErrAccountNotExisted = sdkerrors.Register(string(DefaultCodespace), uint32(CodeAccountNotExisted), "account not exist")
	ErrInsufficientFunds = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInsufficientFunds), "insufficient funds")
	ErrInvalidParam = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidParam), "invalid params")
	ErrInvalidEndPoint = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), "invalid endpoint")
	ErrCdcMarshalFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), "cdc marshal failed")
	ErrCdcUnmarshalFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcUnmarshalFailed), "cdc unmarshal failed")
)

//func ErrAccountNotExisted(desc string) error {
//	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeAccountNotExisted), desc)
//}

//func ErrInsufficientFunds(desc string) error {
//	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInsufficientFunds), desc)
//}

//func ErrInvalidParam (desc string) error {
//	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidParam), desc)
//}

//func ErrInvalidEndPoint(desc string) error {
//	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), desc)
//}

//func ErrCdcMarshalFailed(desc string) error {
//	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), desc)
//}
//
//func ErrCdcUnmarshalFailed(desc string) error {
//	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcUnmarshalFailed), desc)
//}