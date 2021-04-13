package validator

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

type PubKey struct {
	Type string `json:"type"`
	Value string `json:"value"`
}

func NewValidatorKey() (validatorKey, pubKeyStr, address string, err error) {
	var valKey ed25519.PrivKey
	validator := ed25519.GenPrivKey()
	cdc := amino.NewCodec()
	keyByte, err := cdc.MarshalJSON(validator)
	if err != nil {
		return "","", "", err
	}
	validatorKey = string(keyByte[1:len(keyByte)-1])
	privStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, ed25519.PrivKeyName, validatorKey)
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