package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

// IBC client sentinel errors
var (
	ErrClientExists                           = sdkerrors.Register(SubModuleName, 2042, "light client already exists")
	ErrInvalidClient                          = sdkerrors.Register(SubModuleName, 3, "light client is invalid")
	ErrClientNotFound                         = sdkerrors.Register(SubModuleName, 2044, "light client not found")
	ErrClientFrozen                           = sdkerrors.Register(SubModuleName, 2045, "light client is frozen due to misbehaviour")
	ErrInvalidClientMetadata                  = sdkerrors.Register(SubModuleName, 2046, "invalid client metadata")
	ErrConsensusStateNotFound                 = sdkerrors.Register(SubModuleName, 2047, "consensus state not found")
	ErrInvalidConsensus                       = sdkerrors.Register(SubModuleName, 2048, "invalid consensus state")
	ErrClientTypeNotFound                     = sdkerrors.Register(SubModuleName, 2049, "client type not found")
	ErrInvalidClientType                      = sdkerrors.Register(SubModuleName, 2050, "invalid client type")
	ErrRootNotFound                           = sdkerrors.Register(SubModuleName, 2051, "commitment root not found")
	ErrInvalidHeader                          = sdkerrors.Register(SubModuleName, 2052, "invalid client header")
	ErrInvalidMisbehaviour                    = sdkerrors.Register(SubModuleName, 2053, "invalid light client misbehaviour")
	ErrFailedClientStateVerification          = sdkerrors.Register(SubModuleName, 2054, "client state verification failed")
	ErrFailedClientConsensusStateVerification = sdkerrors.Register(SubModuleName, 2055, "client consensus state verification failed")
	ErrFailedConnectionStateVerification      = sdkerrors.Register(SubModuleName, 2056, "connection state verification failed")
	ErrFailedChannelStateVerification         = sdkerrors.Register(SubModuleName, 2057, "channel state verification failed")
	ErrFailedPacketCommitmentVerification     = sdkerrors.Register(SubModuleName, 2058, "packet commitment verification failed")
	ErrFailedPacketAckVerification            = sdkerrors.Register(SubModuleName, 2059, "packet acknowledgement verification failed")
	ErrFailedPacketReceiptVerification        = sdkerrors.Register(SubModuleName, 2060, "packet receipt verification failed")
	ErrFailedNextSeqRecvVerification          = sdkerrors.Register(SubModuleName, 2061, "next sequence receive verification failed")
	ErrSelfConsensusStateNotFound             = sdkerrors.Register(SubModuleName, 2062, "self consensus state not found")
	ErrUpdateClientFailed                     = sdkerrors.Register(SubModuleName, 2063, "unable to update light client")
	ErrInvalidUpdateClientProposal            = sdkerrors.Register(SubModuleName, 2064, "invalid update client proposal")
	ErrInvalidUpgradeClient                   = sdkerrors.Register(SubModuleName, 2065, "invalid client upgrade")
)
//
//func ErrInvalidClient(desc string) error {
//	return sdkerrors.Register(SubModuleName, 2066, desc)
//}

//func ErrInvalidParam(desc string) error {
//	return sdkerrors.Register(SubModuleName, 2067, desc)
//}