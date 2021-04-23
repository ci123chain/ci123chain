package types

import (
	"context"
)

// MsgServer is the server API for Msg service.
type MsgServer interface {
	ValsetConfirm(context.Context, *MsgValsetConfirm) (*MsgValsetConfirmResponse, error)
	SendToEth(context.Context, *MsgSendToEth) (*MsgSendToEthResponse, error)
	RequestBatch(context.Context, *MsgRequestBatch) (*MsgRequestBatchResponse, error)
	ConfirmBatch(context.Context, *MsgConfirmBatch) (*MsgConfirmBatchResponse, error)
	ConfirmLogicCall(context.Context, *MsgConfirmLogicCall) (*MsgConfirmLogicCallResponse, error)
	DepositClaim(context.Context, *MsgDepositClaim) (*MsgDepositClaimResponse, error)
	WithdrawClaim(context.Context, *MsgWithdrawClaim) (*MsgWithdrawClaimResponse, error)
	ERC20DeployedClaim(context.Context, *MsgERC20DeployedClaim) (*MsgERC20DeployedClaimResponse, error)
	LogicCallExecutedClaim(context.Context, *MsgLogicCallExecutedClaim) (*MsgLogicCallExecutedClaimResponse, error)
	SetOrchestratorAddress(context.Context, *MsgSetOrchestratorAddress) (*MsgSetOrchestratorAddressResponse, error)
	CancelSendToEth(context.Context, *MsgCancelSendToEth) (*MsgCancelSendToEthResponse, error)
}