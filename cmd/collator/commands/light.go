package commands

import (
	"encoding/json"
	"fmt"
	helpers "github.com/ci123chain/ci123chain/pkg/collactor/helper"
	tmclient "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

// chainCmd represents the keys command
func lightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "light",
		Aliases: []string{"l"},
		Short:   "manage light clients held by the relayer for each chain",
	}

	//cmd.AddCommand(lightHeaderCmd())
	cmd.AddCommand(initLightCmd())
	//cmd.AddCommand(updateLightCmd())
	//cmd.AddCommand(deleteLightCmd())

	return cmd
}


func initLightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [chain-id]",
		Aliases: []string{"i"},
		Short:   "Initiate the light client",
		Long: `Initiate the light client by:
	1. passing it a root of trust as a --hash/-x and --height
	2. Use --force/-f to initialize from the configured node`,
		Args: cobra.ExactArgs(1),
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s light init ibc-0 --force
$ %s light init ibc-1 --height 1406 --hash <hash>
$ %s l i ibc-2 --force`, appName, appName, appName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			chain, err := config.Chains.Get(args[0])
			if err != nil {
				return err
			}

			force, err := cmd.Flags().GetBool(flagForce)
			if err != nil {
				return err
			}
			//height, err := cmd.Flags().GetInt64(flags.FlagHeight)
			//if err != nil {
			//	return err
			//}
			//hash, err := cmd.Flags().GetBytesHex(flagHash)
			//if err != nil {
			//	return err
			//}

			out, err := helpers.InitLight(chain, force)
			if err != nil {
				return err
			}

			if out != "" {
				fmt.Println(out)
			}

			return nil
		},
	}
	return forceFlag(cmd)
	//return forceFlag(lightFlags(cmd))

}


// API Handlers

// GetLightHeader handles the route
func GetLightHeader(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain, err := config.Chains.Get(vars["chain-id"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	var header *tmclient.Header
	height := strings.TrimSpace(r.URL.Query().Get("height"))

	if len(height) == 0 {
		header, err = helpers.GetLightHeader(chain)
		if err != nil {
			helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
			return
		}
	} else {
		header, err = helpers.GetLightHeader(chain, height)
		if err != nil {
			helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
			return
		}
	}
	helpers.SuccessJSONResponse(http.StatusOK, header, w)
}



// GetLightHeight handles the route
func GetLightHeight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain, err := config.Chains.Get(vars["chain-id"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	height, err := chain.GetLatestLightHeight()
	if err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	helpers.SuccessJSONResponse(http.StatusOK, height, w)
}


type postLightRequest struct {
	Force  bool   `json:"force"`
	//Height int64  `json:"height"`
	//Hash   string `json:"hash"`
}

// PostLight handles the route
func PostLight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain, err := config.Chains.Get(vars["chain-id"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	var request postLightRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	out, err := helpers.InitLight(chain, request.Force)
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	if out == "" {
		out = fmt.Sprintf("successfully created light client for %s", vars["chain-id"])
	}
	helpers.SuccessJSONResponse(http.StatusCreated, out, w)
}


// DeleteLight handles the route
func DeleteLight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain, err := config.Chains.Get(vars["chain-id"])
	if err != nil {
		helpers.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	err = chain.DeleteLightDB()
	if err != nil {
		helpers.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}

	helpers.SuccessJSONResponse(http.StatusOK, fmt.Sprintf("Removed Light DB for %s", vars["chain-id"]), w)
}
