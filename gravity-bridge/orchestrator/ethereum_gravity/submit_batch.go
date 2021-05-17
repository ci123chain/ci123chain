package ethereum_gravity

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"math/big"
	"strings"
	"time"
)

func SendEthTransactionBatch(
	currentValset types.ValSet,
	batch types.TransactionBatch,
	confirms []types.BatchConfirmResponse,
	client *jsonrpc.Client,
	timeout time.Duration,
	gravityContractAddr string,
	gravityId string,
	EthKey *ecdsa.PrivateKey,
	) error {
	newBatchNonce := batch.Nonce
	ethAddress := crypto.PubkeyToAddress(EthKey.PublicKey)

	lg := logger.GetLogger()
	lg.Info(fmt.Sprintf("Ordering signatures and submitting TransactionBatch %s:%d to Ethereum",
		batch.TokenContract.String(), newBatchNonce))

	beforeNonce, err := GetTxBatchNonce(gravityContractAddr, batch.TokenContract, ethAddress, client)
	if err != nil {
		panic(err.Error())
	}

	currentBlockHeight, err := client.Eth().BlockNumber()
	if err != nil {
		panic(err.Error())
	}

	if beforeNonce >= newBatchNonce {
		lg.Info(fmt.Sprintf("Someone else updated the batch to %d, exiting early", beforeNonce))
		return nil
	} else if currentBlockHeight > batch.BatchTimeout {
		lg.Info(fmt.Sprintf("This batch is timed out. timeout block: %d current block: %d, exiting early",
			currentBlockHeight, batch.BatchTimeout))
		return nil
	}

	payload, err := encodeBatchPayload(currentValset, &batch, confirms, gravityId)
	if err != nil {
		panic(err.Error())
	}

	sendTx := gravity_utils.Exec(func() interface{} {
		hash ,err := SendTransaction(client, gravityContractAddr, payload, big.NewInt(0), ethAddress.String(), EthKey)
		if err != nil {
			return err
		}
		return hash
	}).Await()

	hash, ok := sendTx.(*web3.Hash)
	if !ok {
		return sendTx.(error)
	}

	gravity_utils.Exec(func() interface{} {
		tx, err := WaitForTransaction(client, hash, timeout, big.NewInt(0))
		if err != nil {
			return err
		}
		return tx
	}).Await()

	lastNonce, err := GetTxBatchNonce(gravityContractAddr, batch.TokenContract, ethAddress, client)
	if err != nil {
		return err
	}

	if lastNonce != newBatchNonce {
		lg.Error(fmt.Sprintf("Current nonce is %d expected to update to nonce %d", lastNonce, newBatchNonce))
	} else {
		lg.Info(fmt.Sprintf("Successfully updated Batch with new Nonce %d", lastNonce))
	}

	return nil
}

func encodeBatchPayload(
	currentValset types.ValSet,
	batch *types.TransactionBatch,
	confirms []types.BatchConfirmResponse,
	gravityId string,
	) ([]byte, error) {
	currentAddresses, currentPowers := currentValset.FilterEmptyAddress()
	currentValsetNonce := currentValset.Nonce
	newBatchNonce := batch.Nonce
	hash := types.EncodeTxBatchConfirmHashed(gravityId, *batch)

	var confirm []types.Confirm
	for _, v := range confirms {
		confirm = append(confirm, v)
	}

	sigData, err := currentValset.OrderSigs(hash, confirm)
	if err != nil {
		return nil, err
	}
	sigArrays := types.ToArrays(sigData)
	amounts, fees, destinations := batch.GetCheckPointValues()
	sig, err := utils.ParseSignature(strings.Replace("submitBatch(address[],uint256[],uint256,uint8[],bytes32[],bytes32[],uint256[],address[],uint256[],uint256,address,uint256)", " ", "", -1))
	if err != nil {
		return nil, err
	}
	data := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{currentAddresses, currentPowers, currentValsetNonce, sigArrays.V, sigArrays.R, sigArrays.S, amounts, destinations, fees, newBatchNonce, batch.TokenContract.String(), batch.BatchTimeout})...)
	return data, nil
}