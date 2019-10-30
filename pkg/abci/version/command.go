package version

import (
	"fmt"

	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/consts"
	"github.com/spf13/cobra"
)

var (
	// VersionCmd prints out the current sdk version
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		Run:   printVersion,
	}
)

// return version of CLI/node and commit hash
func GetVersion() string {
	return consts.Version
}

// CMD
func printVersion(cmd *cobra.Command, args []string) {
	fmt.Println(GetVersion())
}
