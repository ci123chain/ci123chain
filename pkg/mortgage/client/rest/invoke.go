package rest

import (
	"github.com/gorilla/mux"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"net/http"
)

func RegisterRoutes(cliCtx context.Context, r *mux.Router)  {
	r.HandleFunc("/mortgage/mortgaged", postMortgagedHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/mortgage/done", postMortgageDoneHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/mortgage/cancel", postMortgageCancelHandlerFn(cliCtx)).Methods("POST")
}

func postMortgagedHandlerFn(cliCtx context.Context) http.HandlerFunc {
	return nil
}

func postMortgageDoneHandlerFn(cliCtx context.Context) http.HandlerFunc  {
	return nil
}

func postMortgageCancelHandlerFn(cliCtx context.Context) http.HandlerFunc  {
	return nil
}