package node

import (
	"fmt"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/p2p"
	"io/ioutil"
)

var cdc = amino.NewCodec()

func init()  {
	cryptoAmino.RegisterAmino(cdc)
}

func LoadNodeKey(filePath string) (*p2p.NodeKey, error) {
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	nodeKey := new(p2p.NodeKey)
	err = cdc.UnmarshalJSON(jsonBytes, nodeKey)
	if err != nil {
		return nil, fmt.Errorf("Error reading Nodekey from %v: %v", filePath, err)
	}
	return nodeKey, nil
}

func GenNodeKeyByPrivKey(filePath string, privKey crypto.PrivKey) (*p2p.NodeKey, error) {
	nodeKey := &p2p.NodeKey{
		PrivKey: privKey,
	}
	jsonBytes, err := cdc.MarshalJSON(nodeKey)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filePath, jsonBytes, 0600)
	if err != nil {
		return nil, err
	}
	return nodeKey, nil
}