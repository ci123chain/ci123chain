package ibc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"github.com/tendermint/tendermint/rpc/client"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
	order "github.com/tanhuiya/ci123chain/pkg/order"
)

var testPrivKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

var ip = "127.0.0.1"
var port = "1317"

// 生成 跨链 交易
func TestIBCMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串

	//获取nonce，地址， 端口
	nonce := httpQuery(ip, port, FromAddr)

	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)


	signdata, err := SignIBCTransferMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 6, 20000, nonce, privByte)
	fmt.Println(hex.EncodeToString(signdata))

	//assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}


// 生成 apply 签名交易
const UniqueID = "61849E3829B6B42616BC2736FA44CBBE"
const ObserverID = "1234567812345679"
func TestApplyIBCMsg(t *testing.T)  {
	nonce := httpQuery(ip, port, FromAddr)

	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	uid := []byte(UniqueID)
	signdata, err := SignApplyIBCMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", uid, []byte(ObserverID), 1, nonce, privByte)

	assert.NoError(t, err)
	//fmt.Println(hex.EncodeToString(signdata))
	httpPost(hex.EncodeToString(signdata))
}


// 生成 bank 到 account 扣钱交易
const pkg =
`{"signature":"MEQCIF8Lp+L1p3p0WV5tf68QdN7Mf04mLMLIx9j0FSOefhE1AiBEcPxKCIu7y0B09R5IUHmy5tXeiX5rhPWTXXFyd2Tljw==","ibc_msg_bytes":"eyJ1bmlxdWVfaWQiOiI4QTM3QkE2QjgwMTNBQkU1OUYyNzhFRkUzM0Q1QjE4OCIsIm9ic2VydmVyX2lkIjoiMTIzNDU2NzgxMjM0NTY3OCIsImJhbmtfYWRkcmVzcyI6IjB4NTA1QTc0Njc1ZGM5QzcxZUYzQ0I1REYzMDkyNTY5NTI5MTdFODAxZSIsImFwcGx5X3RpbWUiOiIyMDE5LTExLTEzVDE2OjU4OjMyLjAwNzkyKzA4OjAwIiwic3RhdGUiOiJwcm9jZXNzaW5nIiwiZnJvbV9hZGRyZXNzIjoiMHgyMDRiQ0M0MjU1OUZhZjZERkUxNDg1MjA4Rjc5NTFhYUQ4MDBCMzEzIiwidG9fYWRkcmVzcyI6IjB4RDFhMTQ5NjI2MjdmQWM3NjhGZTg4NUVlYjlGRjA3MjcwNkI1NGMxOSIsImFtb3VudCI6MTB9"}`

func TestBankSendMsg(t *testing.T)  {
	nonce := httpQuery(ip, port, FromAddr)

	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)
	pub := priKey.Public().(*ecdsa.PublicKey)
	addr, _  := cryptoutil.PublicKeyToAddress(pub)

	signdata, err := SignIBCBankSendMsg(addr, []byte(pkg), 1000000, nonce, privByte)


	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}

const pkgReceipt =
`{"unique_id":"8A37BA6B8013ABE59F278EFE33D5B188","observer_id":"1234567812345678","signature":"MEQCIHNGnWa/xk4n+WOERiXphkytHN+iOIfQiwJTozixLBXnAiBPyZbtpUHwbp4ATUg90Tmye/iNZ9sc7Q3jO0RsZ1Jfag=="}`

func TestReceiptMsg(t *testing.T)  {
	nonce := httpQuery(ip, port, FromAddr)

	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignIBCReceiptMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", []byte(pkgReceipt), 1, nonce, privByte)
	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}

const FromAddr  = "0x204bCC42559Faf6DFE1485208F7951aaD800B313"
const ToAddr  = "0xD1a14962627fAc768Fe885Eeb9FF072706B54c19"


type ciRes struct{
	Ret 	uint32 	`json:"ret"`
	Data 	string	`json:"data"`
	Message	string	`json:"message"`
}

func TestAll(t *testing.T)  {

	nonce := httpQuery(ip, port, FromAddr)
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	fmt.Println("---发送跨链消息")
	nonce = httpQuery(ip, port, FromAddr)
	signdata, err := SignIBCTransferMsg(FromAddr,
		ToAddr, 20, 50000, nonce, privByte)
	registRet := httpPost(hex.EncodeToString(signdata))


	fmt.Println("发送跨链消息完成：UniqueID = " + registRet.Data)
	fmt.Println()

	fmt.Println("---申请处理该跨链消息")
	nonce = httpQuery(ip, port, FromAddr)
	uid := []byte(registRet.Data)
	signdata, err = SignApplyIBCMsg(FromAddr, uid, []byte(ObserverID), 50000, nonce, privByte)
	applyRet := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(applyRet.RawLog) < 1)
	fmt.Println("申请处理该跨链消息结束")
	fmt.Println()


	fmt.Println("---第二个申请处理该跨链消息")
	nonce = httpQuery(ip, port, FromAddr)
	ObserverID2 := "12313213213213124321"
	signdata, err = SignApplyIBCMsg(FromAddr, uid, []byte(ObserverID2), 50000, nonce, privByte)

	applyRetErr := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(applyRetErr.RawLog) > 0)
	fmt.Println("申请处理该跨链消息结束失败")
	fmt.Println()


	// bank转账，该交易参数应该是observer 从 fabric 获得，此处是模拟
	fmt.Println("---向对方转账")
	nonce = httpQuery(ip, port, FromAddr)
	pkg := applyRet.Data
	signdata, err = SignIBCBankSendMsg(FromAddr, []byte(pkg), 50000, nonce, privByte)
	fmt.Println(hex.EncodeToString(signdata))
	receiptRet := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(receiptRet.RawLog) < 1)
	fmt.Println("向对方转账成功")
	fmt.Println()

	//发送回执
	fmt.Println("---发送回执")
	nonce = httpQuery(ip, port, FromAddr)
	receivepkg := receiptRet.Data
	signdata, err = SignIBCReceiptMsg(FromAddr, []byte(receivepkg), 50000, nonce, privByte)
	ret := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(ret.RawLog) < 1)
	fmt.Println("发送回执成功")
}

func httpQuery(ip, port, param string) uint64 {
	url := "http://" + ip + ":" + port + "/ibctx/nonce/" + param
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var ret string
	err = json.Unmarshal(body, &ret)
	reps, err := strconv.ParseInt(ret, 10, 64)
	nonce := uint64(reps)

	return nonce
}
func httpPost(param string) retData {
	resp, err := http.PostForm("http://127.0.0.1:1317/tx/broadcast",
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


type retData struct {
	Data string `json:"data"`
	RawLog  string `json:"raw_log"`
}


func TestUpgradeTx(t *testing.T) {

	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)


	signdata, err := order.SignUpgradeTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		20000, 1, "ADD", "asdjqj", 35, privByte)

	assert.NoError(t, err)
	httpPostUpgradeTx(hex.EncodeToString(signdata))
}

func httpPostUpgradeTx(param string) retData{
	resp, err := http.PostForm("http://127.0.0.1:1317/tx/addShard",
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


func TestNginx(t *testing.T) {
	// 将pem 格式私钥转化为 十六机制 字符串

	//获取nonce，地址， 端口
	//nonce := httpQuery(ip, port, FromAddr)

	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	var signdata []byte
	for i := 1; i <= 100; i++ {
		nonce := uint64(i)
		signdata, err = transfer.SignTransferTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
			"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 1, 20000, nonce, privByte)

		assert.NoError(t, err)
		httpPostReqAsync(hex.EncodeToString(signdata))
	}

	/*signdata, err := transfer.SignTransferTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 1, 20000, 0, privByte)
	fmt.Println(hex.EncodeToString(signdata))

	assert.NoError(t, err)
	httpPostReq(hex.EncodeToString(signdata))*/
}

func httpPostReqAsync(param string) retData{

	resp, err := http.PostForm("http://127.0.0.1:8080/tx/broadcast_async",
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

//定义一个实现Job接口的数据
type Score struct {
	Num string
}
//定义对数据的处理
func (s *Score) Do() {
	http.PostForm("http://127.0.0.1:8080/tx/broadcast_async",
		url.Values{"data": {s.Num}})
}

/*func TestRoutine(t *testing.T) {
	num := 200
	// 注册工作池，传入任务
	// 参数1 worker并发个数
	p := c.NewWorkerPool(num)
	p.Run()

	req := make(map[int]string, 20000)
	for j := 1; j <= 5000; j++ {
		//fmt.Println(j)
		req[j] = co(j)
		//req = append(req, req[j])
	}

	//写入一亿条数据
	datanum := 5000
	go func() {
		for i := 1; i <= datanum; i++ {
			//byte := co(i)
			sc := &Score{Num: req[i]}
			p.JobQueue <- sc //数据传进去会被自动执行Do()方法，具体对数据的处理自己在Do()方法中定义
		}
	}()

	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine())
		time.Sleep(2 * time.Second)
	}
}*/

func Test20(t *testing.T) {
	req := make(map[int]string, 80000)
	for j := 1; j <= 20000; j++ {

		req[j] = co(j)
		fmt.Println(j)
		fmt.Println(req[j])
	}
}


func co(i int) string{
	nonce := uint64(i)
	priKey, _ := cryptoutil.DecodePriv([]byte(testPrivKey))
	privByte := cryptoutil.MarshalPrivateKey(priKey)
	signdata, err := transfer.SignTransferTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 1, 20000, nonce, privByte)
	if err != nil {
		panic(err)
	}
	req := hex.EncodeToString(signdata)

	return req
}


var req map[int]string
var rep map[int]string

var Client1 *client.HTTP
var Client2 *client.HTTP


func myFunc1(i interface{}) {
	//n := i.(int32)
	n := i.(int)
	//atomic.AddInt32(&sum, n)
	//fmt.Printf("run with %d\n", n)
	//fmt.Println(req[n])
	param, _ := hex.DecodeString(req[n])

	//Client := client.NewHTTP("http://0.0.0.0:26607", "/http")

	//client.BroadcastTxAsync([]byte("123"))

	Client1.BroadcastTxAsync(param)
	//defer Client.Stop()
}

func myFunc2(i interface{}) {
	//n := i.(int32)
	n := i.(int)
	//atomic.AddInt32(&sum, n)
	//fmt.Printf("run with %d\n", n)
	//fmt.Println(req[n])
	param, _ := hex.DecodeString(rep[n])

	//Client := client.NewHTTP("http://0.0.0.0:26607", "/http")

	//client.BroadcastTxAsync([]byte("123"))

	Client2.BroadcastTxAsync(param)
	//defer Client.Stop()
}

func TestTcentRoutine(t *testing.T) {
	var runTimes = 10000
	var wg sync.WaitGroup
	Client1 = client.NewHTTP("http://0.0.0.0:26607", "/http")

	//req := make(map[int]string, 100000)
	req = make(map[int]string, 20000)
	for j := 501; j <= 10000; j++ {
		req[j] = co(j)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc1(i)
		wg.Done()
	})
	defer p1.Release()
	// Submit tasks one by one.
	for i := 501; i < runTimes; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p1.Running())
	//fmt.Printf("finish all tasks, result is %d\n", sum)
	for {
		fmt.Printf("running goroutines: %d\n", p1.Running())
		time.Sleep(2 * time.Second)
	}
}

func TestTwoRoutine(t *testing.T) {
	var runTimes = 20000
	var wg sync.WaitGroup
	Client2 = client.NewHTTP("http://0.0.0.0:26607", "/http")

	//req := make(map[int]string, 100000)
	rep = make(map[int]string, 50000)
	for j := 1; j <= 20000; j++ {
		rep[j] = co(j)
	}
	p2, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc2(i)
		wg.Done()
	})
	defer p2.Release()
	// Submit tasks one by one.
	for i := 10001; i < runTimes; i++ {
		wg.Add(1)
		_ = p2.Invoke(i)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p2.Running())
	//fmt.Printf("finish all tasks, result is %d\n", sum)
	for {
		fmt.Printf("running goroutines: %d\n", p2.Running())
		time.Sleep(2 * time.Second)
	}
}


func BenchmarkRoutine(b *testing.B) {
	//var runTimes = 10000
	var wg sync.WaitGroup
	Client1 = client.NewHTTP("http://0.0.0.0:26607", "/http")

	//req := make(map[int]string, 100000)
	req = make(map[int]string, 30000)
	for j := 2; j <= 10000; j++ {
		req[j] = co(j)
	}
	p1, _ := ants.NewPoolWithFunc(100, func(i interface{}) {
		myFunc1(i)
		wg.Done()
	})
	defer p1.Release()
	b.ResetTimer()

	/*b.RunParallel(func(pb *testing.PB){
		for i := 1;;pb.Next() {
			wg.Add(1)
			_ = p1.Invoke(i)
		}
	})*/
	// Submit tasks one by one.
	for i := 1; i < b.N; i++ {
		wg.Add(1)
		_ = p1.Invoke(i)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p1.Running())
	//fmt.Printf("finish all tasks, result is %d\n", sum)
	/*for {
		fmt.Printf("running goroutines: %d\n", p1.Running())
		time.Sleep(2 * time.Second)
	}*/
}


func Test11(t *testing.T) {
	client := client.NewHTTP("http://0.0.0.0:26607", "/http")

	health, err := client.Health()
	fmt.Println(health, err)
	client.BroadcastTxAsync([]byte("123"))
}