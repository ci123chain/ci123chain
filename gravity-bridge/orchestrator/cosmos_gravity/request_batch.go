package cosmos_gravity

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types3 "github.com/ci123chain/ci123chain/pkg/app/types"
	types2 "github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func SendRequestBatch(
	privKey *ecdsa.PrivateKey,
	denom string,
	contact Contact,
) (sdk.TxResponse, error) {
	ourAddress := crypto.PubkeyToAddress(privKey.PublicKey)
	msgRequestBatch := &types2.MsgRequestBatch{
		Sender: ourAddress.String(),
		Denom:  denom,
	}

	nonce := contact.GetNonce(ourAddress.String())

	msgs := []sdk.Msg{msgRequestBatch}

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