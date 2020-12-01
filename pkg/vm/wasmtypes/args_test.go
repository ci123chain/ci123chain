package types

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes/utils"
	"gotest.tools/assert"
	"testing"
)

func TestArgsInts(t *testing.T)  {
	args_str := "{\"method\": \"testUint32(uint32)\", \"args\": [123]}"
	var params utils.CallData
	err := json.Unmarshal([]byte(args_str), &params)
	if err != nil {
		panic(err.Error())
	}
	input, err := ArgsToInput(params)
	sink := NewSink(input)
	a, _ := sink.ReadU32()
	assert.Equal(t, uint32(123), a)
}

