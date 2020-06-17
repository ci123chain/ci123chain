package rest

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/account/new", NewAccountRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balance", QueryBalancesRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/account/nonce", QueryNonceRequestHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/account/new_validator", CreateNewValidatorKeyHandleFn(cliCtx)).Methods("GET")
}

type BalanceData struct {
	Balance uint64 		`json:"balance"`
	Proof *merkle.Proof `json:"proof"`
}

type NonceData struct {
	Nonce   uint64   `json:"nonce"`
	Proof *merkle.Proof `json:"proof"`
}

type AccountAddress struct {
	Address string `json:"address"`
	Height  string `json:"height"`
}

type Account struct {
	Address string `json:"address"`
	PrivKey string `json:"privKey"`
}

type Key struct {
	ValidatorKey string `json:"validator_key"`
}

func QueryBalancesRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		address := request.FormValue("address")
		height := request.FormValue("height")
		checkErr := util.CheckStringLength(42, 100, address)
		if checkErr != nil {
			rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, checkErr))
			return
		}
		if height != "" {
			_, Err := util.CheckInt64(height)
			if Err != nil {
				rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, Err))
				return
			}
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, request, height)
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}
		addrBytes, err2 := helper.ParseAddrs(address)
		if len(addrBytes) < 1 || err2 != nil {
			rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, err2))
			return
		}
		//params := types.NewQueryBalanceParams(addr)
		res, proof, err2 := cliCtx.GetBalanceByAddress(addrBytes[0])
		if err2 != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		resp := BalanceData{Balance:res, Proof: proof}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func NewAccountRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		key, err := crypto.GenerateKey()
		if err != nil {
			fmt.Println("Error: ", err.Error());
		}

		address := crypto.PubkeyToAddress(key.PublicKey).Hex()
		privKey := hex.EncodeToString(key.D.Bytes())

		resp := Account{
			Address:	address,
			PrivKey:	privKey,
		}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func QueryNonceRequestHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		address := r.FormValue("address")
		checkErr := util.CheckStringLength(42, 100, address)
		if checkErr != nil {
			rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, checkErr))
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok {
			rest.WriteErrorRes(w, err)
			return
		}
		addrBytes, err2 := helper.ParseAddrs(address)
		if len(addrBytes) < 1 || err2 != nil {
			rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, err2))
			return
		}
		res, proof, err2 := cliCtx.GetNonceByAddress(addrBytes[0])
		if err2 != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		resp := NonceData{Nonce:res, Proof:proof}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func CreateNewValidatorKeyHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		validatorKey := secp256k1.GenPrivKey()

		cdc := amino.NewCodec()
		keyByte, err := cdc.MarshalJSON(validatorKey)
		if err != nil {
			rest.WriteErrorRes(w, client.ErrGenValidatorKey(types.DefaultCodespace, err))
		}
		resp := Key{ValidatorKey:string(keyByte[1:len(keyByte)-1])}

		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}