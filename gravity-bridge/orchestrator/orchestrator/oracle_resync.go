package main

import (
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

const (
	BLOCKS_TO_SEARCH = 5000
	RETRY_TIME = 5 * time.Second
)

// This function retrieves the last event nonce this oracle has relayed to Cosmos
// it then uses the Ethereum indexes to determine what block the last entry
func getLastCheckedBlock(contact cosmos_gravity.Contact, contractAddr string, ourCosmosAddress common.Address, client *jsonrpc.Client) uint64 {
	lastBlock := gravity_utils.Exec(func() interface{} {
		return getBlockNumberWithRetry(client)
	}).Await().(uint64)

	lastEventNonce := gravity_utils.Exec(func() interface{} {
		return getLastEventNonceWithRetry(ourCosmosAddress, contact)
	}).Await().(uint64)

	// zero indicates this oracle has never submitted an event before since there is no
	// zero event nonce (it's pre-incremented in the solidity contract) we have to go
	// and look for event nonce one.
	if lastEventNonce == 0 {
		lastEventNonce = 1
	}

	lg := logger.GetLogger()

	currentBlock := lastBlock
	for {
		if currentBlock == 0 {
			break
		}
		lg.Info(fmt.Sprintf("Oracle is resyncing, looking back into the history to find our last event nonce: %d, on block: %d", lastEventNonce, currentBlock))

		var endSearch uint64
		if currentBlock < BLOCKS_TO_SEARCH {
			endSearch = 0
		} else {
			endSearch = currentBlock - BLOCKS_TO_SEARCH
		}

		getBatchEvents := gravity_utils.Exec(func() interface{} {
			batchEvents, err := ethereum_gravity.CheckForEvents(endSearch,
				currentBlock,
				[]string{contractAddr},
				[]string{"TransactionBatchExecutedEvent(uint256,address,uint256)"},
				client)
			if err != nil {
				return err
			}
			return batchEvents
		}).Await()

		batchEvents, okb := getBatchEvents.([]*web3.Log)

		getSendToCosmosEvents := gravity_utils.Exec(func() interface{} {
			sendToCosmosEvents, err := ethereum_gravity.CheckForEvents(endSearch,
				currentBlock,
				[]string{contractAddr},
				[]string{"SendToCosmosEvent(address,address,bytes32,uint256,uint256)"},
				client)
			if err != nil {
				return err
			}
			return sendToCosmosEvents
		}).Await()

		sendToCosmosEvents, oks := getSendToCosmosEvents.([]*web3.Log)

		getErc20DeployedEvents := gravity_utils.Exec(func() interface{} {
			erc20DeployedEvents, err := ethereum_gravity.CheckForEvents(endSearch,
				currentBlock,
				[]string{contractAddr},
				[]string{"ERC20DeployedEvent(string,address,string,string,uint8,uint256)"},
				client)
			if err != nil {
				return err
			}
			return erc20DeployedEvents
		}).Await()

		erc20DeployedEvents, oke := getErc20DeployedEvents.([]*web3.Log)

		getLogicCallExecutedEvents := gravity_utils.Exec(func() interface{} {
			logicCallExecutedEvents, err := ethereum_gravity.CheckForEvents(endSearch,
				currentBlock,
				[]string{contractAddr},
				[]string{"LogicCallEvent(bytes32,uint256,bytes,uint256)"},
				client)
			if err != nil {
				return err
			}
			return logicCallExecutedEvents
		}).Await()

		logicCallExecutedEvents, okl := getLogicCallExecutedEvents.([]*web3.Log)

		getValSetEvents := gravity_utils.Exec(func() interface{} {
			valSetEvents, err := ethereum_gravity.CheckForEvents(endSearch,
				currentBlock,
				[]string{contractAddr},
				[]string{"ValsetUpdatedEvent(uint256,address[],uint256[])"},
				client)
			if err != nil {
				return err
			}
			return valSetEvents
		}).Await()

		valSetEvents, okv := getValSetEvents.([]*web3.Log)

		if !okb || !oks || !oke || !okl || !okv {
			lg.Error("Failed to get blockchain events while resyncing, is your Eth node working? If you see only one of these it's fine")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(RETRY_TIME)
				return nil
			}).Await()
			continue
		}

		// look for and return the block number of the event last seen on the Cosmos chain
		// then we will play events from that block (including that block, just in case
		// there is more than one event there) onwards. We use valset nonce 0 as an indicator
		// of what block the contract was deployed on.
		for _, event := range batchEvents {
			transactionBatchExecutedEvent, err := types.TransactionBatchExecutedEventFromLog(event)
			if err != nil {
				lg.Error(fmt.Sprintf("Got Batch event that we can't parse: %s", err.Error()))
			}
			if transactionBatchExecutedEvent.EventNonce == lastEventNonce && event.BlockNumber != 0 {
				return event.BlockNumber
			}
		}

		for _, event := range sendToCosmosEvents {
			sendToCosmosEvent, err := types.SendToCosmosEventFromLog(event)
			if err != nil {
				lg.Error(fmt.Sprintf("Got SendToCosmos event that we can't parse: %s", err.Error()))
			}
			if sendToCosmosEvent.EventNonce == lastEventNonce && event.BlockNumber != 0 {
				return event.BlockNumber
			}
		}

		for _, event := range erc20DeployedEvents {
			erc20DeployedEvent, err := types.Erc20DeployedEventFromLog(event)
			if err != nil {
				lg.Error(fmt.Sprintf("Got Erc20Deployed event that we can't parse: %s", err.Error()))
			}
			if erc20DeployedEvent.EventNonce == lastEventNonce && event.BlockNumber != 0 {
				return event.BlockNumber
			}
		}

		for _, event := range logicCallExecutedEvents {
			logicCallExecutedEvent, err := types.LogicCallExecutedEventFromLog(event)
			if err != nil {
				lg.Error(fmt.Sprintf("Got LogicCallExecuted event that we can't parse: %s", err.Error()))
			}
			if logicCallExecutedEvent.EventNonce == lastEventNonce && event.BlockNumber != 0 {
				return event.BlockNumber
			}
		}

		for _, event := range valSetEvents {
			valSetEvent, err := types.ValSetUpdatedEventFromLog(event)
			if err != nil {
				lg.Error(fmt.Sprintf("Got ValsetUpdate event that we can't parse: %s", err.Error()))
			}
			// if we've found this event it is the first possible event from the contract
			// no other events can come before it, therefore either there's been a parsing error
			// or no events have been submitted on this chain yet.
			if valSetEvent.Nonce == 0 && lastEventNonce == 1 {
				return lastBlock
			}
			// if we're looking for a later event nonce and we find the deployment of the contract
			// we must have failed to parse the event we're looking for. The oracle can not start
			if valSetEvent.Nonce == 0 && lastEventNonce > 1 {
				lg.Error(fmt.Sprintf("Could not find the last event relayed by {}, Last Event nonce is {} but no event matching that could be found!", ourCosmosAddress, lastEventNonce))
			}
		}
		currentBlock = endSearch
	}

	panic("You have reached the end of block history without finding the Gravity contract deploy event! You must have the wrong contract address!")
}
