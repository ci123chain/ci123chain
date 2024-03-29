package cmd

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	transfer2 "github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"math/big"
)

const (
	flagFrom    = "from"
	flagTo 		= "to"
	flagAmount  = "amount"
	flagGas 	= "gas"
	flagKey		= "privKey"
	flagIsFabric= "isFabric"
	flagDenom   = "denom"
)

func init()  {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().String(flagTo, "", "Address sending to")
	transferCmd.Flags().Uint(flagAmount, 0, "Amount tbe spent")
	transferCmd.Flags().Uint(flagGas, 0, "gas for tx")
	transferCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
	transferCmd.Flags().String(flagPassword, "", "passphrase")
	transferCmd.Flags().String(flagDenom, "", "coin denom")

	util.CheckRequiredFlag(transferCmd, flagAmount)
	util.CheckRequiredFlag(transferCmd, flagGas)
}

var transferCmd = &cobra.Command{
	Use: "transfer",
	Short: "Build, Sign, and send transfer",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := viper.BindPFlags(cmd.Flags())
		if err != nil {
			panic(err)
		}

		ctx, err := client.NewClientContextFromViper(cdc)
		if err != nil {
			return err
		}
		from := sdk.HexToAddress(viper.GetString(flagFrom))
		tos, err := helper.ParseAddrs(viper.GetString(flagTo))
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid to address")
		}
		if len(tos) == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid to address")
		}
		d := viper.GetString(flagDenom)
		//if d == "" {
		//	return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid denom")
		//}

		gas := uint64((viper.GetInt(flagGas)))
		//amount := uint64(viper.GetInt(flagAmount))
		amount, ok := new(big.Int).SetString(viper.GetString(flagAmount), 10)
		if !ok {
			return sdkerrors.Wrap(sdkerrors.ErrParams, fmt.Sprintf("invalid amount %s", viper.GetString(flagAmount)))
		}
		privKey := viper.GetString(flagKey)

		//coin := sdk.NewUInt64Coin(d, amount)
		coin := sdk.NewCoin(d, sdk.NewIntFromBigInt(amount))
		msg := transfer2.NewMsgTransfer(from, tos[0], sdk.NewCoins(coin))
		nonce, err := transfer2.GetNonceByAddress(from)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrParams, "invalid nonce")
		}

		txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, privKey, cdc)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("sign tx failed: %v", err.Error()))
		}

		res, err := ctx.BroadcastSignedData(txByte)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("boradcast tx failed:%v", err.Error()))
		}
		_ = ctx.PrintOutput(res)
		return nil
	},
}

