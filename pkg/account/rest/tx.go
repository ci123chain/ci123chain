package rest

import (
	"encoding/hex"
	"encoding/json"
	abcitype "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"net/http"
)

var cdc = types2.MakeCodec()
// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/account/new", NewAccountRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/bank/balance", QueryBalancesRequestHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/account/nonce", QueryNonceRequestHandleFn(cliCtx)).Methods("POST")
	r.HandleFunc("/node/new_validator", CreateNewValidatorKey(cliCtx)).Methods("POST")
	r.HandleFunc("/transaction/multi_msgs_tx", rest.MiddleHandler(cliCtx, MultiMsgsRequest, sdkerrors.RootCodespace)).Methods("POST")
}

type BalanceData struct {
	BalanceList interface{} 	 `json:"balance_list"`
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
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrResponse, "generate key failed").Error())
			return
		}
		if key == nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrResponse, "generate empty key").Error())
			return
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
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address").Error())
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, sdkerrors.RootCodespace) {
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, request, height)
		if !ok || err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "parse hetight faield").Error())
			return
		}
		addr, err2 := helper.StrToAddress(address)
		if err2 != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address").Error())
			return
		}
		//params := types.NewQueryBalanceParams(addr)
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, proof, err2 := cliCtx.GetBalanceByAddress(addr, isProve)
		if err2 != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, "query balance failed").Error())
			return
		}
		value := BalanceData{BalanceList:res}
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
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address").Error())
			return
		}
		if !rest.CheckHeightAndProve(w, height, prove, sdkerrors.RootCodespace) {
			return
		}

		cliCtx, ok, err := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r, "")
		if !ok || err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "prase height failed").Error())
			return
		}
		addrBytes, err2 := helper.ParseAddrs(address)
		if len(addrBytes) < 1 || err2 != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid address").Error())
			return
		}
		isProve := false
		if prove == "true" {
			isProve = true
		}
		res, proof, err2 := cliCtx.GetNonceByAddress(addrBytes[0], isProve)
		if err2 != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrResponse, "get nonce failed").Error())
			return
		}
		value := NonceData{Nonce:res}
		resp := rest.BuildQueryRes(height, isProve, value, proof)
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func CreateNewValidatorKey(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		validatorKey := ed25519.GenPrivKey()

		cdc := amino.NewCodec()
		keyByte, err := cdc.MarshalJSON(validatorKey)
		if err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, "cdc marshal failed").Error())
		}
		resp := Key{ValidatorKey:string(keyByte[1:len(keyByte)-1])}
		rest.PostProcessResponseBare(w, cliCtx, resp)
	}
}

func MultiMsgsRequest(cliCtx context.Context, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg_str := r.FormValue("msgs")
	var msgs []abcitype.Msg
	var msgs_str []string
	err := json.Unmarshal([]byte(msg_str), &msgs_str)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error()).Error())
		return
	}
	privKey, from, nonce, gas, err := rest.GetNecessaryParams(cliCtx, r, cdc, true)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
		return
	}
	for _, v := range msgs_str{
		var msg abcitype.Msg
		msg_byte, err := hex.DecodeString(v)
		if err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, err.Error()).Error())
			return
		}
		err = cdc.UnmarshalBinaryLengthPrefixed(msg_byte, &msg)
		if err != nil {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
			return
		}
		msgs = append(msgs, msg)
	}
	txByte, err := types2.SignCommonTx(from, nonce, gas, msgs, privKey, cdc)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	res, err := cliCtx.BroadcastSignedTx(txByte)
	if err != nil {
		rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
		return
	}
	rest.PostProcessResponseBare(w, cliCtx, res)
}