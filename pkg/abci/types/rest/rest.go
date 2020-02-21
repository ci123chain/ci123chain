package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/context"
	"github.com/tanhuiya/ci123chain/pkg/transfer"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	"net/http"
	"net/url"
	"strconv"
)
const (
	DefaultPage  = 1
	DefaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
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


// ParseHTTPArgsWithLimit parses the request's URL and returns a slice containing
// all arguments pairs. It separates page and limit used for pagination where a
// default limit can be provided.
func ParseHTTPArgsWithLimit(r *http.Request, defaultLimit int) (tags []string, page, limit int, err error) {
	tags = make([]string, 0, len(r.Form))
	for key, values := range r.Form {
		if key == "page" || key == "limit" {
			continue
		}
		var value string
		value, err = url.QueryUnescape(values[0])
		if err != nil {
			return tags, page, limit, err
		}

		var tag string
		if key == types.TxHeightKey {
			tag = fmt.Sprintf("%s=%s", key, value)
		} else {
			tag = fmt.Sprintf("%s='%s'", key, value)
		}
		tags = append(tags, tag)
	}

	pageStr := r.FormValue("page")
	if pageStr == "" {
		page = DefaultPage
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return tags, page, limit, err
		} else if page <= 0 {
			return tags, page, limit, errors.New("page must greater than 0")
		}
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return tags, page, limit, err
		} else if limit <= 0 {
			return tags, page, limit, errors.New("limit must greater than 0")
		}
	}

	return tags, page, limit, nil
}