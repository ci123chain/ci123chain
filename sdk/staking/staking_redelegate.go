package staking

import (
	"encoding/hex"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/staking"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func SignRedelegateMsg(from string, amount, gas, nonce uint64, priv string,
	validatorSrcAddress, validatorDstAddress, delegatorAddress string) ([]byte, error) {
	//
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	fromAddr, amt, validatorSrcAddr, validatorDstAddr, delegatorAddr, err := RedelegateParseArgs(from, amount, validatorSrcAddress, validatorDstAddress, delegatorAddress)
	if err != nil {
		return nil, err
	}
	tx := staking.NewRedelegateMsg(fromAddr, gas, nonce, delegatorAddr, validatorSrcAddr, validatorDstAddr, amt)

	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}

func httpRedelegateTx(from, gas, nonce, Type, name, height, priv, validatorSrcAddr, validatorDstAddr, delegatorAddr, proxy string) {
	cli := &http.Client{}
	reqUrl := "http://ciChain:3030/staking/redelegate"
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("privateKey", priv)
	data.Set("validatorSrcAddr", validatorSrcAddr)
	data.Set("validatorDstAddr", validatorDstAddr)
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