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

func EstimateValsetCost(
	latestCosmosValset,
	currentValSet types.ValSet,
	latestCosmosConfirmed []types.ValsetConfirmResponse,
	client *jsonrpc.Client,
	contractAddr,
	gravityId string,
	ethKey *ecdsa.PrivateKey,
	) (GasCost, error){
	ourEthAddress := web3.HexToAddress(crypto.PubkeyToAddress(ethKey.PublicKey).String())
	to := web3.HexToAddress(contractAddr)
	//ourBalance, _ := client.Eth().GetBalance(ourEthAddress, -1)
	//ourNonce, _ := client.Eth().GetNonce(ourEthAddress, -1)
	//var gasLimit uint64
	//if ourBalance.Cmp(ourBalance.SetUint64(math.MaxUint64)) > 0 {
	//	gasLimit = math.MaxUint64
	//} else {
	//	gasLimit = ourBalance.Uint64()
	//}
	gasPrice, _ := client.Eth().GasPrice()
	payload, _ := encodeValsetPayload(latestCosmosValset, currentValSet, latestCosmosConfirmed, gravityId)
	val, err := client.Eth().EstimateGas(&web3.CallMsg{
		From:     ourEthAddress,
		To:       &to,
		Data:     payload,
		GasPrice: gasPrice,
		Value:    big.NewInt(0),
	})
	if err != nil {
		fmt.Println(err)
	}

	return GasCost{
		Gas:      big.NewInt(int64(val)),
		GasPrice: big.NewInt(int64(gasPrice)),
	}, nil
}

func SendEthValsetUpdate(
	newValset,
	oldValset types.ValSet,
	confirms []types.ValsetConfirmResponse,
	client *jsonrpc.Client,
	timeout time.Duration,
	gravityContractAddr string,
	gravityId string,
	ethKey *ecdsa.PrivateKey,
	) error {
	lg := logger.GetLogger()

	oldNonce := oldValset.Nonce
	newNonce := newValset.Nonce
	if newNonce <= oldNonce {
		panic("newNocne <= oldNonce")
	}

	ourEthAddress := crypto.PubkeyToAddress(ethKey.PublicKey)
	beforeNonce, _ := GetValSetNonce(gravityContractAddr, crypto.PubkeyToAddress(ethKey.PublicKey), client)
	if beforeNonce != oldNonce {
		lg.Info(fmt.Sprintf("Someone else updated the valset to %d, exiting early", beforeNonce))
		return nil
	}

	payload, _ := encodeValsetPayload(newValset, oldValset, confirms, gravityId)
	sendTx := gravity_utils.Exec(func() interface{} {
		hash ,err := SendTransaction(client, gravityContractAddr, payload, big.NewInt(0), ourEthAddress.String(), ethKey)
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

	lastNonce, err := GetValSetNonce(gravityContractAddr, ourEthAddress, client)
	if err != nil {
		return err
	}

	if lastNonce != newNonce {
		lg.Error(fmt.Sprintf("Current nonce is %d expected to update to nonce %d", lastNonce, newNonce))
	} else {
		lg.Info(fmt.Sprintf("Successfully updated Valset with new Nonce %d", lastNonce))
	}

	return nil
}

func encodeValsetPayload(
	newValset,
	oldValset types.ValSet,
	confirms []types.ValsetConfirmResponse,
	gravityId string,
) ([]byte, error) {
	oldAddresses, oldPowers := oldValset.FilterEmptyAddress()
	newAddresses, newPowers := newValset.FilterEmptyAddress()
	oldValsetNonce := oldValset.Nonce
	newValsetNonce := newValset.Nonce
	hash := types.EncodeValsetConfirmHashed(gravityId, newValset)

	var confirm []types.Confirm
	for _, v := range confirms {
		confirm = append(confirm, v)
	}

	sigData, err := oldValset.OrderSigs(hash, confirm)
	if err != nil {
		return nil, err
	}
	sigArrays := types.ToArrays(sigData)

	sig, err := utils.ParseSignature(strings.Replace("updateValset(address[],uint256[],uint256,address[],uint256[],uint256,uint8[],bytes32[],bytes32[])", " ", "", -1))
	if err != nil {
		return nil, err
	}
	data := append(utils.MethodID(sig.Method, sig.Args), utils.RawEncode(sig.Args, []interface{}{newAddresses, newPowers, newValsetNonce, oldAddresses, oldPowers, oldValsetNonce, sigArrays.V, sigArrays.R, sigArrays.S})...)
	return data, nil
}