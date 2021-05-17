package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/relayer"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

const (
	feeAmount = 1
	ETH_ORACLE_LOOP_SPEED = 13 * time.Second
	ETH_SIGNER_LOOP_SPEED = 11 * time.Second
	DELAY = 5 * time.Second
)

func orchestratorMainLoop(cosmosKey, ethKey, cosmosRpc, ethRpc, denom, contractAddr string) {
	var fee = sdk.Coin{
		Denom:  denom,
		Amount: sdk.NewInt(feeAmount),
	}

	cosmosPrivKey, err := crypto.HexToECDSA(cosmosKey)
	if err != nil {
		panic(err)
	}

	ethPrivKey, err := crypto.HexToECDSA(ethKey)
	if err != nil {
		panic(err)
	}

	client, err := jsonrpc.NewClient(ethRpc)
	if err != nil {
		panic(err)
	}

	contact := cosmos_gravity.NewContact(cosmosRpc)

	go gravity_utils.Exec(func() interface{} {
		eth_oracle_main_loop(cosmosPrivKey, contact, contractAddr, client, fee)
		return nil
	}).Await()

	go gravity_utils.Exec(func() interface{} {
		eth_signer_main_loop(ethPrivKey, cosmosPrivKey, contact, contractAddr, client, fee)
		return nil
	}).Await()

	go gravity_utils.Exec(func() interface{} {
		relayer.Relayer_main_loop(ethPrivKey, contact, contractAddr, client)
		return nil
	}).Await()
}

func eth_oracle_main_loop(cosmosPrivKey *ecdsa.PrivateKey, contact cosmos_gravity.Contact, contractAddr string, client *jsonrpc.Client, fee sdk.Coin) {
	ourCosmosAddress := crypto.PubkeyToAddress(cosmosPrivKey.PublicKey)

	latestCheckedBlock := gravity_utils.Exec(func() interface{} {
		return getLastCheckedBlock(contact, contractAddr, ourCosmosAddress, client)
	}).Await().(uint64)

	lg := logger.GetLogger()
	lg.Info("Oracle resync complete, Oracle now operational")

	for {
		loopStart := time.Now()
		getLatestEthBlock := gravity_utils.Exec(func() interface{} {
			blockNumber, err := client.Eth().BlockNumber()
			if err != nil {
				return err
			}
			return blockNumber
		}).Await()

		getLatestCosmosBlock := gravity_utils.Exec(func() interface{} {
			chainStatus, err := contact.GetChainStatus()
			if err != nil {
				return err
			}
			return chainStatus
		}).Await()

		latestEthBlock, ok := getLatestEthBlock.(uint64)
		if !ok {
			lg.Error("Could not reach Ethereum rpc!")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(DELAY)
				return nil
			}).Await()
			continue
		}

		latestCosmosBlock, ok := getLatestCosmosBlock.(cosmos_gravity.ChainStatus)
		if !ok {
			lg.Error("Could not reach Cosmos rpc!")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(DELAY)
				return nil
			}).Await()
			continue
		}

		switch latestCosmosBlock.Status {
			case cosmos_gravity.MOVING:
				lg.Info(fmt.Sprintf("Latest Eth block: %d, Latest Cosmos block: %d", latestEthBlock, latestCosmosBlock.BlockHeight))
			case cosmos_gravity.SYNCING:
				lg.Error(fmt.Sprintf("Cosmos node syncing, Eth signer paused"))
				gravity_utils.Exec(func() interface{} {
					time.Sleep(DELAY)
					return nil
				}).Await()
				continue
			case cosmos_gravity.WAITING_TO_START:
				lg.Error(fmt.Sprintf("Cosmos node syncing waiting for chain start, Eth signer paused"))
				gravity_utils.Exec(func() interface{} {
					time.Sleep(DELAY)
					return nil
				}).Await()
				continue
		}

		getLatestCheckedBlock := gravity_utils.Exec(func() interface{} {
			newBlock, err := CheckForEvents(client, contact, contractAddr, cosmosPrivKey, fee, latestCheckedBlock)
			if err != nil {
				return err
			}
			return newBlock
		}).Await()

		_, ok = getLatestCheckedBlock.(uint64)
		if !ok {
			lg.Error(fmt.Sprintf("Failed to get events for block range, Check your Eth node and Cosmos RPC, error: %s", getLatestCheckedBlock.(error).Error()))
		} else {
			latestCheckedBlock = getLatestCheckedBlock.(uint64)
		}

		elapsed := time.Since(loopStart)
		if elapsed < ETH_ORACLE_LOOP_SPEED {
			gravity_utils.Exec(func() interface{} {
				time.Sleep(ETH_ORACLE_LOOP_SPEED - elapsed)
				return nil
			}).Await()
		}
	}
}

func eth_signer_main_loop(ethPrivKey, cosmosPrivKey *ecdsa.PrivateKey, contact cosmos_gravity.Contact, contractAddr string, client *jsonrpc.Client, fee sdk.Coin) {
	ourCosmosAddress := crypto.PubkeyToAddress(cosmosPrivKey.PublicKey)
	ourEthereumAddress := crypto.PubkeyToAddress(ethPrivKey.PublicKey)
	getGravityId := gravity_utils.Exec(func() interface{} {
		gravityIdBz, err := ethereum_gravity.GetGravityId(contractAddr, ourEthereumAddress, client)
		if err != nil {
			return err
		}
		return gravityIdBz
	}).Await()

	lg := logger.GetLogger()

	gravityIdBz, ok := getGravityId.([]byte)
	if !ok {
		lg.Error("Failed to get GravityID, check your Eth node")
		return
	}

	gravityId := string(gravityIdBz)

	for {
		loopStart := time.Now()
		getLatestEthBlock := gravity_utils.Exec(func() interface{} {
			blockNumber, err := client.Eth().BlockNumber()
			if err != nil {
				return err
			}
			return blockNumber
		}).Await()

		getLatestCosmosBlock := gravity_utils.Exec(func() interface{} {
			chainStatus, err := contact.GetChainStatus()
			if err != nil {
				return err
			}
			return chainStatus
		}).Await()

		latestEthBlock, ok := getLatestEthBlock.(uint64)
		if !ok {
			lg.Error("Could not reach Ethereum rpc!")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(DELAY)
				return nil
			}).Await()
			continue
		}

		latestCosmosBlock, ok := getLatestCosmosBlock.(cosmos_gravity.ChainStatus)
		if !ok {
			lg.Error("Could not reach Cosmos rpc!")
			gravity_utils.Exec(func() interface{} {
				time.Sleep(DELAY)
				return nil
			}).Await()
			continue
		}

		switch latestCosmosBlock.Status {
		case cosmos_gravity.MOVING:
			lg.Info(fmt.Sprintf("Latest Eth block: %d, Latest Cosmos block: %d", latestEthBlock, latestCosmosBlock.BlockHeight))
		case cosmos_gravity.SYNCING:
			lg.Error(fmt.Sprintf("Cosmos node syncing, Eth signer paused"))
			gravity_utils.Exec(func() interface{} {
				time.Sleep(DELAY)
				return nil
			}).Await()
			continue
		case cosmos_gravity.WAITING_TO_START:
			lg.Error(fmt.Sprintf("Cosmos node syncing waiting for chain start, Eth signer paused"))
			gravity_utils.Exec(func() interface{} {
				time.Sleep(DELAY)
				return nil
			}).Await()
			continue
		}

		getOldestUnsignedValsets := gravity_utils.Exec(func() interface{} {
			valsets, err := cosmos_gravity.GetOldestUnsignedValsets(contact, ourCosmosAddress)
			if err != nil {
				return err
			}
			return valsets
		}).Await()

		oldestUnsignedValsets, ok := getOldestUnsignedValsets.([]types.ValSet)
		if !ok {
			lg.Error(fmt.Sprintf("Failed to get unsigned valsets, check your Cosmos RPC, error: %s", getOldestUnsignedValsets.(error).Error()))
		}

		if len(oldestUnsignedValsets) == 0 {
			lg.Info("No validator sets to sign, node is caught up!")
		} else {
			lg.Info(fmt.Sprintf("Sending %d valset confirms starting with %d", len(oldestUnsignedValsets), oldestUnsignedValsets[0].Nonce))
			res := gravity_utils.Exec(func() interface{} {
				txResponse, err := cosmos_gravity.SendValsetConfirms(contact, ethPrivKey, fee, oldestUnsignedValsets, cosmosPrivKey, gravityId)
				if err != nil {
					return err
				}
				return txResponse
			}).Await()
			lg.Info(fmt.Sprintf("Valset confirm result is %v", res.(sdk.TxResponse)))
		}

		elapsed := time.Since(loopStart)
		if elapsed < ETH_SIGNER_LOOP_SPEED {
			gravity_utils.Exec(func() interface{} {
				time.Sleep(ETH_SIGNER_LOOP_SPEED - elapsed)
				return nil
			}).Await()
		}
	}
}