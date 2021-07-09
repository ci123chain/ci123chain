package helper

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)
// KeyOutput contains mnemonic and address of key
type KeyOutput struct {
	PrivateKey  string `json:"privatekey" yaml:"privatekey"`
	Address  string `json:"address" yaml:"address"`
}

// KeyAddOrRestore is a helper function for add key and restores key when mnemonic is passed
func KeyAddOrRestore(privatekeys ...string) (KeyOutput, error) {
	var privatekeyStr string
	var err error

	if len(privatekeys) > 0 {
		privatekeyStr = privatekeys[0]
	} else {
		return KeyOutput{}, errors.New("privateKey can not be empty")

	}
	privateKey, err := crypto.HexToECDSA(privatekeyStr)
	if err != nil {
		return KeyOutput{}, errors.Errorf("error format privateKey: %s", privatekeyStr)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	ko := KeyOutput{PrivateKey: privatekeyStr, Address: sdk.ToAccAddress(address[:]).String()}

	return ko, nil
}
