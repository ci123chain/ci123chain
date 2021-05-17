package types

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"math"
	"math/big"
	"sort"
)

type ValSetUpdatedEvent struct {
	Nonce uint64
	Members []ValSetMember
}

func ValSetUpdatedEventFromLog(log *web3.Log) (ValSetUpdatedEvent, error) {
	if len(log.Topics) < 2 {
		return ValSetUpdatedEvent{},  errors.New("Too few topics")
	}

	var nonceBz []byte
	nonceBz, err := log.Topics[1].MarshalText()
	if err != nil {
		return ValSetUpdatedEvent{}, err
	}

	nonceB, err := hex.DecodeString(string(nonceBz[2:]))
	x := new(big.Int)
	y := new(big.Int)
	nonceBig := x.SetBytes(nonceB)
	maxU64 := y.SetUint64(math.MaxUint64)
	if nonceBig.Cmp(maxU64) > 0 {
		return ValSetUpdatedEvent{}, errors.New("Nonce overflow, probably incorrect parsing")
	}
	nonce := nonceBig.Uint64()

	var indexStart uint64 = 2 * 32
	indexEnd := indexStart + 32
	ethAddressOffset := indexStart + 32
	lenEthAddressesBig := x.SetBytes(log.Data[indexStart:indexEnd])
	if lenEthAddressesBig.Cmp(maxU64) > 0 {
		return ValSetUpdatedEvent{}, errors.New("Ethereum array len overflow, probably incorrect parsing")
	}
	lenEthAddresses := lenEthAddressesBig.Uint64()

	indexStart = (3 + lenEthAddresses) * 32
	indexEnd = indexStart + 32
	powersOffset := indexStart + 32
	lenPowersBig := x.SetBytes(log.Data[indexStart:indexEnd])
	if lenPowersBig.Cmp(maxU64) > 0 {
		return ValSetUpdatedEvent{}, errors.New("Powers array len overflow, probably incorrect parsing")
	}
	lenPowers := lenPowersBig.Uint64()
	if lenPowers != lenEthAddresses {
		return ValSetUpdatedEvent{}, errors.New("Array len mismatch, probably incorrect parsing")
	}

	var validators []ValSetMember
	var i uint64
	for ; i < lenEthAddresses ; i++ {
		powerStart := (i * 32) + powersOffset
		powerEnd := powerStart + 32
		addressStart := (i * 32) + ethAddressOffset
		addressEnd := addressStart + 32
		powerBig := x.SetBytes(log.Data[powerStart:powerEnd])
		ethAddress := common.BytesToAddress(log.Data[addressStart + 12:addressEnd])
		if powerBig.Cmp(maxU64) > 0 {
			return ValSetUpdatedEvent{}, errors.New("Power overflow, probably incorrect parsing")
		}
		power := powerBig.Uint64()
		validators = append(validators, ValSetMember{
			Power:      power,
			EthAddress: &ethAddress,
		})
	}

	check := validators
	for i, j := 0, len(check)-1; i < j; i, j = i+1, j-1 {
		check[i], check[j] = check[j], check[i]
	}

	if !sort.SliceIsSorted(check, func(i, j int) bool {
		if check[i].Power < check[j].Power {
			return true
		} else if check[i].Power == check[j].Power {
			if bytes.Compare(check[i].EthAddress.Bytes(), check[j].EthAddress.Bytes()) < 0 {
				return true
			}
		}
		return false
	}) {
		logger.GetLogger().Error(fmt.Sprintf("Someone submitted an unsorted validator set, this means all updates will fail until someone feeds in this unsorted value by hand %v instead of %v",
			validators, check))
	}

	return ValSetUpdatedEvent{
		Nonce:   nonce,
		Members: validators,
	}, nil
}

func ValSetUpdatedEventFromLogs(log []*web3.Log) ([]ValSetUpdatedEvent, error) {
	var res []ValSetUpdatedEvent
	for _, x := range log {
		event, err := ValSetUpdatedEventFromLog(x)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}
	return res, nil
}

type TransactionBatchExecutedEvent struct {
	BatchNonce uint64
	BlockHeight uint64
	Erc20 common.Address
	EventNonce uint64
}

func TransactionBatchExecutedEventFromLog(log *web3.Log) (TransactionBatchExecutedEvent, error) {
	if len(log.Topics) < 3 {
		return TransactionBatchExecutedEvent{}, errors.New("Too few topics")
	}
	x := new(big.Int)

	var batchNonceBz []byte
	batchNonceBz, err := log.Topics[1].MarshalText()
	if err != nil {
		return TransactionBatchExecutedEvent{}, err
	}
	nonceB, err := hex.DecodeString(string(batchNonceBz[2:]))
	batchNonce := x.SetBytes(nonceB)

	var erc20Bz []byte
	erc20Bz, err = log.Topics[2].MarshalText()
	if err != nil {
		return TransactionBatchExecutedEvent{}, err
	}
	erc20 := common.BytesToAddress(erc20Bz[14:34])

	eventNonce := x.SetBytes(log.Data)
	blockHeight := log.BlockNumber
	if blockHeight == 0 {
		return TransactionBatchExecutedEvent{}, errors.New("Log does not have block number, we only search logs already in blocks?")
	}

	maxU64 := x.SetUint64(math.MaxUint64)
	if eventNonce.Cmp(maxU64) > 0 || batchNonce.Cmp(maxU64) > 0 {
		return TransactionBatchExecutedEvent{}, errors.New("Event nonce overflow, probably incorrect parsing")
	}

	return TransactionBatchExecutedEvent{
		BatchNonce:  batchNonce.Uint64(),
		BlockHeight: blockHeight,
		Erc20:       erc20,
		EventNonce:  eventNonce.Uint64(),
	}, nil
}

func TransactionBatchExecutedEventFromLogs(log []*web3.Log) ([]TransactionBatchExecutedEvent, error) {
	var res []TransactionBatchExecutedEvent
	for _, x := range log {
		event, err := TransactionBatchExecutedEventFromLog(x)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}
	return res, nil
}

func TransactionBatchExecutedEventFilterByEventNonce(lastEventNonce uint64, withdraws []TransactionBatchExecutedEvent) []TransactionBatchExecutedEvent {
	var res []TransactionBatchExecutedEvent
	for _, x := range withdraws {
		if x.EventNonce > lastEventNonce {
			res = append(res, x)
		}
	}
	return res
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

func SendToCosmosEventFromLog(log *web3.Log) (SendToCosmosEvent, error) {
	if len(log.Topics) < 4 {
		return SendToCosmosEvent{}, errors.New("Too few topics")
	}
	x := new(big.Int)

	erc20Bz, err := log.Topics[1].MarshalText()
	if err != nil {
		return SendToCosmosEvent{}, err
	}
	erc20 := common.BytesToAddress(erc20Bz[14:34])

	senderBz, err := log.Topics[2].MarshalText()
	if err != nil {
		return SendToCosmosEvent{}, err
	}
	sender := common.BytesToAddress(senderBz[14:34])

	destinationBz, err := log.Topics[3].MarshalText()
	if err != nil {
		return SendToCosmosEvent{}, err
	}
	destination := common.BytesToAddress(destinationBz[14:34])

	amount := x.SetBytes(log.Data[2:34])
	eventNonce := x.SetBytes(log.Data[34:])
	blockHeight := log.BlockNumber
	if blockHeight == 0 {
		return SendToCosmosEvent{}, errors.New("Log does not have block number, we only search logs already in blocks?")
	}

	maxU64 := x.SetUint64(math.MaxUint64)
	if eventNonce.Cmp(maxU64) > 0 {
		return SendToCosmosEvent{}, errors.New("Event nonce overflow, probably incorrect parsing")
	}

	return SendToCosmosEvent{
		Erc20:       erc20,
		Sender:      sender,
		Destination: destination,
		Amount:      amount,
		EventNonce:  eventNonce.Uint64(),
		BlockHeight: blockHeight,
	}, nil
}

func SendToCosmosEventFromLogs(log []*web3.Log) ([]SendToCosmosEvent, error) {
	var res []SendToCosmosEvent
	for _, x := range log {
		event, err := SendToCosmosEventFromLog(x)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}
	return res, nil
}

func SendToCosmosEventFilterByEventNonce(lastEventNonce uint64, deposits []SendToCosmosEvent) []SendToCosmosEvent {
	var res []SendToCosmosEvent
	for _, x := range deposits {
		if x.EventNonce > lastEventNonce {
			res = append(res, x)
		}
	}
	return res
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

func Erc20DeployedEventFromLog(log *web3.Log) (Erc20DeployedEvent, error) {
	if len(log.Topics) < 2 {
		return Erc20DeployedEvent{}, errors.New("Too few topics")
	}

	erc20Bz, err := log.Topics[1].MarshalText()
	if err != nil {
		return Erc20DeployedEvent{}, err
	}
	erc20 := common.BytesToAddress(erc20Bz[14:34])

	x := new(big.Int)
	maxU8 := x.SetUint64(math.MaxUint8)
	maxU32 := x.SetUint64(math.MaxUint32)
	maxU64 := x.SetUint64(math.MaxUint64)

	indexStart := 3 * 32
	indexEnd := indexStart + 32
	decimal := x.SetBytes(log.Data[indexStart:indexEnd])
	if decimal.Cmp(maxU8) > 0 {
		return Erc20DeployedEvent{}, errors.New("Decimals overflow, probably incorrect parsing")
	}

	indexStart = 4 * 32
	indexEnd = indexStart + 32
	nonce := x.SetBytes(log.Data[indexStart:indexEnd])
	if nonce.Cmp(maxU64) > 0 {
		return Erc20DeployedEvent{}, errors.New("Nonce overflow, probably incorrect parsing")
	}

	indexStart = 5 * 32
	indexEnd = indexStart + 32
	denomLen := x.SetBytes(log.Data[indexStart:indexEnd])
	if denomLen.Cmp(maxU32) > 0 {
		return Erc20DeployedEvent{}, errors.New("Denom length overflow, probably incorrect parsing")
	}

	indexStart = 6 * 32
	indexEnd = indexStart + int(denomLen.Uint64())
	denom := string(log.Data[indexStart:indexEnd])

	lg := logger.GetLogger()
	lg.Info(fmt.Sprintf("Denom: %s", denom))

	indexStart = ((indexEnd + 31) / 32) * 32
	indexEnd = indexStart + 32
	erc20NameLen := x.SetBytes(log.Data[indexStart:indexEnd])
	if erc20NameLen.Cmp(maxU32) > 0 {
		return Erc20DeployedEvent{}, errors.New("Erc20 name length overflow, probably incorrect parsing")
	}

	indexStart = indexEnd
	indexEnd = indexStart + int(erc20NameLen.Uint64())
	erc20Name := string(log.Data[indexStart:indexEnd])
	lg.Info(fmt.Sprintf("Erc20Name: %s", erc20Name))

	indexStart = ((indexEnd + 31) / 32) * 32
	indexEnd = indexStart + 32
	symbolLen := x.SetBytes(log.Data[indexStart:indexEnd])
	if symbolLen.Cmp(maxU32) > 0 {
		return Erc20DeployedEvent{}, errors.New("Symbol length overflow, probably incorrect parsing")
	}

	indexStart = indexEnd
	indexEnd = indexStart + int(symbolLen.Uint64())
	symbol := string(log.Data[indexStart:indexEnd])
	lg.Info(fmt.Sprintf("Symbol: %s", symbol))

	blockHeight := log.BlockNumber
	if blockHeight == 0 {
		return Erc20DeployedEvent{}, errors.New("Log does not have block number, we only search logs already in blocks?")
	}

	return Erc20DeployedEvent{
		CosmosDenom: denom,
		Erc20:       erc20,
		Name:        erc20Name,
		Symbol:      symbol,
		Decimals:    uint8(decimal.Uint64()),
		EventNonce:  nonce.Uint64(),
		BlockHeight: blockHeight,
	}, nil
}

func Erc20DeployedEventFromLogs(log []*web3.Log) ([]Erc20DeployedEvent, error) {
	var res []Erc20DeployedEvent
	for _, x := range log {
		event, err := Erc20DeployedEventFromLog(x)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}
	return res, nil
}

func Erc20DeployedEventFilterByEventNonce(lastEventNonce uint64, erc20Deploys []Erc20DeployedEvent) []Erc20DeployedEvent {
	var res []Erc20DeployedEvent
	for _, x := range erc20Deploys {
		if x.EventNonce > lastEventNonce {
			res = append(res, x)
		}
	}
	return res
}

type LogicCallExecutedEvent struct {
	InvalidationId []byte
	InvalidationNonce uint64
	ReturnData []byte
	EventNonce uint64
	BlockHeight uint64
}

//unimplemented!
func LogicCallExecutedEventFromLog(log *web3.Log) (LogicCallExecutedEvent, error) {
	//unimplemented!
	blockHeight := log.BlockNumber
	if blockHeight == 0 {
		return LogicCallExecutedEvent{}, errors.New("Log does not have block number, we only search logs already in blocks?")
	}
	return LogicCallExecutedEvent{
		InvalidationId:    nil,
		InvalidationNonce: 0,
		ReturnData:        nil,
		EventNonce:        0,
		BlockHeight:       0,
	}, nil
}

func LogicCallExecutedEventFromLogs(log []*web3.Log) ([]LogicCallExecutedEvent, error) {
	var res []LogicCallExecutedEvent
	for _, x := range log {
		event, err := LogicCallExecutedEventFromLog(x)
		if err != nil {
			return nil, err
		}
		res = append(res, event)
	}
	return res, nil
}

func LogicCallExecutedEventFilterByEventNonce(lastEventNonce uint64, logicCalls []LogicCallExecutedEvent) []LogicCallExecutedEvent {
	var res []LogicCallExecutedEvent
	for _, x := range logicCalls {
		if x.EventNonce > lastEventNonce {
			res = append(res, x)
		}
	}
	return res
}