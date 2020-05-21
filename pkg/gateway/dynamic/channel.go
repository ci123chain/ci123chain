package dynamic

import (
	"fmt"
	"github.com/pretty66/gosdk"
	"github.com/pretty66/gosdk/cienv"
	"net/http"
)

const APP_ID  = "cc2535f30a87457b91649ecbed0245e6"
const APP_KEY = "dd6a53073bb9447eacaca250cc708c12"
const CALL_APP_ID = "hedlzgp1u48kjf50xtcvwdklminbqe9a"
const CHANNEL = "2"

const CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_MULTIPART = "multipart/form-data"

func init() {
	err := cienv.SetEnv(gosdk.GATEWAY_SERVICE_KEY, "kong:http://127.0.0.1:13800")
	if err != nil {
		panic(err)
	}
}

// 动态添加channel, 向部署空间请求, 利用 go-sdk
func CreateChannel(appid, appkey, channel, version string, h http.Header) error {
	// 获取对象，head是请求的HEAD字段，用来解析HEAD中的Authorization中的token
	client, err := gosdk.GetClientInstance(h)
	if err != nil {
		return err
	}

	// 对Authorization中的token解析，或对SetToken()中token解析
	//一般服务调用时使用token解析
	//client, err = client.SetToken(token)
	// 如果调用方是应用，则通过SetAppInfo进行调用方的信息存储，服务不要使用该方法
	err = client.SetAppInfo(appid, appkey, channel, version)
	if err != nil {
		return err
	}

	// 可以使用SetServices()自定义服务地址，或通过serviceKey从环境变量中寻找服务地址（前者优先级高）
	// services是map[string]string，key是serviceKey，value是服务地址
	//client = client.SetServices(services)
	serviceKey := "mdocxiqnl43hu68a2lrayv9p5fttm0vf"
	method := "POST"
	api := "Channel/assignChannel"
	params := make(map[string]interface{})
	alias := "deployment"

	params["target_appid"] = CALL_APP_ID
	params["instance_name"] = "testDy"
	params["instance_properties"] = "test"
	params["from_id"] = "i1"
	params["from_type"] = "t1"
	params["extra_info"] = "{\"CI_STATEDB\":\"couchdb://admin:zhangchaotest@193.112.144.129:7984/ci123dev02\"}"


	resp, err := client.Call(serviceKey, method, api, params, alias, CONTENT_TYPE_FORM, nil)
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}

func TestCall(appid, appkey, channel, version string, h http.Header) error {
	// 获取对象，head是请求的HEAD字段，用来解析HEAD中的Authorization中的token
	client, err := gosdk.GetClientInstance(h)
	if err != nil {
		return err
	}

	// 对Authorization中的token解析，或对SetToken()中token解析
	//一般服务调用时使用token解析
	//client, err = client.SetToken(token)
	// 如果调用方是应用，则通过SetAppInfo进行调用方的信息存储，服务不要使用该方法
	err = client.SetAppInfo(appid, appkey, channel, version)
	if err != nil {
		return err
	}

	// 可以使用SetServices()自定义服务地址，或通过serviceKey从环境变量中寻找服务地址（前者优先级高）
	// services是map[string]string，key是serviceKey，value是服务地址
	//client = client.SetServices(services)
	serviceKey := "hedlzgp1u48kjf50xtcvwdklminbqe9a"
	appKey := "6f395f69aef7479488d490f77de20e74"
	channelS := "2"

	method := "POST"
	api := "/bank/balance"
	params := make(map[string]interface{})

	params["address"] = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"

	resp, err := client.CallServiceInstance(serviceKey, appKey, channelS, method, api, params, CONTENT_TYPE_FORM, nil)
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	return nil
}


