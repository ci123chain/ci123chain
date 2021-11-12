package util

import (
	"errors"
	"github.com/ci123chain/ci123chain/pkg/libs"
	"gitlab.oneitfarm.com/bifrost/sesdk/discovery"
	"os"
	"strings"
	"time"
)

type DomainInfo struct {
	Host26657   string  `json:"host_26657"`
	Host8546    string  `json:"host_8546"`
}


func GetDomain() (host string, err error) {
	_, _ = libs.RetryI(15, func(retryTimes int) (interface{}, error) {
		host, err = Discovery()
		if err != nil {
		}
		return host, err
	})
	return
}

func Discovery() (string, error) {

	appID := os.Getenv("CI_UNIQUE_KEY")
	region := os.Getenv("IDG_SITEUID")
	env := os.Getenv("CI_SE_ENV")
	zone := os.Getenv("IDG_CLUSTERUID")
	address := "chain-discovery-service-eye.chain-discovery:7171"
	if appID == "" {
		return "", errors.New("appID is empty")
	}
	hn := os.Getenv("PODNAME")
	// 注册中心自身，初始化配置
	conf := &discovery.Config{
		// discovery地址
		Nodes: []string{address},
		Region:   region,
		Zone:     zone,
		Env:      env,
		Host:     hn,               // hostname
		RenewGap: time.Second * 30, // 心跳时间
	}

	// 实例化discovery对象
	dis, err := discovery.New(conf)
	if err != nil {
		return "", err
	}
	// 服务发现：目标服务的唯一识别号
	ep, err := dis.GetEndpoint(appID)
	if err != nil {
		return "", err
	}
	res := strings.Split(ep.Host, "//")
	return strings.Split(res[1], ":")[0], nil
}


//func Discovery(f func(err error)) string {
//
//	appID := os.Getenv("CI_VALIDATOR_KEY")
//	address := "192.168.60.48:80"  //os.Getenv("MSP_SE_NGINX_ADDRESS")
//	region := "sal2" //os.Getenv("IDG_SITEUID")
//	env := "production"   //os.Getenv("MSP_SE_ENV")
//	zone := "aliyun-sh-prod" //os.Getenv("IDG_CLUSTERUID")
//	if appID == "" {
//		f(errors.New("appID is empty"))
//	}
//	hn := os.Getenv("PODNAME")
//	//hn, _ := os.Hostname()
//	// 注册中心自身，初始化配置
//	conf := &discovery.Config{
//		// discovery地址
//		Nodes: []string{address},
//		// Nodes:    []string{"127.0.0.1:7171"},
//		Region:   region,
//		Zone:     zone,
//		Env:      env,
//		Host:     hn,               // hostname
//		RenewGap: time.Second * 30, // 心跳时间
//	}
//
//	// 实例化discovery对象
//	dis, err := discovery.New(conf)
//	if err != nil {
//		f(err)
//	}
//	// 服务发现：目标服务的唯一识别号
//	ep, err := dis.GetEndpoint(appID)
//	if err != nil {
//		f(err)
//	}
//	res := strings.Split(ep.Host, "//")
//	return strings.Split(res[1], ":")[0]
//}
