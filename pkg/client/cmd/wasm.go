package cmd

import (
	"github.com/spf13/cobra"
)

var StoreCodeCmd = &cobra.Command{
	Use: "store [wasm file] --source [source] --builder [builder]",
	Short: "Upload a wasm binary",
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		return nil
	},
}

var InstantiateContractCmd = &cobra.Command{
	Use: "instantiate [code_id_int64] [json_encoded_init_args]",
	Short: "Instantiate a wasm contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		return nil
	},
}

var ExecuteContractCmd = &cobra.Command{
	Use: "execute [contract_addr_bech32] [json_encoded_send_args]",
	Short: "Execute a command on a wasm contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		//
		return nil
	},
}

var WasmCmd = &cobra.Command{
	Use: "wasm",
	Short: "Wasm transaction subcommands",
}

func init() {
	WasmCmd.AddCommand(StoreCodeCmd,
		InstantiateContractCmd,
		ExecuteContractCmd)
	rootCmd.AddCommand(WasmCmd)
}