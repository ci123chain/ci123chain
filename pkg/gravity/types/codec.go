package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

var GravityCodec *codec.Codec
var SubModuleCdc *codec.ProtoCodec

func init() {
	GravityCodec = codec.New()
	RegisterCodec(GravityCodec)
	codec.RegisterCrypto(GravityCodec)
	//GravityCodec.Seal()
}

func SetBinary(registry codectypes.InterfaceRegistry) {
	SubModuleCdc = codec.NewProtoCodec(registry)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*EthereumClaim)(nil), nil)
	cdc.RegisterConcrete(&MsgSetOrchestratorAddress{}, "gravity/MsgSetOrchestratorAddress", nil)
	cdc.RegisterConcrete(&MsgValsetConfirm{}, "gravity/MsgValsetConfirm", nil)
	cdc.RegisterConcrete(&MsgSendToEth{}, "gravity/MsgSendToEth", nil)
	cdc.RegisterConcrete(&MsgRequestBatch{}, "gravity/MsgRequestBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmBatch{}, "gravity/MsgConfirmBatch", nil)
	cdc.RegisterConcrete(&MsgConfirmLogicCall{}, "gravity/MsgConfirmLogicCall", nil)
	cdc.RegisterConcrete(&Valset{}, "gravity/Valset", nil)
	cdc.RegisterConcrete(&MsgDepositClaim{}, "gravity/MsgDepositClaim", nil)
	cdc.RegisterConcrete(&MsgDeposit721Claim{}, "gravity/MsgDeposit721Claim", nil)
	cdc.RegisterConcrete(&MsgWithdrawClaim{}, "gravity/MsgWithdrawClaim", nil)
	cdc.RegisterConcrete(&MsgERC20DeployedClaim{}, "gravity/MsgERC20DeployedClaim", nil)
	cdc.RegisterConcrete(&MsgERC721DeployedClaim{}, "gravity/MsgERC721DeployedClaim", nil)
	cdc.RegisterConcrete(&MsgLogicCallExecutedClaim{}, "gravity/MsgLogicCallExecutedClaim", nil)
	cdc.RegisterConcrete(&MsgValsetConfirmNonceClaim{}, "gravity/MsgValsetConfirmNonceClaim", nil)
	cdc.RegisterConcrete(&OutgoingTxBatch{}, "gravity/OutgoingTxBatch", nil)
	cdc.RegisterConcrete(&MsgCancelSendToEth{}, "gravity/MsgCancelSendToEth", nil)
	cdc.RegisterConcrete(&OutgoingTransferTx{}, "gravity/OutgoingTransferTx", nil)
	cdc.RegisterConcrete(&ERC20Token{}, "gravity/ERC20Token", nil)
	cdc.RegisterConcrete(&ERC20ToDenom{}, "gravity/ERC20ToDenom", nil)
	cdc.RegisterConcrete(&IDSet{}, "gravity/IDSet", nil)
	cdc.RegisterConcrete(&Attestation{}, "gravity/Attestation", nil)
	cdc.RegisterConcrete(&MetaData{}, "gravity/ContractMetaData", nil)
	cdc.RegisterConcrete(&GenesisState{}, "gravity/GenesisState", nil)
	cdc.RegisterConcrete(&GravityData{}, "gravity/GravityData", nil)
	cdc.RegisterConcrete(&Params{}, "gravity/Params", nil)
	cdc.RegisterConcrete(&BridgeValidator{}, "gravity/BridgeValidator", nil)
	cdc.RegisterConcrete(&BridgeValidators{}, "gravity/BridgeValidators", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"gravity.v1.EthereumClaim",
		(*EthereumClaim)(nil),
	)
	registry.RegisterImplementations(
		(*EthereumClaim)(nil),
		&MsgDepositClaim{},
		&MsgDeposit721Claim{},
		&MsgWithdrawClaim{},
		&MsgERC20DeployedClaim{},
		&MsgERC721DeployedClaim{},
		&MsgValsetConfirmNonceClaim{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgValsetConfirm{},
		&MsgSendToEth{},
		&MsgRequestBatch{},
		&MsgConfirmBatch{},
		&MsgDepositClaim{},
		&MsgDeposit721Claim{},
		&MsgWithdrawClaim{},
		&MsgERC20DeployedClaim{},
		&MsgERC721DeployedClaim{},
		&MsgValsetConfirmNonceClaim{},
		&MsgCancelSendToEth{},
	)
}
