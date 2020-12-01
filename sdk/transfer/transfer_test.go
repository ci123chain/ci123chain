package transfer

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
	to = "0x505A74675dc9C71eF3CB5DF309256952917E801e"
	amount = "2"
	//offlineAmount = uint64(2)
	gas = "20000"
	offlineGas = uint64(20000)
	nonce = "2"
	offlineNonce = uint64(2)
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	proxy = "lb"
	reqUrl = "http://ciChain:3030/tx/broadcast"
	onlineReqUrl = "http://ciChain:3030/tx/transfers"
)

func TestSignMsgTransfer(t *testing.T) {

	signdata, err := SignMsgTransfer(from, to, offlineGas, offlineNonce, amount, priv, false)

	assert.NoError(t, err)
	httpTransfer(hex.EncodeToString(signdata))
}
func httpTransfer(param string) {
	cli := &http.Client{}
	data := url.Values{}
	data.Set("txByte", param)
	data.Set("proxy", "lb")
	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
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

func TestSendTransferMSg(t *testing.T) {
	httpSendTransferMSg(from, to, amount, gas, nonce, priv, proxy, onlineReqUrl)
}

func httpSendTransferMSg(from, to, amount, gas, nonce, priv, proxy, reqUrl string) {
	//
	cli := &http.Client{}
	data := url.Values{}
	data.Set("from", from)
	data.Set("to", to)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("amount", amount)
	data.Set("privateKey", priv)
	data.Set("proxy", proxy)


	req2, err := http.NewRequest("POST", reqUrl, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Body = ioutil.NopCloser(strings.NewReader(data.Encode()))

	// set request content type
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