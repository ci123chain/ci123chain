package rest

import (
	"encoding/json"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	cmn "github.com/tendermint/tendermint/libs/common"
	"net/http"
	"strconv"
)

type Response struct {
	Ret 	uint32 	`json:"ret"`
	Data 	interface{}	`json:"data"`
	Message	string	`json:"message"`
}

// ErrorResponse defines the attributes of a JSON error response.
//type ErrorResponse struct {
//	Code  int    `json:"code,omitempty"`
//	Error string `json:"error"`
//}
//
//// NewErrorResponse creates a new ErrorResponse instance.
//func NewErrorResponse(code int ,err string) ErrorResponse {
//	return ErrorResponse{Code:code, Error:err}
//}
//
//// WriteErrorResponse prepares and writes a HTTP error
//// given a status code and an error message.
//func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(status)
//	_, _ = w.Write(codec.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
//}

func NewErrorRes(err sdk.Error) Response {
	return Response{
		Ret:		uint32(err.Code()),
		Data:		err.Data().(cmn.FmtError).Error(),
		Message:	err.Data().(cmn.FmtError).Format(),
	}
}

func WriteErrorRes(w http.ResponseWriter, err sdk.Error) {
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(NewErrorRes(err))
	_, _ = w.Write(resp)
}

func PostProcessResponseBare(w http.ResponseWriter, ctx context.Context, body interface{}) {
	var res Response
	dataJson, err := json.Marshal(body)
	if err != nil {
		res = Response{
			Ret:     0,
			Data:    string(dataJson),
			Message: "",
		}
	} else {
		res = Response{
			Ret:     0,
			Data:    body,
			Message: "",
		}
	}
	resp, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func ParseQueryHeightOrReturnBadRequest(w http.ResponseWriter, cliCtx context.Context, r *http.Request, heightStr string) (context.Context, bool, sdk.Error) {

	if heightStr != "" {
		height, err := strconv.ParseInt(heightStr, 10, 64)
		if err != nil {
			return cliCtx, false , transfer.ErrCheckParams(sdk.CodespaceRoot, "height error")
		}
		if height < 0 {
			return cliCtx, false , transfer.ErrCheckParams(sdk.CodespaceRoot, "height error")
		}
		if height > 0 {
			cliCtx = cliCtx.WithHeight(height)
		}
	} else {
		cliCtx = cliCtx.WithHeight(0)
	}
	return cliCtx, true , nil
}