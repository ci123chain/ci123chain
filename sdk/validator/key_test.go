package validator

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestNewValidatorKey(t *testing.T) {
	validatorKey, pubKey, address, err := NewValidatorKey()
	if err != nil{
		fmt.Println(err)
	}
	var valKey crypto.PubKey
	pubStr := fmt.Sprintf(`{"types":"%s","value":"%s"}`, secp256k1.PubKeyAminoName, pubKey)
	cdc := types.GetCodec()
	err = cdc.UnmarshalJSON([]byte(pubStr), &valKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(validatorKey)
	fmt.Println(pubKey)
	fmt.Println(address)
	fmt.Println(valKey)
}