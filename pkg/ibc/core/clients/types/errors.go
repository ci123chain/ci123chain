package types

import sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

type CodeType = sdk.CodeType
const (
	DefaultCodespace  						sdk.CodespaceType = "ibc"
	CodeInvalidClientState 					CodeType = 800
	CodeInvalidConsensusState 				CodeType = 801
	CodeInvalidClientID						CodeType = 802
	CodeClientNotFound						CodeType = 803
	CodeInvalidClientType					CodeType = 804
	CodeErrInitClientState					CodeType = 805
)

func ErrorClientState(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidClientState, "param client_state invalid: %s", err.Error())
}


func ErrorConsensusState(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidConsensusState, "param consensus_state invalid: %s", err.Error())
}



func ErrorInvalidClientID(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidClientID, "param client_id invalid: %s", err.Error())
}

func ErrorClientNotFound(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeClientNotFound, "param client notfound: %s", err.Error())
}

func ErrInvalidClientType(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidClientType, "invalide client type: %s", err.Error())
}

func ErrInitClientState(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeErrInitClientState, "initial client_state err: %s", err.Error())
}