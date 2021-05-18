package commands

import (
	"github.com/spf13/cobra"

)

// keysCmd represents the keys command
func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "manage keys held by the relayer for each chain",
	}

	//cmd.AddCommand(keysAddCmd())
	//cmd.AddCommand(keysRestoreCmd())
	//cmd.AddCommand(keysDeleteCmd())
	//cmd.AddCommand(keysListCmd())
	//cmd.AddCommand(keysShowCmd())
	//cmd.AddCommand(keysExportCmd())

	return cmd
}



// keysAddCmd respresents the `keys add` command
//func keysAddCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:     "add [chain-id]",
//		Aliases: []string{"a"},
//		Short:   "adds a key to the keychain associated with a particular chain",
//		Args:    cobra.RangeArgs(1, 2),
//		Example: strings.TrimSpace(fmt.Sprintf(`
//$ %s keys add ibc-0
//$ %s keys add ibc-1 key2
//$ %s k a ibc-2 testkey`, appName, appName, appName)),
//		RunE: func(cmd *cobra.Command, args []string) error {
//			//chain, err := config.Chains.Get(args[0])
//			//if err != nil {
//			//	return err
//			//}
//
//			//if chain.KeyExists(keyName) {
//			//	return errKeyExists(keyName)
//			//}
//			//
//			//coinType, _ := cmd.Flags().GetUint32(flagCoinType)
//			privateStr := args[1]
//			// Adding key with key add helper
//			ko, err := KeyAddOrRestore(privateStr)
//			if err != nil {
//				return err
//			}
//
//			home, err := cmd.Flags().GetString(flagHome)
//			if err != nil {
//				return err
//			}
//			err = writeKey(ko, home, args[0])
//			return err
//		},
//	}
//	//cmd.Flags().Uint32(flagCoinType, defaultCoinType, "coin type number for HD derivation")
//
//	return cmd
//}
//
//func writeKey(output collactor.KeyOutput, home, chainid string ) error {
//	out, err := json.Marshal(&output)
//	if err != nil {
//		return err
//	}
//
//	keyDir := path.Join(home, "keys")
//	keyPath := path.Join(keyDir, chainid + ".json")
//	// Then create the file...
//
//	if _, err := os.Stat(keyDir); os.IsNotExist(err) {
//		// Create the home folder
//		if err = os.Mkdir(keyDir, os.ModePerm); err != nil {
//			return err
//		}
//	}
//
//	f, err := os.Create(keyPath)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	// And write the default configs to that location...
//	if _, err = f.Write(out); err != nil {
//		return err
//	}
//	// And return no error...
//	return nil
//
//}


