package relayer

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

const LOOP_SPEED = 17 * time.Second

func Relayer_main_loop(ethPrivKey *ecdsa.PrivateKey, contact cosmos_gravity.Contact, contractAddr string, client *jsonrpc.Client) {
	for {
		loopStart := time.Now()
		ourEthereumAddress := crypto.PubkeyToAddress(ethPrivKey.PublicKey)
		getValSet := gravity_utils.Exec(func() interface{} {
			currentValSet, err := findLatestValset(contact, contractAddr, client, ourEthereumAddress)
			if err != nil {
				return err
			}
			return currentValSet
		}).Await()

		lg := logger.GetLogger()

		currentValset, ok := getValSet.(types.ValSet)
		if !ok {
			lg.Error("Could not get current valset! ", getValSet.(error).Error())
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
			lg.Error("Failed to get GravityID, check your Eth node")
			return
		}

		gravityId := string(gravityIdBz)

		//relayValsets
		gravity_utils.Exec(func() interface{} {
			relayValsets(currentValset, ethPrivKey, client, contact, contractAddr, gravityId, LOOP_SPEED)
			return nil
		}).Await()

		//relayBatches
		gravity_utils.Exec(func() interface{} {
			relayBatches(currentValset, ethPrivKey, client, contact, contractAddr, gravityId, LOOP_SPEED)
			return nil
		}).Await()

		//relayLogicCalls
		gravity_utils.Exec(func() interface{} {
			relayLogicCalls(currentValset, ethPrivKey, client, contact, contractAddr, gravityId, LOOP_SPEED)
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
