package cosmos_gravity

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

func SendValsetConfirms(contact Contact,
	ethPrivKey *ecdsa.PrivateKey,
	fee sdk.Coin,
	valsets []*types.ValSet,
	cosmosPrivateKey *ecdsa.PrivateKey,
	gravityId string) (sdk.TxResponse, error) {

	return sdk.TxResponse{}, nil
}

func SendEthereumClaims(contact Contact,
	cosmosPrivKey *ecdsa.PrivateKey,
	deposits []types.SendToCosmosEvent,
	withdraws []types.TransactionBatchExecutedEvent,
	erc20Deploys []types.Erc20DeployedEvent,
	logicCalls []types.LogicCallExecutedEvent,
	fee sdk.Coin) (sdk.TxResponse, error) {

	return sdk.TxResponse{}, nil
}
