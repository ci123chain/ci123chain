package commands

import (
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/collactor/collactor"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

// keysCmd represents the keys command
func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "manage keys held by the relayer for each chain",
	}

	cmd.AddCommand(keysAddCmd())
	//cmd.AddCommand(keysRestoreCmd())
	//cmd.AddCommand(keysDeleteCmd())
	//cmd.AddCommand(keysListCmd())
	//cmd.AddCommand(keysShowCmd())
	//cmd.AddCommand(keysExportCmd())

	return cmd
}



// keysAddCmd respresents the `keys add` command
func keysAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [chain-id] [[name]]",
		Aliases: []string{"a"},
		Short:   "adds a key to the keychain associated with a particular chain",
		Args:    cobra.RangeArgs(1, 2),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s keys add ibc-0
$ %s keys add ibc-1 key2
$ %s k a ibc-2 testkey`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			var keyName string
			if len(args) == 2 {
				keyName = args[1]
			} else {
				keyName = chain.Key
			}

			//if chain.KeyExists(keyName) {
			//	return errKeyExists(keyName)
			//}
			//
			//coinType, _ := cmd.Flags().GetUint32(flagCoinType)
			privateStr := args[2]
			// Adding key with key add helper
			ko, err := KeyAddOrRestore(chain, keyName, privateStr)
			if err != nil {
				return err
			}

			out, err := json.Marshal(&ko)
			if err != nil {
				return err
			}

			fmt.Println(string(out))
			return nil
		},
	}
	//cmd.Flags().Uint32(flagCoinType, defaultCoinType, "coin type number for HD derivation")

	return cmd
}


// KeyOutput contains mnemonic and address of key
type KeyOutput struct {
	PrivateKey  string `json:"privatekey" yaml:"privatekey"`
	Address  string `json:"address" yaml:"address"`
}


// KeyAddOrRestore is a helper function for add key and restores key when mnemonic is passed
func KeyAddOrRestore(chain *collactor.Chain, keyName string, privatekeys ...string) (KeyOutput, error) {
	var privatekeyStr string
	var err error

	if len(privatekeys) > 0 {
		privatekeyStr = privatekeys[0]
	} else {
		return KeyOutput{}, errors.New("privateKey can not be empty")

	}
	privateKey, err := crypto.HexToECDSA(privatekeyStr)
	if err != nil {
		return KeyOutput{}, errors.Errorf("error format privateKey: %s", privatekeyStr)
	}
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	ko := KeyOutput{PrivateKey: privatekeyStr, Address: sdk.ToAccAddress(address[:]).String()}
	return ko, nil
}

