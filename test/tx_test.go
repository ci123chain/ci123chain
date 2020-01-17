package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	order "github.com/tanhuiya/ci123chain/pkg/order/types"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

var testPrivKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

var TxRequestParam = make(map[int]string, 800000)
var Start = 1
var End = 30000

var Client = client.NewHTTP("127.0.0.1:26607", "/http")
var Client1 = client.NewHTTP("127.0.0.1:26617", "/http")
var Client2 = client.NewHTTP("127.0.0.1:26627", "/http")



func makePrivateKey() []byte {
	priKey, _ := cryptoutil.DecodePriv([]byte(testPrivKey))
	privByte := cryptoutil.MarshalPrivateKey(priKey)
	return privByte
}

func MakeParams(i int, pri []byte) string{
	nonce := uint64(i)
	privByte := pri
	signdata, err := order.SignUpgradeTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		20000, nonce, "ADD", "ty8", 8000, privByte)



	/*signdata, err := transfer.SignTransferTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313", "0x505A74675dc9C71eF3CB5DF309256952917E801e", 1, 20000,
		nonce, privByte)*/
	if err != nil {
		panic(err)
	}

	//assert.NoError(t, err)
	req := hex.EncodeToString(signdata)

	return req
}

func TestSign(t *testing.T) {
	key := makePrivateKey()
	res := MakeParams(18, key)
	fmt.Println(res)
}


func myFunc(i interface{}, ph string) {
	n := i.(int)
	fmt.Println(n)
	//ph := "0"
	//http.PostForm("http://127.0.0.1:131" + ph + "/tx/broadcast_async",
		//url.Values{"data": {TxRequestParam[n]}})

	param,_ := hex.DecodeString(TxRequestParam[n])
	Client.BroadcastTxAsync(param)
}

func myFunc1(i interface{}, ph string) {
	n := i.(int)
	fmt.Println(n)
	//ph := "0"
	//http.PostForm("http://127.0.0.1:131" + ph + "/tx/broadcast_async",
	//url.Values{"data": {TxRequestParam[n]}})

	param,_ := hex.DecodeString(TxRequestParam[n])
	Client1.BroadcastTxAsync(param)
}

func myFunc2(i interface{}, ph string) {
	n := i.(int)
	fmt.Println(n)
	//ph := "0"
	//http.PostForm("http://127.0.0.1:131" + ph + "/tx/broadcast_async",
	//url.Values{"data": {TxRequestParam[n]}})

	param,_ := hex.DecodeString(TxRequestParam[n])
	Client2.BroadcastTxAsync(param)
}

func TestProcess(t *testing.T) {
	var wg sync.WaitGroup
	var ph = "0"
	privateKey := makePrivateKey()
	for j := Start; j <= End; j++ {
		TxRequestParam[j] = MakeParams(j, privateKey)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc(i, ph)
		wg.Done()
	})
	defer p1.Release()
	// Submit tasks one by one.
	time1 := time.Now()
	for i := Start; i < End; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	time2 := time.Now().Sub(time1).Seconds()
	fmt.Println(time2)
	wg.Wait()
}

func TestProcessOne(t *testing.T) {

	var wg sync.WaitGroup
	var ph = "1"
	privateKey := makePrivateKey()
	for j := Start; j <= End; j++ {
		TxRequestParam[j] = MakeParams(j, privateKey)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc1(i, ph)
		wg.Done()
	})
	defer p1.Release()
	// Submit tasks one by one.
	time1 := time.Now()
	for i := Start; i < End; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	time2 := time.Now().Sub(time1).Seconds()
	fmt.Println(time2)
	wg.Wait()
}

func TestProcessTwo(t *testing.T) {
	var wg sync.WaitGroup
	var ph = "2"
	privateKey := makePrivateKey()
	for j := Start; j <= End; j++ {
		TxRequestParam[j] = MakeParams(j, privateKey)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc2(i, ph)
		wg.Done()
	})
	defer p1.Release()
	// Submit tasks one by one.
	time1 := time.Now()
	for i := Start; i < End; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	time2 := time.Now().Sub(time1).Seconds()
	fmt.Println(time2)
	wg.Wait()
}

func TestProcessThree(t *testing.T) {
	var wg sync.WaitGroup
	var ph = "3"
	privateKey := makePrivateKey()
	for j := Start; j <= End; j++ {
		TxRequestParam[j] = MakeParams(j, privateKey)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc(i, ph)
		wg.Done()
	})
	defer p1.Release()
	// Submit tasks one by one.
	time1 := time.Now()
	for i := Start; i < End; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	time2 := time.Now().Sub(time1).Seconds()
	fmt.Println(time2)
	wg.Wait()
}

func TestProcessFour(t *testing.T) {
	var wg sync.WaitGroup
	var ph = "4"
	privateKey := makePrivateKey()
	for j := Start; j <= End; j++ {
		TxRequestParam[j] = MakeParams(j, privateKey)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc(i, ph)
		wg.Done()
	})
	defer p1.Release()
	// Submit tasks one by one.
	time1 := time.Now()
	for i := Start; i < End; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	time2 := time.Now().Sub(time1).Seconds()
	fmt.Println(time2)
	wg.Wait()
}

func TestHexToString(t *testing.T) {

	b, _ := hex.DecodeString("64697374722f2f01")
	fmt.Println(string(b))
}

func TestAddShard(t *testing.T) {

	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)


	signdata, err := order.SignUpgradeTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",

		20000, 1, "ADD", "ci123chain-shared4", 1000, privByte)


	assert.NoError(t, err)
	httpPostUpgradeTx(hex.EncodeToString(signdata))
}

type retData struct {
	Data string `json:"data"`
	RawLog  string `json:"raw_log"`
}
type ciRes struct{
	Ret 	uint32 	`json:"ret"`
	Data 	string	`json:"data"`
	Message	string	`json:"message"`
}

func httpPostUpgradeTx(param string) retData{
	resp, err := http.PostForm("http://localhost:1310/tx/addShard",
		url.Values{"data": {param}})
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(body)
	var ret retData
	var ty ciRes
	err = json.Unmarshal(body, &ty)
	fmt.Println(ty.Data)
	d := []byte(ty.Data)
	err = json.Unmarshal(d, &ret)

	if err != nil {
		fmt.Println(string(body))
	}

	if len(ret.Data) > 0 {
		fmt.Println(ret.Data)
	}
	if len(ret.RawLog) > 0 {
		fmt.Println(ret.RawLog)
	}
	return ret
}

func TestFomate(t *testing.T) {

	var statedb = "couchdb://admin:password@192.168.2.89:5984"
	s := strings.Split(statedb, "://")
	auths := strings.Split(s[1], "@")

	info := auths[0]
	userpass := strings.Split(info, ":")
	fmt.Println(s[0])
	fmt.Println(s[1])
	fmt.Println(auths[0])
	fmt.Println(auths[1])
	fmt.Println(userpass[0])
	fmt.Println(userpass[1])

}