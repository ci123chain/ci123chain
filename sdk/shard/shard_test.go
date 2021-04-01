package shard

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var (
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	gas = "20000"
	offlineGas = uint64(20000)
	nonce = "2"
	offlineNonce = uint64(2)
	ty = "ADD"
	name = "ci"
	offlineHeight = int64(800)
	height = "800"
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	reqUrl = "http://ciChain:3030/tx/broadcast"
	onelineReqUrl = "http://ciChain:3030/tx/addShard"
	proxy = "lb"
)

func TestSignUpgradeMsg(t *testing.T) {

	signdata, err := SignUpgradeMsg(ty, name, offlineHeight, from, offlineGas, offlineNonce, priv)
	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)
	httpSignUpgradeMsg(hex.EncodeToString(signdata))
}

func httpSignUpgradeMsg(param string) {

	cli := &http.Client{}
	data := url.Values{}
	data.Set("txByte", param)
	data.Set("proxy", "lb")
	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content types
	req2.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(b)
}

func TestSendAddShardMsg(t *testing.T) {
	httpSendAddShardMsg(from, gas, nonce, ty, name, height, priv, proxy, onelineReqUrl)

}

func httpSendAddShardMsg(from, gas, nonce, Type, name, height, priv, proxy, reqUrl string) {
	//
	cli := &http.Client{}

	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("types", Type)
	data.Set("name", name)
	data.Set("height", height)
	data.Set("privateKey", priv)
	data.Set("proxy", proxy)


	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content types
	req2.Header.Set("Content-Type", "x-www-form-urlencoded")
	// request
	rep2, err := cli.Do(req2)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(b)
}