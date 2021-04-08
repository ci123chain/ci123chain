package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 			sdk.CodespaceType = "app"
	//CodeGenesisError       		CodeType = 401
	//CodeNewDBError       		CodeType = 402
	//CodeInitWithCfgError       	CodeType = 403
	//CodeTestNetError			CodeType = 404
)

func ErrGenesisFile(codespace sdk.CodespaceType, err error) error{
	return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("GenesisFile error: %s", err.Error()))
}

func ErrNewDB(codespace sdk.CodespaceType, err error) error{
	return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("New DB error: %s", err.Error()))
}

//func ErrInitWithCfg(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeInitWithCfgError,"Init with config error: %s", err.Error())
//}
//
//func ErrTestNet(codespace sdk.CodespaceType, err error) sdk.Error {
//	return sdk.NewError(codespace, CodeTestNetError,"Testnet error: %s", err.Error())
//}