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

func TestSignAddShardMsg(t *testing.T) {

	signdata, err := SignAddShardMsg("0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c",
		20000, 2, "ADD", "ty8", 8000, "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70")
	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)
	httpSignAddShardMsg(hex.EncodeToString(signdata))
}

func httpSignAddShardMsg(param string) {

	cli := &http.Client{}
	reqUrl := "http://ciChain:3030/tx/broadcast"
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

func TestSendAddShardMsg(t *testing.T) {
	from := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	gas := "20000"
	nonce := "2"
	Type := "ADD"
	name := "cichain-shard1"
	height := "900"
	priv := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	proxy := "lb"
	httpSendAddShardMsg(from, gas, nonce, Type, name, height, priv, proxy)

}

func httpSendAddShardMsg(from, gas, nonce, Type, name, height, priv, proxy string) {
	//
	cli := &http.Client{}
	///body := make([]byte, 0)
	reqUrl := "http://ciChain:3030/tx/addShard"
	data := url.Values{}
	data.Set("from", from)
	data.Set("gas", gas)
	data.Set("nonce", nonce)
	data.Set("type", Type)
	data.Set("name", name)
	data.Set("height", height)
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