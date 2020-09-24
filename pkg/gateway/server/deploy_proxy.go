package server

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/gateway/dynamic"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"net/http"
)

const TARGETAPPID = "hedlzgp1u48kjf50xtcvwdklminbqe9a"

type DeployProxy struct {
	ProxyType types.ProxyType
	ResponseChannel chan []byte
}

type deployParam struct {
	Type  string 		`json:"type"`
	Value interface{} 	`json:"value"`
}

type networkValue struct {
	Type  string 				   `json:"type"`
	Hosts []map[string]interface{} `json:"hosts"`
}

func NewDeployProxy(pt types.ProxyType) *DeployProxy {
	dp := &DeployProxy{
		ProxyType: pt,
		ResponseChannel:make(chan []byte),
	}
	return dp
}

func (dp *DeployProxy) Handle(r *http.Request, backends []types.Instance, RequestParams map[string]string) []byte {
	params, err := handleDeployParams(RequestParams)
	if err != nil {
		res := dp.ErrorRes(err)
		return res
	}

	res, err := dynamic.CreateChannel(params)
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

	hosts := make(map[string]interface{})
	var env []deployParam
	err = json.Unmarshal([]byte(deployParams["env"]), &env)
	if err != nil {
		return nil, err
	}
	for _, v := range env {
		if v.Type == "environment" {
			environment := v.Value.(map[string]interface{})
			hosts["domain"] = environment["CI_CHAIN_ID"]
			break
		}
	}

	hosts["backend_protocal"] = "HTTP"
	hosts["need_https"] = 0
	hosts["ssl_certificate_data"] = ""
	hosts["ssl_key_data"] = ""
	hosts["target_port"] = 80
	hosts["type"] = "HTTP"

	networksParam := deployParam{
		Type:  "networks",
		Value: networkValue{
			Type: "DOMAIN",
			Hosts: []map[string]interface{}{hosts},
		},
	}

	extraInfo := []deployParam{networksParam}
	for i := 0; i < len(env); i++ {
		if env[i].Type == "environment" || env[i].Type == "volume_mounts" {
			extraInfo = append(extraInfo, env[i])
		}
	}

	params["extra_info"] = extraInfo
	return params, nil
}

func checkDeployParams(deployParams map[string]string) error {
	_, ok := deployParams["instance_name"]
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