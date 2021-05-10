package main

import (
	"errors"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils"
	"github.com/ci123chain/ci123chain/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagLogLevel 		= "log_level"
	flagCosmosKey 		= "cosmos_key"
	flagEthKey 	 		= "eth_key"
	flagCosmosRpc   	= "cosmos_rpc"
	flagEthRpc 			= "eth_rpc"
	flagFeesDenom   	= "fees_denom"
	flagContractAddr	= "contract_address"
)

func main()  {
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use: 	"orchestrator",
		Short:  "gravity orchestrator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cosmosKey := viper.GetString(flagCosmosKey)
			if len(cosmosKey) == 0 {
				return errors.New("cosmosKey cannot be empty")
			}
			ethKey := viper.GetString(flagEthKey)
			if len(ethKey) == 0 {
				return errors.New("ethKey cannot be empty")
			}
			cosmosRpc := viper.GetString(flagCosmosRpc)
			if len(cosmosRpc) == 0 {
				return errors.New("cosmosRpc cannot be empty")
			}
			ethRpc := viper.GetString(flagEthRpc)
			if len(ethRpc) == 0 {
				return errors.New("ethRpc cannot be empty")
			}
			contractAddr := viper.GetString(flagContractAddr)
			if len(ethRpc) == 0 {
				return errors.New("contractAddr cannot be empty")
			}
			denom := viper.GetString(flagFeesDenom)
			loggerLv := viper.GetString(flagLogLevel)
			logger.GetDefaultLogger(loggerLv)
			gravity_utils.Exec(func() interface{} {
				orchestratorMainLoop(cosmosKey, ethKey, cosmosRpc, ethRpc, denom, contractAddr)
				return nil
			}).Await()
			select {}
		},
	}

	rootCmd.Flags().String(flagLogLevel, "info", "Run abci app with different log level")
	rootCmd.Flags().String(flagCosmosKey, "", "The Cosmos private key of the validator")
	rootCmd.Flags().String(flagEthKey, "", "The Ethereum private key of the validator")
	rootCmd.Flags().String(flagCosmosRpc, "", "The Cosmos RPC url, usually the validator")
	rootCmd.Flags().String(flagEthRpc, "", "The Ethereum RPC url, should be a self hosted node")
	rootCmd.Flags().String(flagFeesDenom, "stake", "The Cosmos Denom in which to pay Cosmos chain fees")
	rootCmd.Flags().String(flagContractAddr, "", "The Ethereum contract address for Gravity, this is temporary")

	rootCmd.Execute()
}

