package server

import (
	"encoding/json"
	"errors"
	"github.com/ci123chain/ci123chain/pkg/gateway/dynamic"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"net/http"
)

const TARGETAPPID = "hedlzgp1u48kjf50xtcvwdklminbqe9a"

type DeployProxy struct {
	ProxyType types.ProxyType
	ResponseChannel chan []byte
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
	logger.Debug("===\n deploy params: %v", params)
	res, err := dynamic.CreateChannel(r, params)
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
		Ret: 0,
		Message:  err.Error(),
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