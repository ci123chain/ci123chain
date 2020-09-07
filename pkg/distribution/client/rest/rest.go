package rest

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/gorilla/mux"
	"net/http"
)
var cdc = app.MakeCodec()

func RegisterRoutes(cliCtx context.Context, r *mux.Router) {
	RegisterQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

func checkDelegatorAddressVar(w http.ResponseWriter, r *http.Request) (string, bool) {
	address := r.FormValue("delegatorAddr")
	checkErr := util.CheckStringLength(42, 100, address)
	if checkErr != nil {
		rest.WriteErrorRes(w,types.ErrBadAddress(types.DefaultCodespace, checkErr))
		return "", false
	}
	return address, true
}

func checkValidatorAddressVar(w http.ResponseWriter, r *http.Request) (string, bool) {
	address := r.FormValue("validatorAddr")
	checkErr := util.CheckStringLength(42, 100, address)
	if checkErr != nil {
		rest.WriteErrorRes(w,types.ErrBadAddress(types.DefaultCodespace, checkErr))
		return "", false
	}
	return address, true
}


func checkWithdrawAddressVar(w http.ResponseWriter, r *http.Request) (string, error) {
	address := r.FormValue("withdrawAddress")
	checkErr := util.CheckStringLength(42, 100, address)
	if checkErr != nil {
		return "", checkErr
	}
	return address, nil
}



/*func checkAccountAddressVar(w http.ResponseWriter, r *http.Request) (sdk.AccAddress, bool) {
	address := r.FormValue("accountAddress")
	checkErr := util.CheckStringLength(42, 100, address)
	if checkErr != nil {
		rest.WriteErrorRes(w,types.ErrBadAddress(types.DefaultCodespace, checkErr))
		return sdk.AccAddress{}, false
	}
	return sdk.HexToAddress(address), true
}*/

func checkFromAddressVar(w http.ResponseWriter, r *http.Request) (string, bool) {
	address := r.FormValue("from")
	checkErr := util.CheckStringLength(42, 100, address)
	if checkErr != nil {
		return "", false
	}
	return address, true
}

func checkAmountVar(w http.ResponseWriter, r *http.Request) (int64, bool) {
	amount := r.FormValue("amount")
	amt, checkErr := util.CheckInt64(amount)
	if checkErr != nil {
		return 0, false
	}
	return amt, true
}

func checkGasVar(w http.ResponseWriter, r *http.Request) (uint64, bool) {
	gas := r.FormValue("gas")
	Gas, checkErr := util.CheckUint64(gas)
	if checkErr != nil {
		return 0, false
	}
	return Gas, true
}

func checkNonce(w http.ResponseWriter, r *http.Request, from sdk.AccAddress) (uint64, error) {

	nonce := r.FormValue("nonce")
	var Nonce uint64
	if nonce == "" {
		//
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return 0, err
		}
		var Err error
		Nonce, _, Err = ctx.GetNonceByAddress(from, false)
		if Err != nil {

			return 0, Err
		}
	}else {
		var checkErr error
		Nonce, checkErr = util.CheckUint64(nonce)
		if checkErr != nil {
			return -0, checkErr
		}
	}
	return Nonce, nil
}

func checkPrivateKey(w http.ResponseWriter, r *http.Request) (string, bool) {
	privKey := r.FormValue("privateKey")
	if privKey == "" {
		return "", false
	}
	return privKey, true
}