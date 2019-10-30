package cmd

import (
	"github.com/tanhuiya/ci123chain/pkg/client"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
	"github.com/tanhuiya/ci123chain/pkg/util"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)


func init()  {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagTo, "", "Address sending to")
	signCmd.Flags().Uint(flagAmount, 0, "Amount tbe spent")
	signCmd.Flags().Uint(flagGas, 0, "gas for tx")
	signCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	signCmd.Flags().String(flagPassword, "", "passphrase")
	util.CheckRequiredFlag(signCmd, flagAmount)
	util.CheckRequiredFlag(signCmd, flagGas)
}

const isFabric = false

var signCmd = &cobra.Command{
	Use: "sign",
	Short: "Build, Sign transaction",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())
		ctx, err := client.NewClientContextFromViper()
		if err != nil {
			return err
		}
		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return nil
		}
		from := addrs[0]
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return err
		}
		if len(tos) == 0 {
			return errors.New("must provide an address to send to")
		}
		nonce, err := transaction.GetNonceByAddress(from)
		if err != nil {
			return err
		}



		tx := transaction.NewTransferTx(from, tos[0], uint64(viper.GetInt(flagGas)), nonce ,uint64(viper.GetInt(flagAmount)), isFabric)


		password := viper.GetString(flagPassword)
		if len(password) < 1 {
			var err error
			password, err = getPassword()
			if err != nil {
				return err
			}
		}

		txByte, err := getSignedDataWithTx(ctx, tx, password, from)
		if err != nil {
			return err
		}
		fmt.Println(hex.EncodeToString(txByte))
		return nil
	},
}

func getSignedDataWithTx(ctx context.Context, tx transaction.Transaction, password string, from common.Address) ([]byte, error) {
	ks := getDefaultKeystore()
	acc := accounts.Account{
		Address: from,
	}
	acct, err := ks.Find(acc)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s address not found", from.Hex()))
	}
	keyjson, err := ioutil.ReadFile(acct.URL.Path)

	pkey, err := keystore.DecryptKey(keyjson, password)
	privByte := crypto.FromECDSA(pkey.PrivateKey)
	signedtx, err := ctx.SignWithTx(tx, privByte, isFabric)
	return signedtx.Bytes(), nil
}

func getDefaultKeystore() *keystore.KeyStore {
	dir := viper.GetString(helper.FlagHomeDir)
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks
}
