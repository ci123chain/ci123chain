package rest

import (
	clientcontext "github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/eth/filters"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/net"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/personal"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/api/web3"
	"github.com/ci123chain/ci123chain/pkg/vm/client/rest/keystore"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

// RPC namespaces and API version
const (
	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"

	apiVersion = "1.0"
)

// GetAPIs returns the list of all APIs from the Ethereum namespaces
func GetAPIs(clientCtx clientcontext.Context) []rpc.API {
	ks := getDefaultKeystore()
	ethAPI := eth.NewAPI(clientCtx, ks)
	backend := filters.New(clientCtx)

	return []rpc.API{
		{
			Namespace: Web3Namespace,
			Version:   apiVersion,
			Service:   web3.NewAPI(),
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   ethAPI,
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   filters.NewAPI(clientCtx, backend),
			Public:    true,
		},
		{
			Namespace: PersonalNamespace,
			Version:   apiVersion,
			Service:   personal.NewAPI(ethAPI, ks),
			Public:    false,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   net.NewAPI(clientCtx),
			Public:    true,
		},
	}
}

func getDefaultKeystore() *keystore.KeyStore {
	dir := viper.GetString(helper.FlagHomeDir)
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks
}
