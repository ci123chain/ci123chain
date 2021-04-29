package client

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"

	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/client/http"
	"os"
)

func NewClientContextFromViper(cdc *codec.Codec) (context.Context, error) {
	nodeURI := viper.GetString(util.FlagNode)

	var rpc client.Client
	var err error
	if nodeURI != "" {
		rpc, err = http.New(nodeURI, "/websocket")
		if err != nil {
			os.Exit(1)
		}
	}
	var addr sdk.AccAddress
	addrs := viper.GetString(util.FlagAddress)
	if len(addrs) > 0 {
		var err error
		addr, err = helper.StrToAddress(addrs)
		if err != nil {
			os.Exit(1)
		}
	}

	//cryptoType := viper.GetInt(helper.FlagWithCrypto)
	return context.Context{
		HomeDir: 	viper.GetString(util.FlagHomeDir),
		Verbose: 	viper.GetBool(util.FlagVerbose),
		Height:  	viper.GetInt64(util.FlagHeight),
		FromAddr: 	addr,
		Blocked:	viper.GetBool(util.FlagBlocked),
		NodeURI: 	nodeURI,
		Client: 	rpc,
		Cdc: 		cdc,
	}, nil
}