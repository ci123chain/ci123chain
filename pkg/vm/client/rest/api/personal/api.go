package personal

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"math"
	"os"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth"
)

// PrivateAccountAPI is the personal_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PrivateAccountAPI struct {
	ethAPI   *eth.PublicEthereumAPI
	logger   log.Logger
	ks       *keystore.KeyStore
}

// NewAPI creates an instance of the public Personal Eth API.
func NewAPI(ethAPI *eth.PublicEthereumAPI, ks *keystore.KeyStore) *PrivateAccountAPI {
	api := &PrivateAccountAPI{
		ethAPI: ethAPI,
		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "json-rpc", "namespace", "personal"),
		ks: ks,
	}
	return api
}

// ImportRawKey armors and encrypts a given raw hex encoded ECDSA key and stores it into the key directory.
// The name of the key will have the format "personal_<length-keys>", where <length-keys> is the total number of
// keys stored on the keyring.
// NOTE: The key will be both armored and encrypted using the same passphrase.
func (api *PrivateAccountAPI) ImportRawKey(privkey string, password string) (common.Address, error) {
	api.logger.Debug("personal_importRawKey")
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := api.ks.ImportECDSA(key, password)
	if err != nil {
		return common.Address{}, err
	}
	return acc.Address, nil
}

// LockAccount will lock the account associated with the given address when it's unlocked.
// It removes the key corresponding to the given address from the API's local keys.
func (api *PrivateAccountAPI) LockAccount(address common.Address) bool {
	api.logger.Debug("personal_lockAccount", "address", address.String())
	return api.ks.Lock(address) == nil
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
// It exports the private key corresponding to the given address from the keyring and stores it in the API's local keys.
func (api *PrivateAccountAPI) UnlockAccount(addr common.Address, password string, duration *uint64) (bool, error) { // nolint: interfacer
	api.logger.Debug("personal_unlockAccount", "address", addr.String())
	// TODO: use duration
	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}
	err := api.ks.TimedUnlock(accounts.Account{Address: addr}, password, d)
	if err != nil {
		return false, err
	}

	return true, nil
}