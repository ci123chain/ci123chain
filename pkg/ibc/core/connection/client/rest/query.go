package rest

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/connection/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router) {
	// Get all validators
	r.HandleFunc("/ibc/connections", queryConnections(cliCtx)).Methods("POST")
	// Get a single validator info
}

func queryConnections(cliCtx context.Context) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		resp, err := utils.QueryConnectionsABCI(cliCtx, 0, 1000)
		if err != nil {
			rest.WriteErrorRes(writer, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error()).Error())
			return
		}
		respQuery := rest.BuildQueryRes("", false, resp, nil)
		rest.PostProcessResponseBare(writer, cliCtx, respQuery)
	}
}
