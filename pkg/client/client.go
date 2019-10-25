package client

import (
	"CI123Chain/pkg/client/context"
	"CI123Chain/pkg/client/helper"
	"CI123Chain/pkg/cryptosuit"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	"os"
)

func NewClientContextFromViper() (context.Context, error) {
	nodeURI := viper.GetString(helper.FlagNode)

	var rpc client.Client
	if nodeURI != "" {
		rpc = client.NewHTTP(nodeURI, "/websocket")
	}
	addrs, err := helper.ParseAddrs(viper.GetString(helper.FlagAddress))
	if err != nil {
		os.Exit(1)
	}
	return context.Context{
		HomeDir: viper.GetString(helper.FlagHomeDir),
		Verbose: viper.GetBool(helper.FlagVerbose),
		Height:  viper.GetInt64(helper.FlagHeight),
		InputAddressed: addrs,
		NodeURI: nodeURI,
		Client: rpc,
		CryptoSuit: cryptosuit.GetSignIdentity(cryptosuit.FabSignType),
	}, nil
}