package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"math/big"
)


type ValSetUpdatedEvent struct {
	Nonce uint64
	Members []ValSetMember
}

func ValSetUpdatedEventFromLog(log *web3.Log) (ValSetUpdatedEvent, error){
	return ValSetUpdatedEvent{}, nil
}

type TransactionBatchExecutedEvent struct {
	BatchNonce uint64
	BlockHeight uint64
	Erc20 common.Address
	EventNonce uint64
}

func TransactionBatchExecutedEventFromLog(log *web3.Log) (TransactionBatchExecutedEvent, error){
	return TransactionBatchExecutedEvent{}, nil
}

type SendToCosmosEvent struct {
	Erc20 common.Address
	Sender common.Address
	// cosmos-address
	Destination common.Address
	Amount *big.Int
	EventNonce uint64
	BlockHeight uint64
}

func SendToCosmosEventFromLog(log *web3.Log) (SendToCosmosEvent, error){
	return SendToCosmosEvent{}, nil
}

type Erc20DeployedEvent struct {
	CosmosDenom string
	Erc20 common.Address
	Name string
	Symbol string
	Decimals uint8
	EventNonce uint64
	BlockHeight uint64
}

func Erc20DeployedEventFromLog(log *web3.Log) (Erc20DeployedEvent, error){
	return Erc20DeployedEvent{}, nil
}

type LogicCallExecutedEvent struct {
	InvalidationId []byte
	InvalidationNonce uint64
	ReturnData []byte
	EventNonce uint64
	BlockHeight uint64
}

func LogicCallExecutedEventFromLog(log *web3.Log) (LogicCallExecutedEvent, error){
	return LogicCallExecutedEvent{}, nil
}

