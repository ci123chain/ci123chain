package cosmos_gravity

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

func SendValsetConfirms(contact Contact,
	ethPrivateKey *ecdsa.PrivateKey,
	fee sdk.Coin,
	valsets []*types.ValSet,
	cosmosPrivateKey *ecdsa.PrivateKey,
	gravityId string) (sdk.TxResponse, error) {

	return sdk.TxResponse{}, nil
}
