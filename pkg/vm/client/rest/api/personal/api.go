package personal

import (
	"os"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth"
)

// PrivateAccountAPI is the personal_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PrivateAccountAPI struct {
	ethAPI   *eth.PublicEthereumAPI
	logger   log.Logger
	keys     map[common.Address]string
}

// NewAPI creates an instance of the public Personal Eth API.
func NewAPI(ethAPI *eth.PublicEthereumAPI, keys map[common.Address]string) *PrivateAccountAPI {
	api := &PrivateAccountAPI{
		ethAPI: ethAPI,
		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc", "namespace", "personal"),
		keys: keys,
	}

	return api
}

// ImportRawKey armors and encrypts a given raw hex encoded ECDSA key and stores it into the key directory.
// The name of the key will have the format "personal_<length-keys>", where <length-keys> is the total number of
// keys stored on the keyring.
// NOTE: The key will be both armored and encrypted using the same passphrase.
func (api *PrivateAccountAPI) ImportRawKey(privkey string) (common.Address, error) {
	api.logger.Debug("personal_importRawKey")
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	address := common.HexToAddress(crypto.PubkeyToAddress(key.PublicKey).Hex())
	api.keys[address] = privkey

	return address, nil
}

// LockAccount will lock the account associated with the given address when it's unlocked.
// It removes the key corresponding to the given address from the API's local keys.
func (api *PrivateAccountAPI) LockAccount(address common.Address) bool {
	api.logger.Debug("personal_lockAccount", "address", address.String())

	keys := api.ethAPI.GetKeys()
	if _, exists := keys[address]; exists {
		delete(keys, address)
		return true
	}
	return false
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
// It exports the private key corresponding to the given address from the keyring and stores it in the API's local keys.
func (api *PrivateAccountAPI) UnlockAccount(address common.Address) (bool, error) { // nolint: interfacer
	api.logger.Debug("personal_unlockAccount", "address", address.String())
	// TODO: use duration

	if key, exists := api.keys[address]; exists {
		keys := api.ethAPI.GetKeys()
		keys[address] = key
		api.ethAPI.SetKeys(keys)
		return true, nil
	}

	return false, nil
}
