package relayer

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/umbracle/go-web3/jsonrpc"
	"time"
)

func relayValsets(currentValSet types.ValSet, ethKey *ecdsa.PrivateKey, client *jsonrpc.Client, contact cosmos_gravity.Contact, contractAddr, gravityId string, timeout time.Duration) {

}
