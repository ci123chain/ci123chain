package main

import (
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/umbracle/go-web3/jsonrpc"
)

func getBlockNumberWithRetry(logger log.Logger, client *jsonrpc.Client) uint64 {
	for {
		getBlockNumber := gravity_utils.Exec(func() interface{} {
			blockNumber, err := client.Eth().BlockNumber()
			if err != nil {
				return err
			}
			return blockNumber
		}).Await()

		blockNumber, ok := getBlockNumber.(uint64)
		if ok {
			return blockNumber
		} else {
			logger.Info("Retry to get eth_block_number")
		}
	}
}

func getLastEventNonceWithRetry(logger log.Logger, ourCosmosAddress common.Address, contact cosmos_gravity.Contact) uint64 {
	for {
		getCosmosLatestEventNonce := gravity_utils.Exec(func() interface{} {
			nonce, err := getLastEventNonce(logger, ourCosmosAddress, contact)
			if err != nil {
				return err
			}
			return nonce
		}).Await()

		nonce, ok := getCosmosLatestEventNonce.(uint64)
		if ok {
			return nonce
		} else {
			logger.Info("Retry to get cosmos_latest_event_nonce")
		}
	}
}

func getLastEventNonce(logger log.Logger, ourCosmosAddress common.Address, contact cosmos_gravity.Contact) (uint64, error) {
	//QueryLastEventNonceByAddrRequest
	return 0, nil
}
