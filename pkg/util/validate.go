package util

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"net/http"
	"strconv"
)

func CheckBool(async string) (bool, error) {
	if async == "" {
		return true, nil
	}
	isAysnc, err := strconv.ParseBool(async)
	if err != nil {
		return true, err
	}
	return isAysnc, nil
}


func CheckInt64(num string) (int64, error) {
	if num == "" {
		return 0, errors.New("it is empty")
	}
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func CheckUint64(num string) (uint64, error) {
	if num == "" {
		return 0, errors.New("it is empty")
	}
	n, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}
//check length of string

func CheckStringLength(min, max int, str string) error {
	if str == "" {
		return errors.New("empty string")
	}

	length := len(str)
	if max == -1 {
		if length < min {
			return errors.New("unexpected length")
		}else {
			return nil
		}
	}else {
		if length < min || length > max {
			return errors.New("unexpected length")
		}
		return nil
	}
}
//check json string

func CheckJsonArgs(str string, param interface{}) (bool, error) {
	if str == "" {
		return false, errors.New("empty string")
	}
	b := []byte(str)
	err := json.Unmarshal(b, &param)
	if err != nil {
		return false, errors.New("error byte")
	}
	return true, nil
}

func CheckFromAddressVar(r *http.Request) (string, bool) {
	address := r.FormValue("from")
	checkErr := CheckStringLength(42, 100, address)
	if checkErr != nil {

		return "", false
	}
	return address, true
}

func CheckAmountVar(r *http.Request) (int64, bool) {
	amount := r.FormValue("amount")
	amt, checkErr := CheckInt64(amount)
	if checkErr != nil {
		return 0, false
	}
	return amt, true
}

func CheckGasVar(r *http.Request) (uint64, bool) {
	gas := r.FormValue("gas")
	Gas, checkErr := CheckUint64(gas)
	if checkErr != nil {
		return 0, false
	}
	return Gas, true
}

func CheckNonce(r *http.Request, from sdk.AccAddress, cdc *codec.Codec) (uint64, bool) {
	nonce := r.FormValue("nonce")
	var Nonce uint64
	if nonce == "" {
		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return 0, false
		}
		Nonce, err = ctx.GetNonceByAddress(from)
		if err != nil {
			return 0, false
		}
	}else {
		var checkErr error
		Nonce, checkErr = CheckUint64(nonce)
		if checkErr != nil {
			return 0, false
		}
	}
	return Nonce, true
}

func CheckPrivateKey(r *http.Request) (string, bool) {
	privKey := r.FormValue("privateKey")
	if privKey == "" {
		return "", false
	}
	return privKey, true
}