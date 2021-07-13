package ethereum_gravity

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"math/big"
	"strings"
	"time"
)

const (
	GasLimit uint64 = 7000000
)

func GetGravityId(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) ([]byte, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	sig, err := utils.ParseSignature(strings.Replace("state_gravityId()", " ", "", -1))
	if err != nil {
		return nil, err
	}
	hash := utils.MethodID(sig.Method, sig.Args)

	val, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     hash,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)

	return hex.DecodeString(val[2:])
}

func GetValSetNonce(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) (uint64, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	sig, err := utils.ParseSignature(strings.Replace("state_lastValsetNonce()", " ", "", -1))
	if err != nil {
		return 0, err
	}
	hash := utils.MethodID(sig.Method, sig.Args)

  	res, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     hash,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)
	if err != nil {
		return 0, err
	}
	nonce, err := hex.DecodeString(res[2:])
	if err != nil {
		return 0, err
	}
	return new(big.Int).SetBytes(nonce).Uint64(), nil
}

func GetEventNonce(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) (uint64, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	sig, err := utils.ParseSignature(strings.Replace("state_lastEventNonce()", " ", "", -1))
	if err != nil {
		return 0, err
	}
	hash := utils.MethodID(sig.Method, sig.Args)

	res, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     hash,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)
	if err != nil {
		return 0, err
	}
	nonce, err := hex.DecodeString(res[2:])
	if err != nil {
		return 0, err
	}
	return new(big.Int).SetBytes(nonce).Uint64(), nil
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
		sig, err := utils.ParseSignature(strings.Replace(event, " ", "", -1))
		if err != nil {
		return nil, err
		}

		hash := utils.EventID(sig.Method, sig.Args)
		eventHash := web3.HexToHash("0x" + hex.EncodeToString(hash))
		finalTopics = append(finalTopics, &eventHash)
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

func GetTxBatchNonce(contractAddr string, erc20Addr, ourEthereumAddress common.Address, client *jsonrpc.Client) (uint64, error) {
	contractAddress := web3.HexToAddress(contractAddr)
	//ourAddress := web3.HexToAddress(ourEthereumAddress.String())
	//ourBalance, err := client.Eth().GetBalance(ourAddress, -1)
	//if err != nil {
	//	return 0, err
	//}
	//nonce, err := client.Eth().GetNonce(ourAddress, -1)
	//if err != nil {
	//	return 0, err
	//}
	sig, err := utils.ParseSignature(strings.Replace("lastBatchNonce(address)", " ", "", -1))
	if err != nil {
		return 0, err
	}
	data := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{erc20Addr.String()})...)

	res, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ourEthereumAddress.String()),
		To:       &contractAddress,
		Data:     data,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)

	nonce, err := hex.DecodeString(res[2:])
	if err != nil {
		return 0, err
	}
	return new(big.Int).SetBytes(nonce).Uint64(), nil
}

type GasCost struct {
	Gas *big.Int
	GasPrice *big.Int
}

func (gc GasCost) GetTotal() *big.Int {
	return gc.Gas.Mul(gc.Gas, gc.GasPrice)
}

func SendTransaction(
	client *jsonrpc.Client,
	toAddress string,
	data []byte,
	value *big.Int,
	ownAddress string,
	secret *ecdsa.PrivateKey,
	) (*web3.Hash, error) {

	nonce, _ := client.Eth().GetNonce(web3.HexToAddress(ownAddress), -1)
	chainId, _ := client.Eth().ChainID()
	balance, _ := client.Eth().GetBalance(web3.HexToAddress(ownAddress), -1)
	gasPrice := big.NewInt(1)
	var gasLimit uint64
	if balance.Cmp(new(big.Int).SetUint64(GasLimit)) > 0 {
		gasLimit = GasLimit
	} else {
		gasLimit = balance.Uint64()
	}
	to := common.HexToAddress(toAddress)
	tx := types.NewMsgEthereumTx(nonce, &to, value, gasLimit, gasPrice, data)
	hash := tx.RLPSignBytes(chainId)
	sig, err := crypto.Sign(hash.Bytes(), secret)
	if err != nil {
		return nil, err
	}
	if len(sig) != 65 {
		return nil, fmt.Errorf("wrong size for signature: got %d, want 65", len(sig))
	}

	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:64])

	var v *big.Int

	if chainId.Sign() == 0 {
		v = new(big.Int).SetBytes([]byte{sig[64] + 27})

	} else {
		v = big.NewInt(int64(sig[64] + 35))
		chainIDMul := new(big.Int).Mul(chainId, big.NewInt(2))

		v.Add(v, chainIDMul)
	}

	tx.Data.V = v
	tx.Data.R = r
	tx.Data.S = s

	txBz, err := rlp.EncodeToBytes(&tx)
	if err != nil {
		return nil, err
	}

	trySendTx := gravity_utils.Exec(func() interface{} {
		hash, err := client.Eth().SendRawTransaction(txBz)
		if err != nil {
			return err
		}
		return hash
	}).Await()

	txHash, ok := trySendTx.(web3.Hash)
	if !ok {
		return nil, trySendTx.(error)
	}

	return &txHash, nil
}

func WaitForTransaction(client *jsonrpc.Client, hash *web3.Hash, timeout time.Duration, blocksToWait *big.Int) (*web3.Transaction, error) {
	loopStart := time.Now()
	response, err := client.Eth().GetTransactionReceipt(*hash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response)
	//return nil, nil
	for {
		gravity_utils.Exec(func() interface{} {
			time.Sleep(1 * time.Second)
			return nil
		}).Await()

		getTransaction := gravity_utils.Exec(func() interface{} {
			response, err := client.Eth().GetTransactionByHash(*hash)
			if err != nil {
				return err
			}
			return response
		}).Await()

		transaction, ok := getTransaction.(*web3.Transaction)
		if !ok {
			return nil, getTransaction.(error)
		}

		if transaction == nil {
			continue
		}

		if blocksToWait.Sign() == 0 && transaction.BlockNumber != 0 {
			return transaction, nil
		}

		currentBlock, _ := client.Eth().BlockNumber()
		if currentBlock > blocksToWait.Uint64() && currentBlock - blocksToWait.Uint64() >= transaction.BlockNumber {
			return transaction, nil
		}

		elapsed := time.Since(loopStart)
		if elapsed > timeout {
			return nil, errors.New("Transaction timeout")
		}
	}
}

func CheckErc20Approved(erc20, ownAddress, targetContract string, client *jsonrpc.Client) (bool, error) {
	sig, err := utils.ParseSignature(strings.Replace("allowance(address,address)", " ", "", -1))
	if err != nil {
		return false, err
	}
	data := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{ownAddress, targetContract})...)

	erc20Address := web3.HexToAddress(erc20)
	res, err := client.Eth().Call(&web3.CallMsg{
		From:     web3.HexToAddress(ownAddress),
		To:       &erc20Address,
		Data:     data,
		GasPrice: 1,
		Value:    big.NewInt(0),
	}, -1)

	allowanceBz, err := hex.DecodeString(res[2:])
	if err != nil {
		return false, err
	}
	allowance := new(big.Int).SetBytes(allowanceBz)
	if allowance.Cmp(big.NewInt(1000000000000)) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func ApproveErc20Transfer(erc20, gravityContract string, privKey *ecdsa.PrivateKey, client *jsonrpc.Client) (*web3.Hash, error) {
	ethAddress := crypto.PubkeyToAddress(privKey.PublicKey).String()

	sig, err := utils.ParseSignature(strings.Replace("approve(address,uint256)", " ", "", -1))
	if err != nil {
		return nil, err
	}
	payload := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{gravityContract, big.NewInt(1000000000000)})...)

	sendTx := gravity_utils.Exec(func() interface{} {
		hash ,err := SendTransaction(client, erc20, payload, big.NewInt(0), ethAddress, privKey)
		if err != nil {
			return err
		}
		return hash
	}).Await()

	hash, ok := sendTx.(*web3.Hash)
	if !ok {
		return nil, sendTx.(error)
	}

	//wait for tx done
	timeout := 30 * time.Second

	var topics []*web3.Hash
	sig, err = utils.ParseSignature(strings.Replace("Approval(address,address,uint256)", " ", "", -1))
	if err != nil {
		return nil, err
	}

	sigBz := utils.EventID(sig.Method, sig.Args)
	eventHash := web3.HexToHash("0x" + hex.EncodeToString(sigBz))
	topics = append(topics, &eventHash)

	gravity_utils.Exec(func() interface{} {
		log, err := WaitForEvent(timeout, erc20, "Approval(address,address,uint256)", topics, client)
		if err != nil {
			return err
		}
		return log
	}).Await()

	return hash, nil
}

func WaitForEvent(timeout time.Duration, contractAddress, event string, topics []*web3.Hash, client *jsonrpc.Client) ([]*web3.Log, error) {
	filter := &web3.LogFilter{
		Address:   []web3.Address{web3.HexToAddress(contractAddress)},
		Topics:    topics,
		BlockHash: nil,
		From:      nil,
		To:        nil,
	}

	gravity_utils.Exec(func() interface{} {
		time.Sleep(timeout)
		return nil
	}).Await()

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