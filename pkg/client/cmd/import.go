package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPrivate = "private"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().String(flagPrivate, "", "private_key")
	importCmd.Flags().String(flagMnemonic, "", "mnemonic string")
	importCmd.Flags().String(flagHDWPath, "m/44'/60'/0'/0/0", "HD Wallet path")
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "import a account by private_key or mnemonic",
	Long:  `import an encrypted account to the keystore.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		dir := viper.GetString(helper.FlagHomeDir)
		ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
		privateKey := viper.GetString(flagPrivate)
		mnemonic := viper.GetString(flagMnemonic)
		hdPath := viper.GetString(flagHDWPath)
		password, err := helper.GetPasswordFromStd()

		if err != nil {
			return err
		}

		var ac *accounts.Account
		if privateKey != "" {
			ac, err = importPrivateKey(ks, privateKey, password)

		} else if mnemonic != "" && hdPath != "" {
			ac, err = importAccountFromHDW(ks, mnemonic, hdPath, password)
		} else {
			err = fmt.Errorf("private_key and mnemonic is null")
		}

		if err != nil {
			return err
		}

		fmt.Println(ac.Address.Hex())
		return nil
	},
}

func importPrivateKey(ks *keystore.KeyStore, privateKey, password string) (*accounts.Account, error) {
	privPub, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}
	key := crypto.ToECDSAUnsafe(privPub)
	ac, err := ks.ImportECDSA(key, password)
	if err != nil {
		return nil, err
	}
	return &ac, nil
}
