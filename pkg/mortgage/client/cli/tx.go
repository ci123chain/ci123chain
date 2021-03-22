package cli

const (
	FlagToAddress = "to"
	FlagAmount      = "amount"
	FlagUniqueID 	= "uniqueID"
	FlagPassword 	= "password"
)
//
//func init()  {
//
//	MortgageCmd.Flags().String(helper.FlagAddress, "", "Address to sign with")
//	MortgageCmd.Flags().String(FlagToAddress, "", "destnation address")
//	MortgageCmd.Flags().Uint64(FlagAmount, 0, "mortgaged coin")
//	MortgageCmd.Flags().String(FlagUniqueID, "", "mortgaged record uniqueID")
//	MortgageCmd.Flags().String(FlagPassword, "", "passphrase")
//
//	util.CheckRequiredFlag(MortgageCmd, FlagAmount)
//	util.CheckRequiredFlag(MortgageCmd, FlagUniqueID)
//}
//
//var MortgageCmd =  &cobra.Command{
//	Use: "mortgage coin",
//	Short: "mortgage coin from account to module",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		viper.BindPFlags(cmd.Flags())
//		ctx, err := clients.NewClientContextFromViper(types.MortgageCdc)
//		if err != nil {
//			return err
//		}
//		tx, err := BuildCreateMortgageMsg(ctx)
//		if err != nil {
//			return err
//		}
//
//		password := viper.GetString(FlagPassword)
//		if len(password) < 1 {
//			var err error
//			password, err = helper.GetPasswordFromStd()
//			if err != nil {
//				return err
//			}
//		}
//
//		return nil
//	},
//}
//
//func BuildCreateMortgageMsg(ctx context.Context) (*types.MsgMortgage, error) {
//	addrs, err := ctx.GetInputAddresses()
//	if err != nil {
//		return nil, err
//	}
//
//	tos, err := helper.ParseAddrs(viper.GetString(FlagToAddress))
//	if err != nil {
//		return nil, err
//	}
//	if len(tos) == 0 {
//		return nil, errors.New("must provide an address to send to")
//	}
//	ucoin := uint64(viper.GetInt(FlagAmount))
//	uniqueID, err := hex.DecodeString(viper.GetString(FlagUniqueID))
//	if err != nil {
//		return nil, err
//	}
//
//	return &types.MsgMortgage{
//		CommonTx: transaction.CommonTx{
//			From: addrs[0],
//		},
//		ToAddress: tos[0],
//		UniqueID:   uniqueID,
//		Coin: 		sdk.Coin(ucoin),
//	}, nil
//}


