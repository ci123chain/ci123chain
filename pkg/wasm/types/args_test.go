package types

import (
	"encoding/json"
	assert "github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestArgsInts(t *testing.T)  {
	args_str := "{\"method\": \"get\", \"args\": {\"uint32\": 23454, \"int64\": 9223372036854775807}}"
	var params CallContractParam
	err := json.Unmarshal([]byte(args_str), &params)
	if err != nil {
		panic(err.Error())
	}
	res := Serialize(params.Args)

	sink := NewSink(res)
	a, _ := sink.ReadU32()
	assert.Equal(t,uint32(23454), a)
	i64 , _ := sink.ReadI64()
	assert.Equal(t,int64(9223372036854775807), i64)
}

func TestArgsInt128(t *testing.T)  {
	args_str := "{\"method\": \"get\", \"args\": {\"int128\": 922337203685477580911777}}"
	var params CallContractParam
	err := json.Unmarshal([]byte(args_str), &params)
	if err != nil {
		panic(err.Error())
	}
	res := Serialize(params.Args)
	sink := NewSink(res)
	i128 , _ := sink.ReadI128()
	b := new(big.Int).SetBytes(ToBigEnd(i128.Bytes()))
	assert.Equal(t,"922337203685477580911777", b.String())
}

func TestArgsString(t *testing.T)  {
	args_str := "{\"method\": \"get\", \"args\": {\"string\": \"123131fad--fafas\"}}"
	var params CallContractParam
	err := json.Unmarshal([]byte(args_str), &params)
	if err != nil {
		panic(err.Error())
	}
	res := Serialize(params.Args)
	sink := NewSink(res)
	str , _ := sink.ReadString()
	assert.Equal(t,"123131fad--fafas", str)
}

func TestArgsBool(t *testing.T)  {
	args_str := "{\"method\": \"get\", \"args\": {\"bool\": true, \"bool\": false}}"
	var params CallContractParam
	err := json.Unmarshal([]byte(args_str), &params)
	if err != nil {
		panic(err.Error())
	}
	res := Serialize(params.Args)
	sink := NewSink(res)
	a , _ := sink.ReadBool()
	assert.Equal(t,true, a)

	b , _ := sink.ReadBool()
	assert.Equal(t,false, b)
}