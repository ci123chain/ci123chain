package cosmos_gravity

import (
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ethereum/go-ethereum/common"
)

func GetValSet(contact Contact, valSetNonce uint64) (*types.ValSet, error) {
	return nil, nil
}

func GetOldestUnsignedValsets(contact Contact, address common.Address) ([]*types.ValSet, error) {
	return nil, nil
}
func GetLastEventNonce(ourCosmosAddress common.Address, contact Contact) (uint64, error) {
	//QueryLastEventNonceByAddrRequest
	return 0, nil
}