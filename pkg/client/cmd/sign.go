package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/client/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Short: "Build, Sign transfer msg",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			panic(err)
		}

		from := sdk.HexToAddress(viper.GetString(flagFrom))
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return types.ErrParseAddr(types.DefaultCodespace, err)
		}
		if len(tos) == 0 {
			return types.ErrNoAddr(types.DefaultCodespace, err)
		}

		gas := uint64((viper.GetInt(flagGas)))
		amount := uint64(viper.GetInt(flagAmount))
		privKey := viper.GetString(flagKey)
		isFabric := viper.GetBool(flagIsFabric)

		coin := sdk.NewUInt64Coin(amount)
		msg := transfer2.NewMsgTransfer(from, tos[0], coin, isFabric)
		nonce, err := transfer2.GetNonceByAddress(from)
		if err != nil {
			return types.ErrParseParam(types.DefaultCodespace, err)
		}

		txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
		if err != nil {
			return types.ErrParseParam(types.DefaultCodespace, err)
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
