package validator

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type PubKey struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

func NewValidatorKey() (validatorKey, pubKeyStr, address string, err error) {
	var valKey secp256k1.PrivKeySecp256k1
	validator := secp256k1.GenPrivKey()
	cdc := amino.NewCodec()
	keyByte, err := cdc.MarshalJSON(validator)
	if err != nil {
		return "","", "", err
	}
	validatorKey = string(keyByte[1:len(keyByte)-1])
	privStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyAminoName, validatorKey)
	cdc = types.MakeCodec()
	err = cdc.UnmarshalJSON([]byte(privStr), &valKey)
	if err != nil {
		return "","", "", err
	}
	valPubKey := valKey.PubKey()
	address = valPubKey.Address().String()

	var pubKey PubKey
	pb, _ := cdc.MarshalJSON(valPubKey)
	err = json.Unmarshal(pb, &pubKey)
	if err != nil {
		return "", "", "", err
	}
	pubKeyStr = pubKey.Value
	return
}