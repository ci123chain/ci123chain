package rest

import (
	"crypto"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/p2p"
	"io/ioutil"
	"net/http"
	"path/filepath"
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

type PubKeyInfo struct {
	Type   string   `json:"type"`
	Value  crypto.PublicKey `json:"value"`
}

type PrivKeyInfo struct{
	Type   string    `json:"type"`
	Value  crypto.PrivateKey `json:"value"`
}

type NodeInfo struct {
	Address  string    `json:"address"`
	PubKey   PubKeyInfo  `json:"pub_key"`
	PriKey   PrivKeyInfo `json:"pri_key"`
}

func RegisterRestNodeRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/node/new_validator", CreateNewValidatorKey(cliCtx)).Methods("POST")
	r.HandleFunc("/node/get_node_info", GetNodeInfo(cliCtx)).Methods("GET")
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

func GetNodeInfo(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		root := viper.GetString("HOME")
		jsonBytes, err := ioutil.ReadFile(filepath.Join(root, "config/node_key.json"))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		nodeKey := new(p2p.NodeKey)
		err = cdc.UnmarshalJSON(jsonBytes, nodeKey)
		validatorKey := nodeKey.PrivKey
		validatorPubKey := validatorKey.PubKey()
		address := validatorPubKey.Address().String()
		resp := NodeInfo{
				Address:address,
				PubKey: PubKeyInfo{
					Type:  ed25519.PubKeyName,
					Value: validatorPubKey,
				},
				PriKey: PrivKeyInfo{
					Type:  ed25519.PrivKeyName,
					Value: validatorKey,
				},
			}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}