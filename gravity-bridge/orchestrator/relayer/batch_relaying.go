package relayer

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/ethereum_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

func relayBatches(
	currentValSet types.ValSet,
	ethKey *ecdsa.PrivateKey,
	client *jsonrpc.Client,
	contact cosmos_gravity.Contact,
	contractAddr, gravityId string,
	timeout time.Duration,
) {
	lg := logger.GetLogger()
	ourEthereumAddress := crypto.PubkeyToAddress(ethKey.PublicKey)
	getLatestBatches := gravity_utils.Exec(func() interface{} {
		batches, err := cosmos_gravity.GetLatestTransactionBatches(contact)
		if err != nil {
			return err
		}
		return batches
	}).Await()

	latestBatches, ok := getLatestBatches.([]types.TransactionBatch)
	if !ok {
		return
	}

	for _, batch := range latestBatches {
		getSigs := gravity_utils.Exec(func() interface{} {
			sigs, err := cosmos_gravity.GetTransactionBatchSignatures(contact, batch.Nonce, batch.TokenContract)
			if err != nil {
				return err
			}
			return sigs
		}).Await()

		sigs, ok := getSigs.([]types.BatchConfirmResponse)
		if !ok {
			lg.Error(fmt.Sprintf("could not get signatures for %v:%d with %v", batch.TokenContract, batch.Nonce, sigs))
			return
		}

		var confirm []types.Confirm
		for _, v := range sigs {
			confirm = append(confirm, v)
		}

		hash := types.EncodeTxBatchConfirmHashed(gravityId, batch)
		if _, err := currentValSet.OrderSigs(hash, confirm); err != nil{
			lg.Error("Batch can not be submitted yet, waiting for more signatures")
			return
		}

		oldestSignedBatch := batch
		oldestSignatures := sigs
		erc20Contract := oldestSignedBatch.TokenContract

		getLatestEthereumBatch := gravity_utils.Exec(func() interface{} {
			latestEthereumBatch, err := ethereum_gravity.GetTxBatchNonce(contractAddr, erc20Contract, ourEthereumAddress, client)
			if err != nil {
				return err
			}
			return latestEthereumBatch
		}).Await()

		latestEthereumBatch, ok := getLatestEthereumBatch.(uint64)
		if !ok {
			lg.Error("Failed to get latest Ethereum batch")
			return
		}

		latestCosmosBatchNonce := oldestSignedBatch.Nonce
		if latestCosmosBatchNonce > latestEthereumBatch {
			res := gravity_utils.Exec(func() interface{} {
				err := ethereum_gravity.SendEthTransactionBatch(
					currentValSet,
					oldestSignedBatch,
					oldestSignatures,
					client,
					timeout,
					contractAddr,
					gravityId,
					ethKey)
				if err != nil {
					return err
				}
				return nil
			}).Await()
			if _, ok := res.(error); ok {
				lg.Error("Batch submission failed")
			}
		}
	}
}