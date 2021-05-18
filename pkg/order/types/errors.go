package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 		       sdk.CodespaceType = "order"
	CodeCheckParamsError		   CodeType = 1501
	CodeBadMarshal  			   CodeType = 1502
	CodeQueryTxError			   CodeType = 1503
	CodeInvalidEndPoint            CodeType = 1504
	CodeGetOrderBookFailed         CodeType = 1505
)


func ErrCheckParams(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), desc)
}

func ErrFailedMarshal(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), desc)
}

func ErrQueryTx(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeCheckParamsError), desc)
}

func ErrInvalidEndPoint(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidEndPoint), desc)
}

func ErrGetOrderBookFailed(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeGetOrderBookFailed), desc)
}