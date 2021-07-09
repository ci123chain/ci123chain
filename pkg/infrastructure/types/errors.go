package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace  sdk.CodespaceType = "infrastructure"

	CodeInvalidEndPoint	CodeType = 1901
	CodeCdcMarshalFailed CodeType = 1902
	CodeGetContentFailed CodeType = 1903
)

var (
	ErrInvalidEndPoint = sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), "invalid endpoint")
	ErrCdcMarshaFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), "cdc marshal failed")
	ErrGetContentFailed = sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetContentFailed), "get content failed")
)

//func ErrInvalidEndPoint(desc string) error {
////	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), desc)
////}
////
////func ErrCdcMarshaFailed(desc string) error {
////	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCdcMarshalFailed), desc)
////}
////
////func ErrGetContentFailed(desc string) error {
////	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetContentFailed), desc)
////}