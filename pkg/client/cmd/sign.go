package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)


func init()  {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().String(util.FlagTo, "", "Address sending to")
	signCmd.Flags().Uint(util.FlagAmount, 0, "Amount tbe spent")
	signCmd.Flags().Uint(util.FlagGas, 0, "gas for tx")
	signCmd.Flags().String(util.FlagAddress, "", "Address to sign with")
	signCmd.Flags().String(util.FlagPassword, "", "passphrase")
	signCmd.Flags().String(util.FlagCoinName, "", "coin denom")
	util.CheckRequiredFlag(signCmd, util.FlagAmount)
	util.CheckRequiredFlag(signCmd, util.FlagGas)
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

		from := sdk.HexToAddress(viper.GetString(util.FlagFrom))
		tos, err := helper.ParseAddrs(viper.GetString(util.FlagTo))
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid to address")
		}
		if len(tos) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid to address")
		}
		d := viper.GetString(util.FlagCoinName)
		if d == "" {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid denom")
		}

		gas := uint64((viper.GetInt(util.FlagGas)))
		amount := uint64(viper.GetInt(util.FlagAmount))
		privKey := viper.GetString(util.FlagKey)
		isFabric := viper.GetBool(util.FlagIsFabric)

		coin := sdk.NewUInt64Coin(d, amount)
		msg := transfer2.NewMsgTransfer(from, tos[0], sdk.NewCoins(coin), isFabric)
		nonce, err := transfer2.GetNonceByAddress(from)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid nonce")
		}
		txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("sign tx failed: %v", err.Error()))
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
	dir := viper.GetString(util.FlagHomeDir)
	ks := keystore.NewKeyStore(dir, keystore.StandardScryptN, keystore.StandardScryptP)
	return ks
}
