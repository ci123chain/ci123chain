package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 			sdk.CodespaceType = "app"
	CodeGenesisError       		CodeType = 401
	CodeNewDBError       		CodeType = 402
	//CodeInitWithCfgError       	CodeType = 403
	//CodeTestNetError			CodeType = 404
	CodeInvalidParam         CodeType = 405
	CodeUnexpectedResponse     CodeType = 406
)

func ErrGenesisFile(codespace sdk.CodespaceType, err error) error{
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeGenesisError), err.Error())
}

func ErrNewDB(codespace sdk.CodespaceType, err error) error{
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeNewDBError), err.Error())
}

func ErrInvalidParam(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeInvalidParam), desc)
}

func ErrUnexpectedResponse(desc string) error {
	return sdkerrors.Register(string(DefaultCodespace), uint32(CodeUnexpectedResponse), desc)
}

//func ErrInitWithCfg(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeInitWithCfgError,"Init with configs error: %s", err.Error())
//}
//
//func ErrTestNet(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeTestNetError,"Testnet error: %s", err.Error())
//}