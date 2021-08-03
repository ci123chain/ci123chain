package cosmos_gravity

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	gravity "github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
)

func GetValSet(contact Contact, valSetNonce uint64) (*types.ValSet, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/valset_request/%d", gravity.StoreKey, valSetNonce))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var cval gravity.Valset
	if err = json.Unmarshal(data, &cval); err != nil {
		return nil, err
	}

	var member []types.ValSetMember
	for _, m := range cval.Members {
		addr := common.HexToAddress(m.EthereumAddress)
		member = append(member, types.ValSetMember{
			Power:      m.Power,
			EthAddress: &addr,
		})
	}

	oval := &types.ValSet{
		Nonce:   cval.Nonce,
		Members: member,
	}
	return oval, nil
}

func GetAllValsetConfirms(contact Contact, latestNonce uint64) ([]types.ValsetConfirmResponse, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/valset_confirm/%d", gravity.StoreKey, latestNonce))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var confirms []*gravity.MsgValsetConfirm
	if err = json.Unmarshal(data, &confirms); err != nil {
		return nil, err
	}

	var response []types.ValsetConfirmResponse
	for _, confirm := range confirms {
		sigBz, err := hex.DecodeString(confirm.Signature)
		if err != nil {
			return nil, err
		}
		sig, err := types.FromBytesToEthSignature(sigBz)
		if err != nil {
			return nil, err
		}
		confirmResponse := types.ValsetConfirmResponse{
			Nonce:        confirm.Nonce,
			Orchestrator: common.HexToAddress(confirm.Orchestrator),
			EthAddress:   common.HexToAddress(confirm.EthAddress),
			Signature:    sig,
		}

		response = append(response, confirmResponse)
	}
	return response, nil
}

func GetLatestValSets(contact Contact) ([]*types.ValSet, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/valset_requests", gravity.StoreKey))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var cval []gravity.Valset
	var oval []*types.ValSet
	if err = json.Unmarshal(data, &cval); err != nil {
		return nil, err
	}

	for _, v := range cval {
		var member []types.ValSetMember
		for _, m := range v.Members {
			addr := common.HexToAddress(m.EthereumAddress)
			member = append(member, types.ValSetMember{
				Power:      m.Power,
				EthAddress: &addr,
			})
		}

		oval = append(oval, &types.ValSet{
			Nonce:   v.Nonce,
			Members: member,
		})
	}

	return oval, nil
}

func GetOldestUnsignedValsets(contact Contact, address common.Address) ([]types.ValSet, error) {
	//QueryLastPendingValsetRequestByAddr
	res, err := contact.Get(fmt.Sprintf("/%s/pending_valset_requests/%s", gravity.StoreKey, address.String()))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var cval []gravity.Valset
	var oval []types.ValSet
	if err = json.Unmarshal(data, &cval); err != nil {
		return nil, err
	}

	for _, v := range cval {
		var member []types.ValSetMember
		for _, m := range v.Members {
			addr := common.HexToAddress(m.EthereumAddress)
			member = append(member, types.ValSetMember{
				Power:      m.Power,
				EthAddress: &addr,
			})
		}

		oval = append(oval, types.ValSet{
			Nonce:   v.Nonce,
			Members: member,
		})
	}

	return oval, nil
}

func GetOldestUnsignedBatch(contact Contact, address common.Address) (*types.TransactionBatch, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/pending_batch_requests/%s", gravity.StoreKey, address.String()))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var cbatch gravity.OutgoingTxBatch
	if err = json.Unmarshal(data, &cbatch); err != nil {
		return nil, err
	}
	if cbatch.Block == 0 {
		return nil, nil
	}

	var batchTransactions []types.BatchTransaction
	totalFee := types.Erc20Token{
		Amount:               sdk.NewInt(0),
		TokenContractAddress: common.Address{},
	}
	for _, v := range cbatch.Transactions {
		if totalFee.Amount.IsZero() {
			totalFee = types.Erc20Token{
				Amount:               v.Erc20Fee.Amount,
				TokenContractAddress: common.HexToAddress(v.Erc20Fee.Contract),
			}
		} else {
			totalFee.Amount = totalFee.Amount.Add(v.Erc20Fee.Amount)
		}
		batchTransactions = append(batchTransactions, types.BatchTransaction{
			Id:          v.Id,
			Sender:      common.HexToAddress(v.Sender),
			Destination: common.HexToAddress(v.DestAddress),
			Erc20Token:  types.Erc20Token{
				Amount:               v.Erc20Token.Amount,
				TokenContractAddress: common.HexToAddress(v.Erc20Token.Contract),
			},
			Erc20Fee:    types.Erc20Token{
				Amount:               v.Erc20Fee.Amount,
				TokenContractAddress: common.HexToAddress(v.Erc20Fee.Contract),
			},
		})
	}

	obatch := &types.TransactionBatch{
		Nonce:         cbatch.BatchNonce,
		BatchTimeout:  cbatch.BatchTimeout,
		Transactions:  batchTransactions,
		TotalFee:      totalFee,
		TokenContract: common.HexToAddress(cbatch.TokenContract),
	}

	return obatch, nil
}


func GetLastEventNonce(contact Contact, ourCosmosAddress common.Address) (uint64, error) {
	//QueryLastEventNonceByAddrRequest
	res, err := contact.Get(fmt.Sprintf("/%s/last_event_nonce/%s", gravity.StoreKey, ourCosmosAddress.String()))
	if err != nil {
		return 0, err
	}

	var nonceRes rest.Response
	var nonce uint64
	json.Unmarshal(res, &nonceRes)
	json.Unmarshal(nonceRes.Data, &nonce)
	return nonce, nil
}

func GetLatestTransactionBatches(contact Contact) ([]types.TransactionBatch, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/transaction_batches", gravity.StoreKey))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var batchs []*gravity.OutgoingTxBatch
	if err = json.Unmarshal(data, &batchs); err != nil {
		return nil, err
	}

	var response []types.TransactionBatch
	for _, batch := range batchs {
		totalFee := types.Erc20Token{
			Amount:               sdk.NewInt(0),
			TokenContractAddress: common.Address{},
		}
		var transactions []types.BatchTransaction
		for _, tx := range batch.Transactions {
			transaction := types.BatchTransaction{
				Id:          tx.Id,
				Sender:      common.HexToAddress(tx.Sender),
				Destination: common.HexToAddress(tx.DestAddress),
				Erc20Token:  types.Erc20Token{
					Amount:               tx.Erc20Token.Amount,
					TokenContractAddress: common.HexToAddress(tx.Erc20Token.Contract),
				},
				Erc20Fee:    types.Erc20Token{
					Amount:               tx.Erc20Fee.Amount,
					TokenContractAddress: common.HexToAddress(tx.Erc20Fee.Contract),
				},
			}
			if totalFee.Amount.IsZero() {
				totalFee = types.Erc20Token{
					Amount:               tx.Erc20Fee.Amount,
					TokenContractAddress: common.HexToAddress(tx.Erc20Fee.Contract),
				}
			} else {
				totalFee.Amount = totalFee.Amount.Add(tx.Erc20Fee.Amount)
			}
			transactions = append(transactions, transaction)
		}
		batchResponse := types.TransactionBatch{
			Nonce:         batch.BatchNonce,
			BatchTimeout:  batch.BatchTimeout,
			Transactions:  transactions,
			TotalFee:      totalFee,
			TokenContract: common.HexToAddress(batch.TokenContract),
		}
		response = append(response, batchResponse)
	}
	return response, nil
}

func GetTransactionBatchSignatures(contact Contact, nonce uint64, tokenContract common.Address) ([]types.BatchConfirmResponse, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/batch_confirm/%d/%s", gravity.StoreKey, nonce, tokenContract.String()))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var confirms []*gravity.MsgConfirmBatch
	if err = json.Unmarshal(data, &confirms); err != nil {
		return nil, err
	}

	var response []types.BatchConfirmResponse
	for _, confirm := range confirms {
		sigBz, err := hex.DecodeString(confirm.Signature)
		if err != nil {
			return nil, err
		}
		sig, err := types.FromBytesToEthSignature(sigBz)
		if err != nil {
			return nil, err
		}
		confirmResponse := types.BatchConfirmResponse{
			Nonce:          confirm.Nonce,
			Orchestrator:   common.HexToAddress(confirm.Orchestrator),
			TokenContract:  common.HexToAddress(confirm.TokenContract),
			EthereumSigner: common.HexToAddress(confirm.EthSigner),
			EthSignature:   sig,
		}

		response = append(response, confirmResponse)
	}
	return response, nil
}

func GetLatestLogicCalls(contact Contact) ([]types.LogicCall, error){
	res, err := contact.Get(fmt.Sprintf("/%s/lastLogicCalls", gravity.StoreKey))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var calls []*gravity.OutgoingLogicCall
	if err = json.Unmarshal(data, &calls); err != nil {
		return nil, err
	}
	var response []types.LogicCall
	for _, call := range calls {
		var transfers, fees []types.Erc20Token
		for _, tx := range call.Transfers {
			transfer := types.Erc20Token{
				Amount:               tx.Amount,
				TokenContractAddress: common.HexToAddress(tx.Contract),
			}
			transfers = append(transfers, transfer)
		}
		for _, cost := range call.Fees {
			fee := types.Erc20Token{
				Amount:               cost.Amount,
				TokenContractAddress: common.HexToAddress(cost.Contract),
			}
			fees = append(fees, fee)
		}
		batchResponse := types.LogicCall{
			Transfers:            transfers,
			Fees:    			  fees,
			LogicContractAddress: common.HexToAddress(call.LogicContractAddress),
			PayLoad:			  call.Payload,
			Timeout:              call.Timeout,
			InvalidationId:       call.InvalidationId,
			InvalidationNonce:    call.InvalidationNonce,
		}
		response = append(response, batchResponse)
	}
	return response, nil
}

func GetLogicCallSignatures(contact Contact, invalidId []byte, invalidNonce uint64) ([]types.LogicCallConfirmResponse, error) {
	res, err := contact.Get(fmt.Sprintf("/%s/logicCall_confirm/%s/%s", gravity.StoreKey, hex.EncodeToString(invalidId), strconv.FormatUint(invalidNonce, 10)))
	if err != nil {
		return nil, err
	}
	data, err := resToData(res)
	if err != nil {
		return nil, err
	}
	var confirms []*gravity.MsgConfirmLogicCall
	if err = json.Unmarshal(data, &confirms); err != nil {
		return nil, err
	}

	var response []types.LogicCallConfirmResponse
	for _, confirm := range confirms {
		sigBz, err := hex.DecodeString(confirm.Signature)
		if err != nil {
			return nil, err
		}
		sig, err := types.FromBytesToEthSignature(sigBz)
		if err != nil {
			return nil, err
		}
		invalidId, err := hex.DecodeString(confirm.InvalidationId)
		if err != nil {
			return nil, err
		}
		confirmResponse := types.LogicCallConfirmResponse{
			InvalidationId:    invalidId,
			InvalidationNonce: confirm.InvalidationNonce,
			EthereumSigner:    common.HexToAddress(confirm.EthSigner),
			Orchestrator:      common.HexToAddress(confirm.Orchestrator),
			EthSignature:      sig,
		}
		response = append(response, confirmResponse)
	}
	return response, nil
}

func resToData(res []byte) (json.RawMessage, error) {
	var response rest.Response
	if err := json.Unmarshal(res, &response); err != nil {
		return nil, err
	}
	//if response.Ret != 1 {
	//	return nil, errors.New(response.Message)
	//}
	return response.Data, nil
}