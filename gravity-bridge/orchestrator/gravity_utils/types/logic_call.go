package types

import "github.com/ethereum/go-ethereum/common"

type LogicCall struct {
	Transfers 			[]Erc20Token
	Fees 				[]Erc20Token
	LogicContractAddress common.Address
	PayLoad 			[]byte
	Timeout 			uint64
	InvalidationId 		[]byte
	InvalidationNonce 	uint64
}

type LogicCallConfirmResponse struct {
	InvalidationId 		[]byte
	InvalidationNonce 	uint64
	EthereumSigner 		common.Address
	Orchestrator		common.Address
	EthSignature		EthSignature
}

func (lcc LogicCallConfirmResponse) GetEthAddress() common.Address {
	return lcc.Orchestrator
}

func (lcc LogicCallConfirmResponse) GetSignature() EthSignature {
	return lcc.EthSignature
}