package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

var (
	ErrInvalidDenomForTransfer = sdkerrors.Register(ModuleName, 2003, "invalid denomination for cross-chain transfer")
	ErrInvalidVersion          = sdkerrors.Register(ModuleName, 2004, "invalid ICS20 version")
	ErrInvalidAmount           = sdkerrors.Register(ModuleName, 2005, "invalid token amount")
	ErrTraceNotFound           = sdkerrors.Register(ModuleName, 2006, "denomination trace not found")
	ErrSendDisabled            = sdkerrors.Register(ModuleName, 2007, "fungible token transfers from this chain are disabled")
	ErrReceiveDisabled         = sdkerrors.Register(ModuleName, 2008, "fungible token transfers to this chain are disabled")
	ErrMaxTransferChannels     = sdkerrors.Register(ModuleName, 2009, "max transfer channels")
	ErrInvalidToken            = sdkerrors.Register(ModuleName, 2011, "invalid token")
	ErrInvalidSender           = sdkerrors.Register(ModuleName, 2012, "invalid sender")
	ErrInvalidReceiver         = sdkerrors.Register(ModuleName, 2013, "invalid receiver")
	ErrInvalidPortID           = sdkerrors.Register(ModuleName, 2014, "invalid port ID")
	ErrInvalidChannelID        = sdkerrors.Register(ModuleName, 2015, "invalid channel ID")
)

//func ErrInvalidParam(desc string) error {
//	return sdkerrors.Register(ModuleName, 2010, desc)
//}


//func ErrInvalidClient(desc string) error {
//	return sdkerrors.Register(ModuleName, 2011, desc)
//}