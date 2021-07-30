package rest

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"net/http"
)


type KeyValue struct {
	Type  string  `json:"type"`
	Value string   `json:"value"`
}

type Validator struct {
	Address string 	  `json:"address"`
	PubKey  KeyValue  `json:"pub_key"`
	PriKey  KeyValue  `json:"pri_key"`
}

func RegisterRestNodeRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/node/new_validator", CreateNewValidatorKey(cliCtx)).Methods("POST")
}


func CreateNewValidatorKey(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		validatorKey := ed25519.GenPrivKey()
		validatorPubKey := validatorKey.PubKey()

		cdc := amino.NewCodec()
		keyByte, err := cdc.MarshalJSON(validatorKey)
		if err != nil {
			rest.WriteErrorRes(w, "cdc marshal validatorKey failed")
		}
		pubKeyByte, err := cdc.MarshalJSON(validatorPubKey)
		if err != nil {
			rest.WriteErrorRes(w, "cdc marshal validatorKey failed")
		}
		address := validatorPubKey.Address().String()
		resp := Validator{
			Address: address,
			PubKey:  KeyValue{
				Type:  ed25519.PubKeyName,
				Value: string(pubKeyByte[1:len(pubKeyByte)-1]),
			},
			PriKey:  KeyValue{
				Type:  ed25519.PrivKeyName,
				Value: string(keyByte[1:len(keyByte)-1]),
			},
		}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}