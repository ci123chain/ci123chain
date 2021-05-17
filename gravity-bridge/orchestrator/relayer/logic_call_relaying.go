package relayer

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

//need to implement
func relayLogicCalls(
	currentValSet types.ValSet,
	ethKey *ecdsa.PrivateKey,
	client *jsonrpc.Client,
	contact cosmos_gravity.Contact,
	contractAddr,
	gravityId string,
	timeout time.Duration,
) {
	//lg := logger.GetLogger()
	//ourEthereumAddress := crypto.PubkeyToAddress(ethKey.PublicKey)
	//latestCalls, err := cosmos_gravity.GetLatestLogicCalls(contact)
	//if err != nil {
	//	return
	//}
	//
	//var oldestSignedCall types.LogicCall
	//var oldestSignatures types.LogicCallConfirmResponse
	//for _, call := range latestCalls {
	//	sigs, err := cosmos_gravity.GetLogicCallSignatures(contact, call.InvalidationId, call.InvalidationNonce)
	//	if err != nil {
	//		lg.Error(fmt.Sprintf("could not get signatures for %s/%s, error: %s", hex.EncodeToString(call.InvalidationId), strconv.FormatUint(call.InvalidationNonce, 10), err.Error()))
	//	}
	//
	//	hash := types.EncodeValsetConfirm()
	//
	//}

}