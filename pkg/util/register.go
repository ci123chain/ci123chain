package util

import (
	ctx "context"
	"encoding/json"
	"errors"
	"github.com/tendermint/tendermint/libs/log"
	"gitlab.oneitfarm.com/bifrost/sesdk"
	"gitlab.oneitfarm.com/bifrost/sesdk/discovery"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)


func SetupRegisterCenter(f func(err error, lg log.Logger)) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "register-center")
	appID := os.Getenv("CI_VALIDATOR_KEY")
	address := "192.168.60.48:80"  //os.Getenv("MSP_SE_NGINX_ADDRESS")
	region := "sal2" //os.Getenv("IDG_SITEUID")
	env := "production"   //os.Getenv("MSP_SE_ENV")
	zone := "aliyun-sh-prod" //os.Getenv("IDG_CLUSTERUID")
	if appID == "" {
		logger.Error("CI_VALIDATOR_KEY can not be empty")
		os.Exit(1)
	}
	hn := os.Getenv("PODNAME")
	serviceName := os.Getenv("IDG_SERVICE_NAME")
	weight := os.Getenv("IDG_WEIGHT")
	rt := os.Getenv("IDG_RUNTIME")
	// 注册中心自身，初始化配置
	conf := &discovery.Config{
		// discovery地址
		Nodes:    []string{address},
		Region:   region,
		Zone:     zone,
		Env:      env,
		Host:     hn,               // hostname
		RenewGap: time.Second * 30, // 心跳时间
	}
	// 自身实例信息
	ins := &sesdk.Instance{
		Region:   region,
		Zone:     zone,
		Env:      env,
		AppID:    appID, // 自身唯一识别号
		Hostname: hn,
		Addrs: []string{ // 可上报任意服务监听地址，供发现方连接
			"http://127.0.0.1:80",
			//"https://127.0.0.1:443",
			//"tcp://192.168.2.88:3030",
		},
		// 上报任意自身属性信息
		Metadata: map[string]string{
			"weight":       weight, // 负载均衡权重
			"runtime":      rt,
			"service_name": serviceName,
		},
	}
	// 实例化discovery对象
	dis, err := discovery.New(conf)
	if err != nil {
		f(err, logger)
	}
	// 注册自身
	_, err = dis.Register(ctx.Background(), ins)
	if err != nil {
		f(err, logger)
	}
	// 监听系统信号，服务下线
	dis.ExitSignal(func(s os.Signal) {
		logger.Info("got exit signal, exit now", "signal", s.String())
	})
}


func Discovery(f func(err error)) string {

	appID := os.Getenv("CI_VALIDATOR_KEY")
	address := "192.168.60.48:80"  //os.Getenv("MSP_SE_NGINX_ADDRESS")
	region := "sal2" //os.Getenv("IDG_SITEUID")
	env := "production"   //os.Getenv("MSP_SE_ENV")
	zone := "aliyun-sh-prod" //os.Getenv("IDG_CLUSTERUID")
	if appID == "" {
		f(errors.New("appID is empty"))
	}
	hn := os.Getenv("PODNAME")
	// 注册中心自身，初始化配置
	conf := &discovery.Config{
		// discovery地址
		Nodes:    []string{address},
		Region:   region,
		Zone:     zone,
		Env:      env,
		Host:     hn,               // hostname
		RenewGap: time.Second * 30, // 心跳时间
	}

	// 实例化discovery对象
	dis, err := discovery.New(conf)
	if err != nil {
		f(err)
	}
	// 服务发现：目标服务的唯一识别号
	ep, err := dis.GetEndpoint(appID)
	if err != nil {
		f(err)
	}
	// 根据发现的地址测试调用
	resp, err := http.Get(ep.Host + "/info")
	if err != nil {
		f(err)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		f(err)
	}
	var result map[string]string
	err = json.Unmarshal(b, &result)
	if err != nil {
		f(err)
	}
	if len(result) == 0 {
		f(errors.New("empty result from remote discovery"))
	}
	return result["host"]
}
