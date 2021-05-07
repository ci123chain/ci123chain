package relayer

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

const LOOP_SPEED = 17 * time.Second

func Relayer_main_loop(logger log.Logger, ethPrivKey *ecdsa.PrivateKey, contact cosmos_gravity.Contact, contractAddr string, client *jsonrpc.Client) {
	for {
		loopStart := time.Now()
		ourEthereumAddress := crypto.PubkeyToAddress(ethPrivKey.PublicKey)
		getValSet := gravity_utils.Exec(func() interface{} {
			currentValSet, err := findLatestValset(logger, contact, contractAddr, client, ourEthereumAddress)
			if err != nil {
				return err
			}
			return currentValSet
		}).Await()

		currentValset, ok := getValSet.(types.ValSet)
		if !ok {
			logger.Error("Could not get current valset! ", getValSet.(error).Error())
			continue
		}

		getGravityId := gravity_utils.Exec(func() interface{} {
			gravityIdBz, err := ethereum_gravity.GetGravityId(contractAddr, ourEthereumAddress, client)
			if err != nil {
				return err
			}
			return gravityIdBz
		}).Await()

		gravityIdBz, ok := getGravityId.([]byte)
		if !ok {
			logger.Error("Failed to get GravityID, check your Eth node")
			return
		}

		gravityId := string(gravityIdBz)

		//relayValsets
		gravity_utils.Exec(func() interface{} {
			relayValsets(logger, currentValset, ethPrivKey, client, contact, contractAddr, gravityId, LOOP_SPEED)
			return nil
		}).Await()

		//relayBatches
		gravity_utils.Exec(func() interface{} {
			relayBatches(logger, currentValset, ethPrivKey, client, contact, contractAddr, gravityId, LOOP_SPEED)
			return nil
		}).Await()

		//relayLogicCalls
		gravity_utils.Exec(func() interface{} {
			relayLogicCalls(logger, currentValset, ethPrivKey, client, contact, contractAddr, gravityId, LOOP_SPEED)
			return nil
		}).Await()

		elapsed := time.Since(loopStart)
		if elapsed < LOOP_SPEED {
			gravity_utils.Exec(func() interface{} {
				time.Sleep(LOOP_SPEED - elapsed)
				return nil
			}).Await()
		}
	}
}
