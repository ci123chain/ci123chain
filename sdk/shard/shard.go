package shard

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ci123chain/ci123chain/pkg/order"
)

//off line
func SignAddShardMsg(from string, gas, nonce uint64,t, name string, height int64, priv string) ([]byte, error){

	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	tx := order.NewAddShardTx(fromAddr, gas, nonce, t, name, height)
	sid := cryptosuit.NewFabSignIdentity()
	privPub, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	pub, err  := sid.GetPubKey(privPub)
	if err != nil {
		return nil, err
	}

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privPub)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

//on line
func HttpAddShardTx(from, gas, nonce, Type, name, height, priv, proxy, reqUrl string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("type", Type)
	data.Set("name", name)
	data.Set("height", height)
	data.Set("privateKey", priv)
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