package dynamic

import (
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/pretty66/gosdk"
	"github.com/pretty66/gosdk/cienv"
	"net/http"
)

const CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_MULTIPART = "multipart/form-data"
const DPSPACEKEY = "mdocxiqnl43hu68a2lrayv9p5fttm0vf"
const METHOD = "POST"
const API = "Channel/assignChannel"
const ALIAS = "deployment"

var client gosdk.Client
func Init() {
	var err error
	h := http.Header{}
	client, err = gosdk.GetClientInstance(h)
	if err != nil {
		logger.Error("init dynamic error: %v", err)
	}
	err = client.SetAppInfo(cienv.GetEnv("IDG_APPID"), cienv.GetEnv("IDG_APPKEY"), cienv.GetEnv("IDG_CHANNEL"), cienv.GetEnv("IDG_VERSION"))
	if err != nil {
		logger.Error("init dynamic error: %v", err)
	}
}

func CreateChannel(params map[string]interface{}) ([]byte, error) {
	res, err := client.Call(DPSPACEKEY, METHOD, API, params, ALIAS, CONTENT_TYPE_FORM, nil)
	if err != nil {
		return nil, err
	}
	return res, nil
}


