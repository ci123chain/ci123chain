package types

import (
	"encoding/json"
)

const (
	communityContractAddress = "0xfffffffffffffffffffffffffffffffffffffff1"
	//supportRequiredPct = "600000000000000000"  //60*10^16 :60%
	//minAcceptQuorumPct = "\"500000000000000000\""  //50*10^16 :50%
	//openTime = "\"1290600000\""  //2 weeks
	uploadMethod = "upload()"
	InitMethodStr = "init(string)"
	invoker    = communityContractAddress
	initParams = "{\"init_apps\":[{\"app_name\":\"voting_app\",\"init_args\":[\"6000000000\", \"5000000000\", \"1290600000\", \"true\"]}]}"
)

type Contract struct {
	Index    int        `json:"index"`
	Code     string     `json:"code"`
	Method   string     `json:"method"`
	Params   []json.RawMessage   `json:"params"`
	Address  string     `json:"address"`
}


type GenesisState struct {
	Contracts []Contract `json:"contracts"`
	Invoker  string      `json:"invoker"` //invoker.
	Name     string      `json:"name"`
	Version  string      `json:"version"`
	Author   string      `json:"author"`
	Email    string      `json:"email"`
	Describe string      `json:"describe"`
}

func DefaultGenesisState() GenesisState{

	paramByte, _ := json.Marshal(initParams)

	contracts := []Contract{
		{
			Index:   0,
			Code:    aclCode,
			Method:  uploadMethod,
			Params:  []json.RawMessage{},
			Address: "",
		},
		{
			Index:   1,
			Code:    votingCode,
			Method:  uploadMethod,
			Params:  []json.RawMessage{},
			Address: "",
		},
		{
			Index:   2,
			Code:    communityCode,
			Method:  InitMethodStr,
			Params:  []json.RawMessage{json.RawMessage(paramByte)},
			Address: communityContractAddress,
		},

	}

	return GenesisState{
		Contracts:contracts,
		Invoker:invoker,
		Name: "OfficialContract",
		Version: "v0.0.1",
		Author: "Official",
		Email: "ci123chain@corp-ci.com",
		Describe: "OfficialContract",
	}
}