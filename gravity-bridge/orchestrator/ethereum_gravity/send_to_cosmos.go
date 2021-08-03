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

func SendToCosmos(erc20, gravityContract string,
	amount *big.Int,
	cosmosDestination common.Address,
	senderSecret *ecdsa.PrivateKey,
	waitTimeout time.Duration,
	client *jsonrpc.Client,
) (*web3.Hash, error) {
	senderAddress := crypto.PubkeyToAddress(senderSecret.PublicKey)
	getProved := gravity_utils.Exec(func() interface{} {
		proved, err := CheckErc20Approved(erc20, senderAddress.String(), gravityContract, client)
		if err != nil {
			return err
		}
		return proved
	}).Await()

	proved, ok := getProved.(bool)
	if !ok {
		return nil, getProved.(error)
	}

	if !proved {
		hash := gravity_utils.Exec(func() interface{} {
			hash, err := ApproveErc20Transfer(erc20, gravityContract, senderSecret, client)
			if err != nil {
				return err
			}
			return hash
		}).Await()

		gravity_utils.Exec(func() interface{} {
			tx, err := WaitForTransaction(client, hash.(*web3.Hash), waitTimeout, big.NewInt(0))
			if err != nil {
				return err
			}
			return tx
		}).Await()
	}


	sig, err := utils.ParseSignature(strings.Replace("sendToCosmos(address,address,uint256)", " ", "", -1))
	if err != nil {
		return nil, err
	}
	payload := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{erc20, cosmosDestination.String(), amount})...)

	sendTx := gravity_utils.Exec(func() interface{} {
		hash ,err := SendTransaction(client, gravityContract, payload, big.NewInt(0), senderAddress.String(), senderSecret)
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
