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
)

//const (
//	flagFrom    = "from"
//	flagTo 		= "to"
//	flagAmount  = "amount"
//	flagGas 	= "gas"
//	flagKey		= "privKey"
//	flagIsFabric= "isFabric"
//	flagDenom   = "denom"
//)

func init()  {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().String(util.FlagTo, "", "Address sending to")
	transferCmd.Flags().Uint(util.FlagAmount, 0, "Amount tbe spent")
	transferCmd.Flags().Uint(util.FlagGas, 0, "gas for tx")
	transferCmd.Flags().String(util.FlagAddress, "", "Address to sign with")
	transferCmd.Flags().String(util.FlagPassword, "", "passphrase")
	transferCmd.Flags().String(util.FlagCoinName, "", "coin denom")

	util.CheckRequiredFlag(transferCmd, util.FlagAmount)
	util.CheckRequiredFlag(transferCmd, util.FlagGas)
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

		res, err := ctx.BroadcastSignedData(txByte)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInternal, fmt.Sprintf("boradcast tx failed:%v", err.Error()))
		}
		_ = ctx.PrintOutput(res)
		return nil
	},
}

