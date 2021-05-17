package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
)

func CheckForEvents(client *jsonrpc.Client,
	contact cosmos_gravity.Contact,
	contractAddr string,
	cosmosPrivKey *ecdsa.PrivateKey,
	fee sdk.Coin,
	startingBlock uint64) (uint64, error) {

	lg := logger.GetLogger()
	ourCosmosAddress := crypto.PubkeyToAddress(cosmosPrivKey.PublicKey)
	latestBlock := gravity_utils.Exec(func() interface{} {
		return getBlockNumberWithRetry(client)
	}).Await().(uint64)

	blockDelay := gravity_utils.Exec(func() interface{} {
		return getBlockDelay(client)
	}).Await().(uint64)

	latestBlock -= blockDelay

	getDeposits := gravity_utils.Exec(func() interface{} {
		events, err := ethereum_gravity.CheckForEvents(
			startingBlock,
			latestBlock,
			[]string{contractAddr},
			[]string{"SendToCosmosEvent(address,address,bytes32,uint256,uint256)"},
			client)
		if err != nil {
			return err
		}
		return events
	}).Await()

	lg.Info(fmt.Sprintf("Deposits %v", getDeposits))
	depositsEvents, okd := getDeposits.([]*web3.Log)

	getBatches := gravity_utils.Exec(func() interface{} {
		events, err := ethereum_gravity.CheckForEvents(
			startingBlock,
			latestBlock,
			[]string{contractAddr},
			[]string{"TransactionBatchExecutedEvent(uint256,address,uint256)"},
			client)
		if err != nil {
			return err
		}
		return events
	}).Await()

	lg.Info(fmt.Sprintf("Batches %v", getBatches))
	batchEvents, okb := getBatches.([]*web3.Log)

	getValsets := gravity_utils.Exec(func() interface{} {
		events, err := ethereum_gravity.CheckForEvents(
			startingBlock,
			latestBlock,
			[]string{contractAddr},
			[]string{"ValsetUpdatedEvent(uint256,address[],uint256[])"},
			client)
		if err != nil {
			return err
		}
		return events
	}).Await()

	lg.Info(fmt.Sprintf("Valsets %v", getValsets))
	valsetsEvents, okv := getValsets.([]*web3.Log)

	getErc20Deployed := gravity_utils.Exec(func() interface{} {
		events, err := ethereum_gravity.CheckForEvents(
			startingBlock,
			latestBlock,
			[]string{contractAddr},
			[]string{"ERC20DeployedEvent(string,address,string,string,uint8,uint256)"},
			client)
		if err != nil {
			return err
		}
		return events
	}).Await()

	lg.Info(fmt.Sprintf("Erc20 Deployments %v", getErc20Deployed))
	erc20DeployedEvents, oke := getErc20Deployed.([]*web3.Log)

	getLogicCallExecuted := gravity_utils.Exec(func() interface{} {
		events, err := ethereum_gravity.CheckForEvents(
			startingBlock,
			latestBlock,
			[]string{contractAddr},
			[]string{"LogicCallEvent(bytes32,uint256,bytes,uint256)"},
			client)
		if err != nil {
			return err
		}
		return events
	}).Await()

	lg.Info(fmt.Sprintf("Logic call executions %v", getLogicCallExecuted))
	logicCallExecutedEvents, okl := getLogicCallExecuted.([]*web3.Log)

	if okd && okb && okv && oke && okl {
		valsets, _ := types.ValSetUpdatedEventFromLogs(valsetsEvents)
		lg.Info(fmt.Sprintf("Parsed valsets: %v", valsets))

		withdraws, _ := types.TransactionBatchExecutedEventFromLogs(batchEvents)
		lg.Info(fmt.Sprintf("Parsed batches: %v", withdraws))

		deposits, _ := types.SendToCosmosEventFromLogs(depositsEvents)
		lg.Info(fmt.Sprintf("Parsed deposits: %v", deposits))

		erc20Deploys, _ := types.Erc20DeployedEventFromLogs(erc20DeployedEvents)
		lg.Info(fmt.Sprintf("Parsed erc20 deploys: %v", erc20Deploys))

		logicCalls, _ := types.LogicCallExecutedEventFromLogs(logicCallExecutedEvents)
		lg.Info(fmt.Sprintf("Logic call executions: %v", logicCalls))

		// note that starting block overlaps with our last checked block, because we have to deal with
		// the possibility that the relayer was killed after relaying only one of multiple events in a single
		// block, so we also need this routine so make sure we don't send in the first event in this hypothetical
		// multi event block again. In theory we only send all events for every block and that will pass of fail
		// atomicly but lets not take that risk.
		lastEventNonce := gravity_utils.Exec(func() interface{} {
			nonce, err := cosmos_gravity.GetLastEventNonce(contact, ourCosmosAddress)
			if err != nil {
				return err
			}
			return nonce
		}).Await().(uint64)

		deposits = types.SendToCosmosEventFilterByEventNonce(lastEventNonce, deposits)
		withdraws = types.TransactionBatchExecutedEventFilterByEventNonce(lastEventNonce, withdraws)
		erc20Deploys = types.Erc20DeployedEventFilterByEventNonce(lastEventNonce, erc20Deploys)
		logicCalls = types.LogicCallExecutedEventFilterByEventNonce(lastEventNonce, logicCalls)

		if len(deposits) != 0 {
			lg.Info(fmt.Sprintf("Oracle observed deposit with sender %v, destination %v, amount %v, and event nonce %d",
				deposits[0].Sender, deposits[0].Destination, deposits[0].Amount, deposits[0].EventNonce))
		}

		if len(withdraws) != 0 {
			lg.Info(fmt.Sprintf("Oracle observed batch with nonce %d, contract %v, and event nonce %d",
				withdraws[0].BatchNonce, withdraws[0].Erc20, withdraws[0].EventNonce))
		}

		if len(erc20Deploys) != 0 {
			lg.Info(fmt.Sprintf("Oracle observed ERC20 deployment with denom %s erc20 name %s and symbol %s and event nonce %d",
				erc20Deploys[0].CosmosDenom, erc20Deploys[0].Name, erc20Deploys[0].Symbol, erc20Deploys[0].EventNonce,
			))
		}

		if len(logicCalls) != 0 {
			lg.Info(fmt.Sprintf("Oracle observed logic call execution with ID %s Nonce %d and event nonce %d",
				hex.EncodeToString(logicCalls[0].InvalidationId), logicCalls[0].InvalidationNonce, logicCalls[0].EventNonce))
		}

		if len(deposits) != 0 || len(withdraws) != 0 || len(erc20Deploys) != 0 || len(logicCalls) != 0{
			res := gravity_utils.Exec(func() interface{} {
				txResponse, err := cosmos_gravity.SendEthereumClaims(contact, cosmosPrivKey, deposits, withdraws, erc20Deploys, logicCalls, fee)
				if err != nil {
					return err
				}
				return txResponse
			}).Await().(sdk.TxResponse)

			lg.Info(fmt.Sprintf("Claims response: %v", res))

			newEventNonce := gravity_utils.Exec(func() interface{} {
				nonce, err := cosmos_gravity.GetLastEventNonce(contact, ourCosmosAddress)
				if err != nil {
					return err
				}
				return nonce
			}).Await().(uint64)

			if newEventNonce == lastEventNonce {
				return 0, errors.New(fmt.Sprintf("Claims did not process, trying to update but still on %d, trying again in a moment, check txhash %s for errors", lastEventNonce, res.TxHash))
			} else {
				lg.Info(fmt.Sprintf("Claims processed, new nonce %d", newEventNonce))
			}
		}
		return latestBlock, nil
	} else {
		lg.Error("Failed to get events")
		return 0, errors.New("Failed to get logs!")
	}
}

func getBlockDelay(client *jsonrpc.Client) uint64 {
	netVersion := gravity_utils.Exec(func() interface{} {
		return 	getNetVersionWithRetry(client)
	}).Await().(uint64)

	switch netVersion {
	// Mainline Ethereum, Ethereum classic, or the Ropsten, Mordor testnets
	// all POW Chains
	case 1, 3, 7:
		return 6
	// Rinkeby, Goerli, Dev, our own Gravity Ethereum testnet, Kotti and Cichain respectively
	// all non-pow chains
	case 4, 5, 2018, 15, 6, 444900:
		return 0
	// assume the safe option (POW) where we don't know
	default:
		return 6
	}
}


