package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

var GravityCodec *codec.Codec
var SubModuleCdc *codec.ProtoCodec

func init(){
	GravityCodec = codec.New()
	RegisterCodec(GravityCodec)
	codec.RegisterCrypto(GravityCodec)
	GravityCodec.Seal()
}

func SetBinary(registry codectypes.InterfaceRegistry) {
	SubModuleCdc = codec.NewProtoCodec(registry)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterInterface((*EthereumClaim)(nil), nil)
	cdc.RegisterConcrete(&MsgSetOrchestratorAddress{}, "gravity/MsgSetOrchestratorAddress", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "gravity/MsgValsetConfirm", nil)
	cdc.RegisterConcrete(&MsgSendToEth{}, "gravity/MsgSendToEth", nil)
	cdc.RegisterConcrete(&MsgRequestBatch{}, "gravity/MsgRequestBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, "gravity/MsgConfirmBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmLogicCall{}, "gravity/MsgConfirmLogicCall", nil)
	cdc.RegisterConcrete(&Valset{}, "gravity/Valset", nil)
	cdc.RegisterConcrete(&MsgDepositClaim{}, "gravity/MsgDepositClaim", nil)
	cdc.RegisterConcrete(&MsgWithdrawClaim{}, "gravity/MsgWithdrawClaim", nil)
	cdc.RegisterConcrete(&MsgERC20DeployedClaim{}, "gravity/MsgERC20DeployedClaim", nil)
	cdc.RegisterConcrete(&MsgLogicCallExecutedClaim{}, "gravity/MsgLogicCallExecutedClaim", nil)
	cdc.RegisterConcrete(&OutgoingTxBatch{}, "gravity/OutgoingTxBatch", nil)
	cdc.RegisterConcrete(&MsgCancelSendToEth{}, "gravity/MsgCancelSendToEth", nil)
	cdc.RegisterConcrete(&OutgoingTransferTx{}, "gravity/OutgoingTransferTx", nil)
	cdc.RegisterConcrete(&ERC20Token{}, "gravity/ERC20Token", nil)
	cdc.RegisterConcrete(&IDSet{}, "gravity/IDSet", nil)
	cdc.RegisterConcrete(&Attestation{}, "gravity/Attestation", nil)
	cdc.RegisterConcrete(&MetaData{}, "gravity/ContractMetaData", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"gravity.v1.EthereumClaim",
		(*EthereumClaim)(nil),
	)
	registry.RegisterImplementations(
		(*EthereumClaim)(nil),
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
		&MsgERC20DeployedClaim{},
		&MsgLogicCallExecutedClaim{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSetOrchestratorAddress{},
		&MsgValsetConfirm{},
		&MsgSendToEth{},
		&MsgRequestBatch{},
		&MsgConfirmBatch{},
		&MsgConfirmLogicCall{},
		&MsgDepositClaim{},
		&MsgWithdrawClaim{},
		&MsgERC20DeployedClaim{},
		&MsgLogicCallExecutedClaim{},
		&MsgCancelSendToEth{},
	)
}