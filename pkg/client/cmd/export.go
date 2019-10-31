package cmd

import (
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

const (
)

func init()  {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().String(flagPassword, "", "passphrase")
	exportCmd.Flags().String(helper.FlagAddress, "", "Address to export")
}


var exportCmd = &cobra.Command{
	Use: "export",
	Short: "export privatekey of a account",
	Long:  `export privatekey of a account from keystore`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		dir := viper.GetString(helper.FlagHomeDir)
		addr := viper.GetString(helper.FlagAddress)
		password := viper.GetString(flagPassword)
		if len(password) < 1 {
			var err error
			password, err = getPassword()
			if err != nil {
				return err
			}
		}
		ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
		acc := accounts.Account{
			Address: common.HexToAddress(addr),
		}
		acct, err := ks.Find(acc)

		keyjson, err := ioutil.ReadFile(acct.URL.Path)
		pkey, err := keystore.DecryptKey(keyjson, password)
		if err != nil {
			return err
		}
		privByte := crypto.FromECDSA(pkey.PrivateKey)

		fmt.Println(hex.EncodeToString(privByte))
		return nil
	},
}
