package rest

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/util"
	"io/ioutil"
	"net/http"
)

var cdc = types.MakeCodec()

func getCode(r *http.Request) (code []byte, err error){
	codeStr := r.FormValue("code_str")
	if codeStr != "" {
		Byte, err := hex.DecodeString(codeStr)
		if err != nil {
			return nil, errors.New("invalid code")
		}
		code = Byte
	}else {
		file, _, err := r.FormFile("code_file")
		if err != nil {
			return nil, errors.New("cannot get code file: " + err.Error())
		}
		code, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}
	return
}

func adjustInstantiateParams(r *http.Request) (name, version, author, email, describe string, err error) {
	name, err = validParam(r, "name")
	if err != nil {
		return
	}
	version, err = validParam(r, "version")
	if err != nil {
		return
	}
	author, err = validParam(r, "author")
	if err != nil {
		return
	}
	email, err = validParam(r, "email")
	if err != nil {
		return
	}
	describe, err = validParam(r, "describe")
	if err != nil {
		return
	}
	return
}

func validParam(r *http.Request, value string) (param string, err error) {
	param = r.FormValue(value)
	if param == "" {
		param = value
	}else {
		err := util.CheckStringLength(1, 100, param)
		if err != nil {
			return "", errors.New(fmt.Sprintf("error %s", value))
		}
	}
	return
}