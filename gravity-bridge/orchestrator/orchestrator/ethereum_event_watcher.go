package main

import (
	"crypto/ecdsa"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/umbracle/go-web3/jsonrpc"
)

func CheckForEvents(client *jsonrpc.Client,
	contact cosmos_gravity.Contact,
	contractAddr string,
	cosmosPrivKey *ecdsa.PrivateKey,
	fee sdk.Coin,
	startingBlock uint64) (uint64, error) {

	return 0, nil
}

