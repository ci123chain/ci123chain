package main

import (
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

func getBlockNumberWithRetry(client *jsonrpc.Client) uint64 {
	for {
		getBlockNumber := gravity_utils.Exec(func() interface{} {
			blockNumber, err := client.Eth().BlockNumber()
			if err != nil {
				return err
			}
			return blockNumber
		}).Await()

		lg := logger.GetLogger()
		blockNumber, ok := getBlockNumber.(uint64)
		if ok {
			return blockNumber
		} else {
			lg.Info("Retry to get eth_block_number")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(RETRY_TIME)
				return nil
			}).Await()
		}
	}
}

func getLastEventNonceWithRetry(ourCosmosAddress common.Address, contact cosmos_gravity.Contact) uint64 {
	for {
		getCosmosLatestEventNonce := gravity_utils.Exec(func() interface{} {
			nonce, err := cosmos_gravity.GetLastEventNonce(ourCosmosAddress, contact)
			if err != nil {
				return err
			}
			return nonce
		}).Await()

		lg := logger.GetLogger()
		nonce, ok := getCosmosLatestEventNonce.(uint64)
		if ok {
			return nonce
		} else {
			lg.Info("Retry to get cosmos_latest_event_nonce")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(RETRY_TIME)
				return nil
			}).Await()
		}
	}
}

func getNetVersionWithRetry(client *jsonrpc.Client) uint64 {
	for {
		getNetVersion := gravity_utils.Exec(func() interface{} {
			netVersion, err := client.Net().Version()
			if err != nil {
				return err
			}
			return netVersion
		}).Await()

		lg := logger.GetLogger()
		version, ok := getNetVersion.(uint64)
		if ok {
			return version
		} else {
			lg.Info("Retry to get eth net_version")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(RETRY_TIME)
				return nil
			}).Await()
		}
	}
}


