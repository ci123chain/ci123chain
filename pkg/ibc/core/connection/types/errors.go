package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

// IBC connection sentinel errors
var (
	ErrConnectionExists              = sdkerrors.Register(SubModuleName, 2082, "connection already exists")
	ErrConnectionNotFound            = sdkerrors.Register(SubModuleName, 2083, "connection not found")
	ErrClientConnectionPathsNotFound = sdkerrors.Register(SubModuleName, 2084, "light client connection paths not found")
	ErrConnectionPath                = sdkerrors.Register(SubModuleName, 2085, "connection path is not associated to the given light client")
	ErrInvalidConnectionState        = sdkerrors.Register(SubModuleName, 2086, "invalid connection state")
	ErrInvalidCounterparty           = sdkerrors.Register(SubModuleName, 2087, "invalid counterparty connection")
	ErrInvalidConnection             = sdkerrors.Register(SubModuleName, 2088, "invalid connection")
	ErrInvalidVersion                = sdkerrors.Register(SubModuleName, 2089, "invalid connection version")
	ErrVersionNegotiationFailed      = sdkerrors.Register(SubModuleName, 2090, "connection version negotiation failed")
	ErrInvalidConnectionIdentifier   = sdkerrors.Register(SubModuleName, 2091, "invalid connection identifier")
)
func ErrInvalidParam(desc string) error {
	return sdkerrors.Register(SubModuleName, 2092, desc)
}
