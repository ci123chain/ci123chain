package server

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"github.com/pretty66/gosdk"
	"github.com/pretty66/gosdk/cienv"
	"net/http"
)

const GATEWAY_SERVICE_KEY = "GATEWAY_HOST_SERVICE"
const GATEWAY_URL = "kong:http://127.0.0.1:13800"
const CONTENT_TYPE_FORM = "application/x-www-form-urlencoded"
const CONTENT_TYPE_JSON = "application/json"
const CONTENT_TYPE_MULTIPART = "multipart/form-data"
const TARGETAPPID = "hedlzgp1u48kjf50xtcvwdklminbqe9a"
const DPSPACEKEY = "mdocxiqnl43hu68a2lrayv9p5fttm0vf"
const METHOD = "POST"
const API = "Channel/assignChannel"
const ALIAS = "deployment"

type DeployProxy struct {
	ProxyType types.ProxyType
	ResponseChannel chan []byte
}

func NewDeployProxy(pt types.ProxyType) *DeployProxy {
	dp := &DeployProxy{
		ProxyType: pt,
		ResponseChannel:make(chan []byte),
	}
	err := cienv.SetEnv(GATEWAY_SERVICE_KEY, GATEWAY_URL)
	if err != nil {
		panic(err)
	}
	return dp
}

func (dp *DeployProxy) Handle(r *http.Request, backends []types.Instance, RequestParams map[string]string) []byte {
	params, err := handleDeployParams(RequestParams)
	if err != nil {
		res := dp.ErrorRes(err)
		return res
	}


	client, err := gosdk.GetClientInstance(r.Header)
	if err != nil {
		res := dp.ErrorRes(err)
		return res
	}

	appkey := RequestParams["appkey"]
	channel := RequestParams["channel"]
	version := RequestParams["version"]

	err = client.SetAppInfo(TARGETAPPID, appkey, channel, version)
	if err != nil {
		res := dp.ErrorRes(err)
		return res
	}

	res, err := client.Call(DPSPACEKEY, METHOD, API, params, ALIAS, CONTENT_TYPE_FORM, nil)
	if err != nil {
		res := dp.ErrorRes(err)
		return res
	}

	dp.ResponseChannel <- res
	return res
}

func (dp *DeployProxy) Response() *chan []byte {
	return &dp.ResponseChannel
}

func (dp *DeployProxy) ErrorRes(err error) []byte {
	res, _ := json.Marshal(types.ErrorResponse{
		Err:  err.Error(),
	})
	dp.ResponseChannel <- res
	return res
}

func handleDeployParams(deployParams map[string]string) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	err := checkDeployParams(deployParams)
	if err != nil {
		return nil, err
	}

	params["target_appid"] = TARGETAPPID
	params["instance_name"] = deployParams["instance_name"]
	params["instance_properties"] = deployParams["instance_properties"]
	params["from_id"] = deployParams["idempotent_id"]
	params["from_type"] = deployParams["idempotent_type"]
	params["extra_info"] = deployParams["env"]

	return params, nil
}

func checkDeployParams(deployParams map[string]string) error {
	_, ok := deployParams["appkey"]
	if !ok {
		return errors.New("missing instance_name")
	}
	_, ok = deployParams["channel"]
	if !ok {
		return errors.New("missing instance_name")
	}
	_, ok = deployParams["instance_name"]
	if !ok {
		return errors.New("missing instance_name")
	}
	_, ok = deployParams["instance_properties"]
	if !ok {
		return errors.New("missing instance_properties")
	}
	_, ok = deployParams["idempotent_id"]
	if !ok {
		return errors.New("missing idempotent_id")
	}
	_, ok = deployParams["idempotent_type"]
	if !ok {
		return errors.New("missing idempotent_type")
	}
	_, ok = deployParams["env"]
	if !ok {
		return errors.New("missing env")
	}

	return nil
}