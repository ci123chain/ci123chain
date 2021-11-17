package cmd

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/node"
)

func StartRestServer(cdc *codec.Codec, tmNode *node.Node, addr string) error {
	rs := NewRestServer(cdc, tmNode)
	err := rs.Start(
		addr,
		viper.GetInt(FlagMaxOpenConnections),
		uint(viper.GetInt(FlagRPCReadTimeout)),
		uint(viper.GetInt(FlagRPCWriteTimeout)),
	)
	return err
}
