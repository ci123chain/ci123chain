package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)
type CodeType = sdk.CodeType
const (
	DefaultCodespace 			sdk.CodespaceType = "app"
	CodeGenesisError       		CodeType = 401
	CodeNewDBError       		CodeType = 402
	CodeInitWithCfgError       	CodeType = 403
	CodeTestNetError			CodeType = 404
)

func ErrGenesisFile(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeGenesisError,"GenesisFile error: %s", err.Error())
}

func ErrNewDB(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeNewDBError,"New DB error: %s", err.Error())
}

func ErrInitWithCfg(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInitWithCfgError,"Init with config error: %s", err.Error())
}

func ErrTestNet(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeTestNetError,"Testnet error: %s", err.Error())
}