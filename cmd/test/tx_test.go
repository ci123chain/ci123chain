package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/ci123chain/pkg/order"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"io/ioutil"
	"net/http"
	"net/url"
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
var End = 300


func makePrivateKey() []byte {
	priKey, _ := cryptoutil.DecodePriv([]byte(testPrivKey))
	privByte := cryptoutil.MarshalPrivateKey(priKey)
	return privByte
}

func MakeParams(i int, pri []byte) string{
	nonce := uint64(i)
	privByte := pri
	signdata, err := transfer.SignTransferTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 1, 20000, nonce, privByte)
	if err != nil {
		panic(err)
	}
	req := hex.EncodeToString(signdata)

	return req
}


func myFunc(i interface{}, ph string) {
	n := i.(int)
	fmt.Println(n)
	http.PostForm("http://127.0.0.1:131" + ph + "/tx/broadcast_async",
		url.Values{"data": {TxRequestParam[n]}})
}

func TestProcess(t *testing.T) {
	var ph = "0"
	var wg sync.WaitGroup
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
	var ph = "1"
	var wg sync.WaitGroup
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

func TestProcessTwo(t *testing.T) {
	var ph = "2"
	var wg sync.WaitGroup
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

func TestProcessThree(t *testing.T) {
	var ph = "3"
	var wg sync.WaitGroup
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
	var ph = "4"
	var wg sync.WaitGroup
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


func TestAddShard(t *testing.T) {
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)


	signdata, err := order.SignUpgradeTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		20000, 1, "ADD", "test-chain-tL56l9", 29, privByte)

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
	resp, err := http.PostForm("http://127.0.0.1:1310/tx/addShard",
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