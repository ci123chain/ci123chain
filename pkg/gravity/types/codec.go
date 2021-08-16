package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

var GravityCodec *codec.Codec

func init(){
	GravityCodec = codec.New()
	RegisterCodec(GravityCodec)
	codec.RegisterCrypto(GravityCodec)
	GravityCodec.Seal()
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
