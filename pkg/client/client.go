package client

import (

	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	"os"
)

func NewClientContextFromViper(cdc *codec.Codec) (context.Context, error) {
	nodeURI := viper.GetString(helper.FlagNode)

	var rpc client.Client
	if nodeURI != "" {
		rpc = client.NewHTTP(nodeURI, "/websocket")
	}
	var addr sdk.AccAddress
	addrs := viper.GetString(helper.FlagAddress)
	if len(addrs) > 0 {
		var err error
		addr, err = helper.StrToAddress(addrs)
		if err != nil {
			os.Exit(1)
		}
	}

	//cryptoType := viper.GetInt(helper.FlagWithCrypto)
	return context.Context{
		HomeDir: 	viper.GetString(helper.FlagHomeDir),
		Verbose: 	viper.GetBool(helper.FlagVerbose),
		Height:  	viper.GetInt64(helper.FlagHeight),
		FromAddr: 	addr,
		Blocked:	viper.GetBool(helper.FlagBlocked),
		NodeURI: 	nodeURI,
		Client: 	rpc,
		Cdc: 		cdc,
	}, nil
}