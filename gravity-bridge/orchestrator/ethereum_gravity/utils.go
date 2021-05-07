package ethereum_gravity

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
)

func GetGravityId(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) ([]byte, error) {
	return nil, nil
}

func GetValSetNonce(contractAddr string, ourEthereumAddress common.Address, client *jsonrpc.Client) (uint64, error) {
	return 0, nil
}

func CheckForEvents(startBlock, endBlock uint64, contractAddr string, events []string, client *jsonrpc.Client) ([]*web3.Log, error) {
	return nil, nil
}
