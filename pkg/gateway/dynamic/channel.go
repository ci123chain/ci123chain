package dynamic

import (
	"github.com/pretty66/gosdk"
	"github.com/pretty66/gosdk/cienv"
	"net/http"
)

const APP_ID  = "cc2535f30a87457b91649ecbed0245e6"

func init()  {
	err := cienv.SetEnv(gosdk.GATEWAY_SERVICE_KEY, "kong:http://127.0.0.1:13800")
	if err != nil {
		panic(err)
	}
}

// 动态添加channel, 向部署空间请求, 利用 go-sdk
func createChannel(appid, appkey, channel, version string) error {
	var _header http.Header
	// 获取对象，head是请求的HEAD字段，用来解析HEAD中的Authorization中的token
	client, err := gosdk.GetClientInstance(_header)
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
	client = client.SetServices(services)
}
