package register_delegate_keys

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types3 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ethereum/go-ethereum/common"
	"encoding/json"
)

func UpdateGravityDelegateAddresses(contact cosmos_gravity.Contact,
	ethAddress, cosmosAddress common.Address,
	privateKey *ecdsa.PrivateKey) (sdk.TxResponse, error) {
	msgSetOrchestratorAddress := &types.MsgSetOrchestratorAddress{
		Validator:    cosmosAddress.String(),
		Orchestrator: cosmosAddress.String(),
		EthAddress:   ethAddress.String(),
	}

	nonce := contact.GetNonce(cosmosAddress.String())

	msgs := []sdk.Msg{msgSetOrchestratorAddress}

	txBz, err := types3.SignCommonTx(sdk.HexToAddress(cosmosAddress.String()), nonce, cosmos_gravity.COMMON_GAS, msgs, hex.EncodeToString(privateKey.D.Bytes()), types3.GetCodec())
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