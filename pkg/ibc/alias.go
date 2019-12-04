package ibc

import (
	"github.com/tanhuiya/ci123chain/pkg/ibc/client/rest"
	"github.com/tanhuiya/ci123chain/pkg/ibc/handler"
	"github.com/tanhuiya/ci123chain/pkg/ibc/keeper"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
)

var (
	StoreKey  = types.StoreKey
	RouterKey  = types.RouterKey
	ModuleName = types.ModuleName

	NewHandler = handler.NewHandler
	NewKeeper = keeper.NewIBCKeeper

	NewIBCTransfer = types.NewIBCTransferMsg
	NewApplyIBCTx = types.NewApplyIBCTx
	NewIBCMsgBankSendMsg = types.NewIBCMsgBankSendMsg
	NewIBCReceiveReceiptMsg = types.NewIBCReceiveReceiptMsg

	RegisterCodec = types.RegisterCodec
	RegisterRoutes = rest.RegisterTxRoutes
	NewQuerier = keeper.NewQuerier

	ErrBadBankSignature       	= types.ErrBadBankSignature
	ErrBadReceiptSignature		= types.ErrBadReceiptSignature
	ErrBadUnmarshal      		= types.ErrFailedUnmarshal
	ErrBadMarshal      			= types.ErrFailedMarshal
	ErrGetBankAddr				= types.ErrGetBankAddr
	ErrMakeIBCMsg				= types.ErrMakeIBCMsg
	ErrSetIBCMsg				= types.ErrSetIBCMsg
	ErrApplyIBCMsg				= types.ErrApplyIBCMsg
	ErrMakeBankReceipt			= types.ErrMakeBankReceipt
	ErrBankSend					= types.ErrBankSend
	ErrReceiveReceipt			= types.ErrReceiveReceipt
	ErrState					= types.ErrState
)


