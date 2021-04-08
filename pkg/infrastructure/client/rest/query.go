package rest

import (
	"encoding/json"
	"fmt"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	//sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/rest"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterQueryRoutes(cliCtx context.Context, r *mux.Router) {
	r.HandleFunc("/infrastructure/query_content", queryStoredContentFn(cliCtx)).Methods("POST")
}

func queryStoredContentFn(cliCtx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyStr := r.FormValue("key")
		if keyStr == "" {
			rest.WriteErrorRes(w, sdkerrors.Wrap(sdkerrors.ErrParams, "invalid key").Error())
			return
		}
		b := cliCtx.Cdc.MustMarshalJSON(types.NewContentParams([]byte(keyStr)))

		route := fmt.Sprintf("/custom/%s/%s", types.QuerierRoute, types.QueryContent)
		res, _, _, err := cliCtx.Query(route, b, false)

		if err != nil {
			rest.WriteErrorRes(w, err.Error())
			return
		}
		var result types.StoredContent
		_ = json.Unmarshal(res, &result)
		rest.PostProcessResponseBare(w, cliCtx, result)
	}
}
