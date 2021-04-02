package errors

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type CodeType = sdk.CodeType
const (
	DefaultCodespace          sdk.CodespaceType = "ibc"
	CodeInvalidClientState    CodeType          = 800
	CodeInvalidConsensusState CodeType          = 801
	CodeInvalidClientID       CodeType          = 802
	CodeClientNotFound        CodeType          = 803
	CodeInvalidClientType     CodeType          = 804
	CodeErrInitClientState    CodeType          = 805

	CodeInvalidConnectionID       CodeType = 806
	CodeInvalidConnectionVersion  CodeType = 807
	CodeInvalidCounterpartyPrefix CodeType = 808
	CodeInvalidCounterparty       CodeType = 809
	CodeInvalidChannel			  CodeType = 810
	CodeInvalidChannelState 	  CodeType = 811
	CodeInvalidChannelOrder		  CodeType = 812
	CodeInvalidCounterpartyChannel CodeType = 813
)

func ErrorClientState(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidClientState, "param client_state invalid: %s", err.Error())
}


func ErrorConsensusState(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidConsensusState, "param consensus_state invalid: %s", err.Error())
}

func ErrorClientNotFound(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeClientNotFound, "param client notfound: %s", err.Error())
}

func ErrInvalidClientType(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidClientType, "invalid client types: %s", err.Error())
}

func ErrInitClientState(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeErrInitClientState, "initial client_state err: %s", err.Error())
}


func ErrorConnectionID(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidConnectionID, "invalid connectionid: %s", err.Error())
}


func ErrorCounterpartyConnectionID(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidConnectionID, "invalid counterparty connectionid: %s", err.Error())
}

func ErrorInvalidConnectionVersion(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidConnectionVersion, "invalid connection version: %s", err.Error())

}

func ErrorCounterpartyPrefix(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidCounterpartyPrefix, "invalid counterparty prefix: %s", err.Error())

}

func ErrorCounterparty(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidCounterparty, "invalid counterparty: %s", err.Error())
}


func ErrInvalidChannel(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidChannel, "invalid channel: %s", err.Error())
}


func ErrInvalidChannelState (cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidChannelState , "invalid channel state: %s", err.Error())
}

func ErrInvalidChannelOrder (cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidChannelOrder , "invalid channel order: %s", err.Error())
}

func ErrInvalidCounterpartyChannel(cs sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(cs, CodeInvalidCounterpartyChannel , "invalid counter party channel: %s", err.Error())

}

var (
	ErrInvalidPacket = errors.New("invalid packet")
)