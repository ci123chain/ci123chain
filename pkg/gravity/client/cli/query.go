package cli

import (
	"github.com/spf13/cobra"
)

func QueryObserved() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "observed",
		Short:                      "observed ETH events",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastObservedNonceRequest(storeKey, cdc),
		// CmdGetLastObservedNoncesRequest(storeKey, cdc),
		// CmdGetLastObservedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetAllBridgedDenominatorsRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}
func QueryApproved() *cobra.Command {
	testingTxCmd := &cobra.Command{
		Use:                        "approved",
		Short:                      "approved cosmos operation",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
	}
	testingTxCmd.AddCommand([]*cobra.Command{
		// CmdGetLastApprovedNoncesRequest(storeKey, cdc),
		// CmdGetLastApprovedMultiSigUpdateRequest(storeKey, cdc),
		// CmdGetInflightBatchesRequest(storeKey, cdc),
	}...)

	return testingTxCmd
}