package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/client/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"io/ioutil"
	//ac "github.com/ci123chain/ci123chain/pkg/account/keeper"
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
	Short: "Build, Sign transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.BindPFlags(cmd.Flags())

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return types.ErrNewClientCtx(types.DefaultCodespace, err)
		}
		addrs, err := ctx.GetInputAddresses()
		if err != nil {
			return types.ErrGetInputAddrCtx(types.DefaultCodespace, err)
		}
		from := addrs[0]
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return types.ErrParseAddr(types.DefaultCodespace, err)
		}
		if len(tos) == 0 {
			return types.ErrNoAddr(types.DefaultCodespace, err)
		}
		//直接getNonce
		// todo err
		nonce, err := ctx.GetNonceByAddress(from)
		if err != nil {
			return err
		}
		ucoin := uint64(viper.GetInt(flagAmount))

		tx := transfer.NewTransferTx(from, tos[0], uint64(viper.GetInt(flagGas)), nonce , sdk.NewUInt64Coin(ucoin), isFabric)
		password := viper.GetString(flagPassword)
		if len(password) < 1 {
			var err error
			password, err = helper.GetPassphrase(from)
			if err != nil {
				return types.ErrGetPassPhrase(types.DefaultCodespace, err)
			}
		}

		txByte, err := getSignedDataWithTx(ctx, tx, password, from)
		if err != nil {
			return types.ErrGetSignData(types.DefaultCodespace, err)
		}
		fmt.Println(hex.EncodeToString(txByte))
		return nil
	},
}

func getSignedDataWithTx(ctx context.Context, tx transaction.Transaction, password string, from sdk.AccAddress) ([]byte, error) {
	ks := getDefaultKeystore()
	acc := accounts.Account{
		Address: from.Address,
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
