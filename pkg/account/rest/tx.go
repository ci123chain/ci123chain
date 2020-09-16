package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	abcitype "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/account/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"net/http"
)

var cdc = app.MakeCodec()
// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/account/new", NewAccountRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balance", QueryBalancesRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/account/nonce", QueryNonceRequestHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/node/new_validator", rest.MiddleHandler(cliCtx, CreateNewValidatorKey, types.DefaultCodespace)).Methods("GET")
	r.HandleFunc("/transaction/multi_msgs_tx", rest.MiddleHandler(cliCtx, MultiMsgsRequest, types.DefaultCodespace)).Methods("POST")
}

type BalanceData struct {
	Balance uint64 	 `json:"balance"`
}

type NonceData struct {
	Nonce   uint64   `json:"nonce"`
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

func QueryBalancesRequestHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		address := request.FormValue("address")
		height := request.FormValue("height")
		prove := request.FormValue("prove")
		checkErr := util.CheckStringLength(42, 100, address)
		if checkErr != nil {
			rest.WriteErrorRes(w, client.ErrParseParam(types.DefaultCodespace, checkErr))
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, types.DefaultCodespace) {
			return
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
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, proof, err2 := cliCtx.GetBalanceByAddress(addrBytes[0], isProve)
		if err2 != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		value := BalanceData{Balance:res}
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func QueryNonceRequestHandleFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		address := r.FormValue("address")
		height := r.FormValue("height")
		prove := r.FormValue("prove")
		checkErr := util.CheckStringLength(42, 100, address)
		if checkErr != nil {
			rest.WriteErrorRes(w, client.ErrParseAddr(types.DefaultCodespace, checkErr))
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, types.DefaultCodespace) {
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
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, proof, err2 := cliCtx.GetNonceByAddress(addrBytes[0], isProve)
		if err2 != nil {
			rest.WriteErrorRes(w, transfer.ErrQueryTx(types.DefaultCodespace, err2.Error()))
			return
		}
		value := NonceData{Nonce:res}
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func CreateNewValidatorKey(cliCtx context.Context, w http.ResponseWriter, r *http.Request) {
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

func MultiMsgsRequest(cliCtx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg_str := r.FormValue("msgs")
	var msgs []abcitype.Msg
	var msgs_str []string
	err := json.Unmarshal([]byte(msg_str), &msgs_str)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrParseParam(types.DefaultCodespace, err))
		return
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, true)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrParseParam(types.DefaultCodespace, err))
		return
	}
	for _, v := range msgs_str{
		var msg abcitype.Msg
		msg_byte, err := hex.DecodeString(v)
		if err != nil {
			rest.WriteErrorRes(w, client.ErrParseParam(types.DefaultCodespace, err))
			return
		}
		err = cdc.UnmarshalBinaryLengthPrefixed(msg_byte, &msg)
		if err != nil {
			rest.WriteErrorRes(w, client.ErrParseParam(types.DefaultCodespace, err))
			return
		}
		msgs = append(msgs, msg)
	}
	txByte, err := app.SignCommonTx(from, nonce, gas, msgs, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrParseParam(types.DefaultCodespace, err))
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, client.ErrBroadcast(types.DefaultCodespace, err))
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}