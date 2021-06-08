package cosmos_gravity

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types3 "github.com/ci123chain/ci123chain/pkg/app/types"
	types2 "github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	COMMON_GAS = 100000
)

func SendValsetConfirms(contact Contact,
	ethPrivKey *ecdsa.PrivateKey,
	fee sdk.Coin,
	valsets []types.ValSet,
	cosmosPrivKey *ecdsa.PrivateKey,
	gravityId string) (sdk.TxResponse, error) {

	ourCosmosAddress := crypto.PubkeyToAddress(cosmosPrivKey.PublicKey)
	ourEthAddress := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	var msgs []sdk.Msg
	lg := logger.GetLogger()
	for _, valset := range valsets {
		lg.Info(fmt.Sprintf("Submitting signature for valset: %v", valset.Nonce))
		msg := types.EncodeValsetConfirm(gravityId, valset)
		msgHash := types.GetEthereumMsgHash(msg)
		sig, err := crypto.Sign(msgHash[:], ethPrivKey)
		if err != nil {
			return sdk.TxResponse{}, err
		}
		lg.Info(fmt.Sprintf("Sending valset update with address %s and sig %s", ourCosmosAddress.String(), hex.EncodeToString(sig)))

		confirm := &types2.MsgValsetConfirm{
			Nonce:        valset.Nonce,
			Orchestrator: ourCosmosAddress.String(),
			EthAddress:   ourEthAddress.String(),
			Signature:    hex.EncodeToString(sig),
		}
		msgs = append(msgs, confirm)
	}

	nonce := contact.GetNonce(ourCosmosAddress.String())

	txBz, err := types3.SignCommonTx(sdk.HexToAddress(ourCosmosAddress.String()), nonce, COMMON_GAS, msgs, hex.EncodeToString(cosmosPrivKey.D.Bytes()), types3.GetCodec())
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var result sdk.TxResponse
	res := gravity_utils.Exec(func() interface{} {
		res := contact.BroadcastTx(txBz)
		return res
	}).Await().([]byte)
	json.Unmarshal(res, &result)

	return result, nil
}

func SendBatchConfirm(contact Contact,
	ethPrivKey *ecdsa.PrivateKey,
	fee sdk.Coin,
	transactionBatches []*types.TransactionBatch,
	cosmosPrivKey *ecdsa.PrivateKey,
	gravityId string) (sdk.TxResponse, error) {

	ourCosmosAddress := crypto.PubkeyToAddress(cosmosPrivKey.PublicKey)
	ourEthAddress := crypto.PubkeyToAddress(ethPrivKey.PublicKey)

	var msgs []sdk.Msg
	lg := logger.GetLogger()
	for _, batch := range transactionBatches {
		lg.Info(fmt.Sprintf("Submitting signature for batch: %v", batch.Nonce))
		msg := types.EncodeTxBatchConfirmHashed(gravityId, *batch)
		sig, err := crypto.Sign(msg[:], ethPrivKey)
		if err != nil {
			return sdk.TxResponse{}, err
		}
		lg.Info(fmt.Sprintf("Sending batch update with address %s and sig %s", ourCosmosAddress.String(), hex.EncodeToString(sig)))

		confirm := &types2.MsgConfirmBatch{
			Nonce:         batch.Nonce,
			TokenContract: batch.TokenContract.String(),
			EthSigner:     ourEthAddress.String(),
			Orchestrator:  ourCosmosAddress.String(),
			Signature:     hex.EncodeToString(sig),
		}
		msgs = append(msgs, confirm)
	}

	nonce := contact.GetNonce(ourCosmosAddress.String())

	txBz, err := types3.SignCommonTx(sdk.HexToAddress(ourCosmosAddress.String()), nonce, COMMON_GAS, msgs, hex.EncodeToString(cosmosPrivKey.D.Bytes()), types3.GetCodec())
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var result sdk.TxResponse
	res := gravity_utils.Exec(func() interface{} {
		res := contact.BroadcastTx(txBz)
		return res
	}).Await().([]byte)
	json.Unmarshal(res, &result)

	return result, nil
}

func SendEthereumClaims(contact Contact,
	cosmosPrivKey *ecdsa.PrivateKey,
	deposits []types.SendToCosmosEvent,
	withdraws []types.TransactionBatchExecutedEvent,
	erc20Deploys []types.Erc20DeployedEvent,
	logicCalls []types.LogicCallExecutedEvent,
	fee sdk.Coin) (sdk.TxResponse, error) {

	ourCosmosAddress := crypto.PubkeyToAddress(cosmosPrivKey.PublicKey)

	var msgs []sdk.Msg
	//lg := logger.GetLogger()
	for _, deposit := range deposits {
		claim := &types2.MsgDepositClaim{
			EventNonce:   deposit.EventNonce,
			BlockHeight:  deposit.BlockHeight,
			TokenContract: deposit.Erc20.String(),
			Amount: sdk.NewIntFromBigInt(deposit.Amount),
			CosmosReceiver: deposit.Destination.String(),
			EthereumSender: deposit.Sender.String(),
			Orchestrator: ourCosmosAddress.String(),
		}
		msgs = append(msgs, claim)
	}

	for _, withdraw := range withdraws {
		claim := &types2.MsgWithdrawClaim{
			EventNonce:   withdraw.EventNonce,
			BlockHeight:  withdraw.BlockHeight,
			TokenContract: withdraw.Erc20.String(),
			BatchNonce:  withdraw.BatchNonce,
			Orchestrator: ourCosmosAddress.String(),
		}
		msgs = append(msgs, claim)
	}

	for _, deploy := range erc20Deploys {
		claim := &types2.MsgERC20DeployedClaim{
			EventNonce:   deploy.EventNonce,
			BlockHeight:  deploy.BlockHeight,
			TokenContract: deploy.Erc20.String(),
			Name: deploy.Name,
			Symbol: deploy.Symbol,
			Decimals: uint64(deploy.Decimals),
			Orchestrator: ourCosmosAddress.String(),
		}
		msgs = append(msgs, claim)
	}

	for _, call := range logicCalls {
		claim := &types2.MsgLogicCallExecutedClaim{
			EventNonce:   call.EventNonce,
			BlockHeight:  call.BlockHeight,
			InvalidationId: call.InvalidationId,
			InvalidationNonce: call.InvalidationNonce,
			Orchestrator: ourCosmosAddress.String(),
		}
		msgs = append(msgs, claim)
	}

	//should sort msgs???

	nonce := contact.GetNonce(ourCosmosAddress.String())

	txBz, err := types3.SignCommonTx(sdk.HexToAddress(ourCosmosAddress.String()), nonce, COMMON_GAS, msgs, hex.EncodeToString(cosmosPrivKey.D.Bytes()), types3.GetCodec())
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var result sdk.TxResponse
	res := gravity_utils.Exec(func() interface{} {
		res := contact.BroadcastTx(txBz)
		return res
	}).Await().([]byte)
	json.Unmarshal(res, &result)

	return result, nil
}

func SendToEth(privKey *ecdsa.PrivateKey, destination *common.Address, amount sdk.Coin, fee sdk.Coin, contact Contact) (sdk.TxResponse, error) {
	ourAddress := crypto.PubkeyToAddress(privKey.PublicKey)
	if amount.Denom != fee.Denom {
		return sdk.TxResponse{}, errors.New("amount denom must equal fee denom")
	}
	balances := contact.GetBalance(ourAddress.String())

	found := false
	for _, balance := range balances {
		if balance.Denom == amount.Denom {
			feeTotal := fee.Amount.Mul(sdk.NewInt(2))
			totalAmount := amount.Amount.Add(feeTotal)
			if balance.Amount.LT(totalAmount) {
				return sdk.TxResponse{}, errors.New("Insufficient balance")
			}
			found = true
			break
		}
	}

	if !found {
		return sdk.TxResponse{}, errors.New("No balance to send")
	}

	msgSendToEth := &types2.MsgSendToEth{
		Sender:    ourAddress.String(),
		EthDest:   destination.String(),
		Amount:    amount,
		BridgeFee: fee,
	}

	nonce := contact.GetNonce(ourAddress.String())

	msgs := []sdk.Msg{msgSendToEth}

	txBz, err := types3.SignCommonTx(sdk.HexToAddress(ourAddress.String()), nonce, COMMON_GAS, msgs, hex.EncodeToString(privKey.D.Bytes()), types3.GetCodec())
	if err != nil {
		return sdk.TxResponse{}, err
	}

	var result sdk.TxResponse
	res := gravity_utils.Exec(func() interface{} {
		res := contact.BroadcastTx(txBz)
		return res
	}).Await().([]byte)
	json.Unmarshal(res, &result)

	return result, nil
}