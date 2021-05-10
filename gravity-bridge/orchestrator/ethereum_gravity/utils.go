package ethereum_gravity

import (
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"math/big"
	"strconv"
)

func GetGravityId(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) ([]byte, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	data := Digest("state_gravityId()")

	val, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     data,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)

	return []byte(val), err
}

func GetValSetNonce(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) (uint64, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	data := Digest("state_lastValsetNonce()")

	res, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     data,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)
	if err != nil {
		return 0, err
	}
	nonce, err := strconv.ParseUint(res, 10, 64)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func GetEventNonce(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) (uint64, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	data := Digest("state_lastEventNonce()")

	res, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     data,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)
	if err != nil {
		return 0, err
	}
	nonce, err := strconv.ParseUint(res, 10, 64)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func CheckForEvents(startBlock, endBlock uint64, contractAddr []string, events []string, client *jsonrpc.Client) ([]*web3.Log, error) {
	var finalTopics []*web3.Hash
	var addresses []web3.Address

	fromBlock := web3.BlockNumber(startBlock)
	toBlock := web3.BlockNumber(endBlock)

	for _, contract := range contractAddr {
		addresses = append(addresses, web3.HexToAddress(contract))
	}

	for _, event := range events {
		sig := Digest(event)
		var eventHash *web3.Hash
		eventHash.UnmarshalText(sig)
		finalTopics = append(finalTopics, eventHash)
	}

	filter := &web3.LogFilter{
		Address:   addresses,
		Topics:    finalTopics,
		BlockHash: nil,
		From:      &fromBlock,
		To:        &toBlock,
	}

	res := gravity_utils.Exec(func() interface{} {
		logs, err := client.Eth().GetLogs(filter)
		if err != nil {
			return err
		}
		return logs
	}).Await()

	logs, ok := res.([]*web3.Log)
	if !ok {
		return nil, res.(error)
	}

	return logs, nil
}
