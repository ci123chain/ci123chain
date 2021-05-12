package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

// IBC port sentinel errors
var (
	ErrPortExists   = sdkerrors.Register(SubModuleName, 2102, "port is already binded")
	ErrPortNotFound = sdkerrors.Register(SubModuleName, 2103, "port not found")
	ErrInvalidPort  = sdkerrors.Register(SubModuleName, 2104, "invalid port")
	ErrInvalidRoute = sdkerrors.Register(SubModuleName, 2105, "route not found")
)
