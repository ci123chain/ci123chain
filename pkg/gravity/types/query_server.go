package types

import (
	"context"
)

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Deployments queries deployments
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	CurrentValset(context.Context, *QueryCurrentValsetRequest) (*QueryCurrentValsetResponse, error)
	ValsetRequest(context.Context, *QueryValsetRequestRequest) (*QueryValsetRequestResponse, error)
	ValsetConfirm(context.Context, *QueryValsetConfirmRequest) (*QueryValsetConfirmResponse, error)
	ValsetConfirmsByNonce(context.Context, *QueryValsetConfirmsByNonceRequest) (*QueryValsetConfirmsByNonceResponse, error)
	LastValsetRequests(context.Context, *QueryLastValsetRequestsRequest) (*QueryLastValsetRequestsResponse, error)
	LastPendingValsetRequestByAddr(context.Context, *QueryLastPendingValsetRequestByAddrRequest) (*QueryLastPendingValsetRequestByAddrResponse, error)
	LastPendingBatchRequestByAddr(context.Context, *QueryLastPendingBatchRequestByAddrRequest) (*QueryLastPendingBatchRequestByAddrResponse, error)
	LastPendingLogicCallByAddr(context.Context, *QueryLastPendingLogicCallByAddrRequest) (*QueryLastPendingLogicCallByAddrResponse, error)
	LastEventNonceByAddr(context.Context, *QueryLastEventNonceByAddrRequest) (*QueryLastEventNonceByAddrResponse, error)
	BatchFees(context.Context, *QueryBatchFeeRequest) (*QueryBatchFeeResponse, error)
	OutgoingTxBatches(context.Context, *QueryOutgoingTxBatchesRequest) (*QueryOutgoingTxBatchesResponse, error)
	OutgoingLogicCalls(context.Context, *QueryOutgoingLogicCallsRequest) (*QueryOutgoingLogicCallsResponse, error)
	BatchRequestByNonce(context.Context, *QueryBatchRequestByNonceRequest) (*QueryBatchRequestByNonceResponse, error)
	BatchConfirms(context.Context, *QueryBatchConfirmsRequest) (*QueryBatchConfirmsResponse, error)
	LogicConfirms(context.Context, *QueryLogicConfirmsRequest) (*QueryLogicConfirmsResponse, error)
	ERC20ToDenom(context.Context, *QueryERC20ToDenomRequest) (*QueryERC20ToDenomResponse, error)
	DenomToERC20(context.Context, *QueryDenomToERC20Request) (*QueryDenomToERC20Response, error)
	GetDelegateKeyByValidator(context.Context, *QueryDelegateKeysByValidatorAddress) (*QueryDelegateKeysByValidatorAddressResponse, error)
	GetDelegateKeyByEth(context.Context, *QueryDelegateKeysByEthAddress) (*QueryDelegateKeysByEthAddressResponse, error)
	GetDelegateKeyByOrchestrator(context.Context, *QueryDelegateKeysByOrchestratorAddress) (*QueryDelegateKeysByOrchestratorAddressResponse, error)
	GetPendingSendToEth(context.Context, *QueryPendingSendToEth) (*QueryPendingSendToEthResponse, error)
}

type QueryParamsRequest struct {
}

type QueryParamsResponse struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

type QueryCurrentValsetRequest struct {
}

type QueryCurrentValsetResponse struct {
	Valset *Valset `protobuf:"bytes,1,opt,name=valset,proto3" json:"valset,omitempty"`
}

func (m *QueryCurrentValsetResponse) GetValset() *Valset {
	if m != nil {
		return m.Valset
	}
	return nil
}

type QueryValsetRequestRequest struct {
	Nonce uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *QueryValsetRequestRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

type QueryValsetRequestResponse struct {
	Valset *Valset `protobuf:"bytes,1,opt,name=valset,proto3" json:"valset,omitempty"`
}

func (m *QueryValsetRequestResponse) GetValset() *Valset {
	if m != nil {
		return m.Valset
	}
	return nil
}

type QueryValsetConfirmRequest struct {
	Nonce   uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryValsetConfirmRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *QueryValsetConfirmRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryValsetConfirmResponse struct {
	Confirm *MsgValsetConfirm `protobuf:"bytes,1,opt,name=confirm,proto3" json:"confirm,omitempty"`
}

func (m *QueryValsetConfirmResponse) GetConfirm() *MsgValsetConfirm {
	if m != nil {
		return m.Confirm
	}
	return nil
}

type QueryValsetConfirmsByNonceRequest struct {
	Nonce uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (m *QueryValsetConfirmsByNonceRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

type QueryValsetConfirmsByNonceResponse struct {
	Confirms []*MsgValsetConfirm `protobuf:"bytes,1,rep,name=confirms,proto3" json:"confirms,omitempty"`
}

func (m *QueryValsetConfirmsByNonceResponse) GetConfirms() []*MsgValsetConfirm {
	if m != nil {
		return m.Confirms
	}
	return nil
}

type QueryLastValsetRequestsRequest struct {
}

type QueryLastValsetRequestsResponse struct {
	Valsets []*Valset `protobuf:"bytes,1,rep,name=valsets,proto3" json:"valsets,omitempty"`
}

func (m *QueryLastValsetRequestsResponse) GetValsets() []*Valset {
	if m != nil {
		return m.Valsets
	}
	return nil
}

type QueryLastPendingValsetRequestByAddrRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryLastPendingValsetRequestByAddrRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryLastPendingValsetRequestByAddrResponse struct {
	Valsets []*Valset `protobuf:"bytes,1,rep,name=valsets,proto3" json:"valsets,omitempty"`
}

func (m *QueryLastPendingValsetRequestByAddrResponse) GetValsets() []*Valset {
	if m != nil {
		return m.Valsets
	}
	return nil
}

type QueryBatchFeeRequest struct {
}

type QueryBatchFeeResponse struct {
	BatchFees []*BatchFees `protobuf:"bytes,1,rep,name=batch_fees,json=batchFees,proto3" json:"batch_fees,omitempty"`
}

func (m *QueryBatchFeeResponse) GetBatchFees() []*BatchFees {
	if m != nil {
		return m.BatchFees
	}
	return nil
}

type QueryLastPendingBatchRequestByAddrRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryLastPendingBatchRequestByAddrRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryLastPendingBatchRequestByAddrResponse struct {
	Batch *OutgoingTxBatch `protobuf:"bytes,1,opt,name=batch,proto3" json:"batch,omitempty"`
}

func (m *QueryLastPendingBatchRequestByAddrResponse) GetBatch() *OutgoingTxBatch {
	if m != nil {
		return m.Batch
	}
	return nil
}

type QueryLastPendingLogicCallByAddrRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryLastPendingLogicCallByAddrRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryLastPendingLogicCallByAddrResponse struct {
	Call *OutgoingLogicCall `protobuf:"bytes,1,opt,name=call,proto3" json:"call,omitempty"`
}

func (m *QueryLastPendingLogicCallByAddrResponse) GetCall() *OutgoingLogicCall {
	if m != nil {
		return m.Call
	}
	return nil
}

type QueryOutgoingTxBatchesRequest struct {
}

type QueryOutgoingTxBatchesResponse struct {
	Batches []*OutgoingTxBatch `protobuf:"bytes,1,rep,name=batches,proto3" json:"batches,omitempty"`
}

func (m *QueryOutgoingTxBatchesResponse) GetBatches() []*OutgoingTxBatch {
	if m != nil {
		return m.Batches
	}
	return nil
}

type QueryOutgoingLogicCallsRequest struct {
}

type QueryOutgoingLogicCallsResponse struct {
	Calls []*OutgoingLogicCall `protobuf:"bytes,1,rep,name=calls,proto3" json:"calls,omitempty"`
}

func (m *QueryOutgoingLogicCallsResponse) GetCalls() []*OutgoingLogicCall {
	if m != nil {
		return m.Calls
	}
	return nil
}

type QueryBatchRequestByNonceRequest struct {
	Nonce           uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	ContractAddress string `protobuf:"bytes,2,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
}

func (m *QueryBatchRequestByNonceRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *QueryBatchRequestByNonceRequest) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

type QueryBatchRequestByNonceResponse struct {
	Batch *OutgoingTxBatch `protobuf:"bytes,1,opt,name=batch,proto3" json:"batch,omitempty"`
}

func (m *QueryBatchRequestByNonceResponse) GetBatch() *OutgoingTxBatch {
	if m != nil {
		return m.Batch
	}
	return nil
}

type QueryBatchConfirmsRequest struct {
	Nonce           uint64 `protobuf:"varint,1,opt,name=nonce,proto3" json:"nonce,omitempty"`
	ContractAddress string `protobuf:"bytes,2,opt,name=contract_address,json=contractAddress,proto3" json:"contract_address,omitempty"`
}

func (m *QueryBatchConfirmsRequest) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *QueryBatchConfirmsRequest) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

type QueryBatchConfirmsResponse struct {
	Confirms []*MsgConfirmBatch `protobuf:"bytes,1,rep,name=confirms,proto3" json:"confirms,omitempty"`
}

func (m *QueryBatchConfirmsResponse) GetConfirms() []*MsgConfirmBatch {
	if m != nil {
		return m.Confirms
	}
	return nil
}

type QueryLogicConfirmsRequest struct {
	InvalidationId    []byte `protobuf:"bytes,1,opt,name=invalidation_id,json=invalidationId,proto3" json:"invalidation_id,omitempty"`
	InvalidationNonce uint64 `protobuf:"varint,2,opt,name=invalidation_nonce,json=invalidationNonce,proto3" json:"invalidation_nonce,omitempty"`
}

func (m *QueryLogicConfirmsRequest) GetInvalidationId() []byte {
	if m != nil {
		return m.InvalidationId
	}
	return nil
}

func (m *QueryLogicConfirmsRequest) GetInvalidationNonce() uint64 {
	if m != nil {
		return m.InvalidationNonce
	}
	return 0
}

type QueryLogicConfirmsResponse struct {
	Confirms []*MsgConfirmLogicCall `protobuf:"bytes,1,rep,name=confirms,proto3" json:"confirms,omitempty"`
}

func (m *QueryLogicConfirmsResponse) GetConfirms() []*MsgConfirmLogicCall {
	if m != nil {
		return m.Confirms
	}
	return nil
}

type QueryLastEventNonceByAddrRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *QueryLastEventNonceByAddrRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

type QueryLastEventNonceByAddrResponse struct {
	EventNonce uint64 `protobuf:"varint,1,opt,name=event_nonce,json=eventNonce,proto3" json:"event_nonce,omitempty"`
}

func (m *QueryLastEventNonceByAddrResponse) GetEventNonce() uint64 {
	if m != nil {
		return m.EventNonce
	}
	return 0
}

type QueryERC20ToDenomRequest struct {
	Erc20 string `protobuf:"bytes,1,opt,name=erc20,proto3" json:"erc20,omitempty"`
}

func (m *QueryERC20ToDenomRequest) GetErc20() string {
	if m != nil {
		return m.Erc20
	}
	return ""
}

type QueryERC20ToDenomResponse struct {
	Denom            string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	CosmosOriginated bool   `protobuf:"varint,2,opt,name=cosmos_originated,json=cosmosOriginated,proto3" json:"cosmos_originated,omitempty"`
}


func (m *QueryERC20ToDenomResponse) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *QueryERC20ToDenomResponse) GetCosmosOriginated() bool {
	if m != nil {
		return m.CosmosOriginated
	}
	return false
}

type QueryDenomToERC20Request struct {
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
}


func (m *QueryDenomToERC20Request) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

type QueryDenomToERC20Response struct {
	Erc20            string `protobuf:"bytes,1,opt,name=erc20,proto3" json:"erc20,omitempty"`
	CosmosOriginated bool   `protobuf:"varint,2,opt,name=cosmos_originated,json=cosmosOriginated,proto3" json:"cosmos_originated,omitempty"`
}

func (m *QueryDenomToERC20Response) GetErc20() string {
	if m != nil {
		return m.Erc20
	}
	return ""
}

func (m *QueryDenomToERC20Response) GetCosmosOriginated() bool {
	if m != nil {
		return m.CosmosOriginated
	}
	return false
}

type QueryDelegateKeysByValidatorAddress struct {
	ValidatorAddress string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
}

func (m *QueryDelegateKeysByValidatorAddress) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

type QueryDelegateKeysByValidatorAddressResponse struct {
	EthAddress          string `protobuf:"bytes,1,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty"`
	OrchestratorAddress string `protobuf:"bytes,2,opt,name=orchestrator_address,json=orchestratorAddress,proto3" json:"orchestrator_address,omitempty"`
}

func (m *QueryDelegateKeysByValidatorAddressResponse) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

func (m *QueryDelegateKeysByValidatorAddressResponse) GetOrchestratorAddress() string {
	if m != nil {
		return m.OrchestratorAddress
	}
	return ""
}

type QueryDelegateKeysByEthAddress struct {
	EthAddress string `protobuf:"bytes,1,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty"`
}

func (m *QueryDelegateKeysByEthAddress) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

type QueryDelegateKeysByEthAddressResponse struct {
	ValidatorAddress    string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	OrchestratorAddress string `protobuf:"bytes,2,opt,name=orchestrator_address,json=orchestratorAddress,proto3" json:"orchestrator_address,omitempty"`
}

func (m *QueryDelegateKeysByEthAddressResponse) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *QueryDelegateKeysByEthAddressResponse) GetOrchestratorAddress() string {
	if m != nil {
		return m.OrchestratorAddress
	}
	return ""
}

type QueryDelegateKeysByOrchestratorAddress struct {
	OrchestratorAddress string `protobuf:"bytes,1,opt,name=orchestrator_address,json=orchestratorAddress,proto3" json:"orchestrator_address,omitempty"`
}

func (m *QueryDelegateKeysByOrchestratorAddress) GetOrchestratorAddress() string {
	if m != nil {
		return m.OrchestratorAddress
	}
	return ""
}

type QueryDelegateKeysByOrchestratorAddressResponse struct {
	ValidatorAddress string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	EthAddress       string `protobuf:"bytes,2,opt,name=eth_address,json=ethAddress,proto3" json:"eth_address,omitempty"`
}

func (m *QueryDelegateKeysByOrchestratorAddressResponse) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

func (m *QueryDelegateKeysByOrchestratorAddressResponse) GetEthAddress() string {
	if m != nil {
		return m.EthAddress
	}
	return ""
}

type QueryPendingSendToEth struct {
	SenderAddress string `protobuf:"bytes,1,opt,name=sender_address,json=senderAddress,proto3" json:"sender_address,omitempty"`
}

func (m *QueryPendingSendToEth) GetSenderAddress() string {
	if m != nil {
		return m.SenderAddress
	}
	return ""
}

type QueryPendingSendToEthResponse struct {
	TransfersInBatches []*OutgoingTransferTx `protobuf:"bytes,1,rep,name=transfers_in_batches,json=transfersInBatches,proto3" json:"transfers_in_batches,omitempty"`
	UnbatchedTransfers []*OutgoingTransferTx `protobuf:"bytes,2,rep,name=unbatched_transfers,json=unbatchedTransfers,proto3" json:"unbatched_transfers,omitempty"`
}

func (m *QueryPendingSendToEthResponse) GetTransfersInBatches() []*OutgoingTransferTx {
	if m != nil {
		return m.TransfersInBatches
	}
	return nil
}

func (m *QueryPendingSendToEthResponse) GetUnbatchedTransfers() []*OutgoingTransferTx {
	if m != nil {
		return m.UnbatchedTransfers
	}
	return nil
}





































