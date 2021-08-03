package ethereum_gravity

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"math/big"
	"strings"
	"time"
)

func DeployErc20(
	cosmosDenom,
	erc20Name,
	erc20Symbol string,
	decimals uint8,
	gravityContract common.Address,
	client *jsonrpc.Client,
	waitTimeout time.Duration,
	senderSecret *ecdsa.PrivateKey,
	) (*web3.Hash, error) {
	senderAddress := crypto.PubkeyToAddress(senderSecret.PublicKey)

	sig, err := utils.ParseSignature(strings.Replace("deployERC20(string,string,string,uint8)", " ", "", -1))
	if err != nil {
		return nil, err
	}
	payload := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{cosmosDenom, erc20Name, erc20Symbol, decimals})...)

	sendTx := gravity_utils.Exec(func() interface{} {
		hash ,err := SendTransaction(client, gravityContract.String(), payload, big.NewInt(0), senderAddress.String(), senderSecret)
		if err != nil {
			return err
		}
		return hash
	}).Await()

	hash, ok := sendTx.(*web3.Hash)
	if !ok {
		return nil, sendTx.(error)
	}

	gravity_utils.Exec(func() interface{} {
		tx, err := WaitForTransaction(client, hash, waitTimeout, big.NewInt(0))
		if err != nil {
			return err
		}
		return tx
	}).Await()



	return hash, nil
}