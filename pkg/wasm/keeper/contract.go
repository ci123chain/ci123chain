package keeper

import "encoding/json"

type ContractResult struct {
	_map  map[string]interface{}
	_data interface{}
}

type Response struct {
	Data   []byte    `json:"data"`
}

func (result *ContractResult) Ok() bool {
	if data, exist := result._map["ok"]; exist {
		result._data = data
		return true
	}
	return false
}

func (result *ContractResult) Err() bool {
	if data, exist := result._map["err"]; exist {
		result._data = data
		return true
	}
	return false
}

func (result *ContractResult) Parse() Response {
	var response Response
	b, _ := json.Marshal(result._data)
	_ = json.Unmarshal(b, &response)
	return response
}

func (result *ContractResult) ParseError() string {
	return result._data.(string)
}
