package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
)


func init()  {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(flagTo, "", "Address sending to")
	signCmd.Flags().Uint(flagAmount, 0, "Amount tbe spent")
	signCmd.Flags().Uint(flagGas, 0, "gas for tx")
	signCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	signCmd.Flags().String(flagPassword, "", "passphrase")
	signCmd.Flags().String(flagDenom, "", "coin denom")
	util.CheckRequiredFlag(signCmd, flagAmount)
	util.CheckRequiredFlag(signCmd, flagGas)
}

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
			return errors.New("invalid to addresss")
		}
		if len(tos) == 0 {
			return errors.New("invalid to address")
		}
		d := viper.GetString(flagDenom)
		//if d == "" {
		//	return errors.New("invalid denom")
		//}

		gas := uint64((viper.GetInt(flagGas)))
		//amount := uint64(viper.GetInt(flagAmount))
		amount, ok := new(big.Int).SetString(viper.GetString(flagAmount), 10)
		if !ok {
			return errors.New("invalid amount")
		}
		privKey := viper.GetString(flagKey)

		//coin := sdk.NewUInt64Coin(d, amount)
		coin := sdk.NewCoin(d, sdk.NewIntFromBigInt(amount))
		msg := transfer2.NewMsgTransfer(from, tos[0], sdk.NewCoins(coin))
		nonce, err := transfer2.GetNonceByAddress(from)
		if err != nil {
			return errors.New("invalid nonce")
		}
		txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
		if err != nil {
			return errors.New(fmt.Sprintf("sign tx failed: %v", err.Error()))
		}

		fmt.Println(hex.EncodeToString(txByte))
		return nil
	},
}


func getDefaultKeystore() *keystore.KeyStore {
	dir := viper.GetString(helper.FlagHomeDir)
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks
}
