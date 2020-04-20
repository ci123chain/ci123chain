package staking

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/staking"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SignDelegateMsg(from string, amount, gas, nonce uint64, priv string,
	validatorAddress, delegatorAddress string) ([]byte, error) {

	//
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	fromAddr, amt, validatorAddr, delegatorAddr, err := CommonParseArgs(from, amount, validatorAddress, delegatorAddress)
	if err != nil {
		return nil, err
	}
	tx := staking.NewDelegateMsg(fromAddr, gas, nonce, delegatorAddr, validatorAddr,amt)

	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)

	return tx.Bytes(), nil

}

func HttpDelegateTx(from, gas, nonce, amount,priv, validatorAddr, delegatorAddr, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("privateKey", priv)
	data.Set("amount", amount)
	data.Set("validatorAddr", validatorAddr)
	data.Set("delegatorAddr", delegatorAddr)
	data.Set("proxy", proxy)


	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
	req.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep, err := cli.Do(req)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(b))
}