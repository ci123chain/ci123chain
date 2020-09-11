package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/client"
	"github.com/ci123chain/ci123chain/pkg/client/context"
	"github.com/ci123chain/ci123chain/pkg/client/helper"
	"github.com/ci123chain/ci123chain/pkg/transfer"
	"github.com/ci123chain/ci123chain/pkg/util"
	"github.com/tendermint/tendermint/crypto/merkle"
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
	Ret 	int64 	`json:"ret"`
	Data    json.RawMessage `json:"data"`
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

type QueryRes struct {
	Value interface{}   `json:"value"`
	Height int64		`json:"height,omitempty"`
	Proof *merkle.Proof	`json:"proof,omitempty"`
}

func BuildQueryRes(height string, isProve bool, value interface{}, proof *merkle.Proof) QueryRes {
	resp := QueryRes{
		Value: value,
	}
	if height != "" {
		queryHeight, _ := strconv.ParseInt(height, 10, 64)
		resp.Height = queryHeight
	}
	if isProve {
		resp.Proof = proof
	}
	return resp
}

func CheckHeightAndProve(w http.ResponseWriter, height, prove string, codespace sdk.CodespaceType) (isValid bool) {
	if height != "" {
		_, Err := util.CheckInt64(height)
		if Err != nil {
			WriteErrorRes(w, client.ErrParseParam(codespace, Err))
			isValid = false
			return
		}
	}
	if prove != "" && prove != "true" && prove != "false"{
		WriteErrorRes(w, client.ErrParseParam(codespace, errors.New("prove need true or false")))
		isValid = false
		return
	}
	return true
}

func NewErrorRes(err sdk.Error) Response {
	buildData := struct {
		Code sdk.CodeType `json:"code"`
		CodeSpace sdk.CodespaceType `json:"code_space"`
	}{
		err.Code(),
		err.Codespace(),
	}
	data, _ := json.Marshal(buildData)
	return Response{
		Ret:		-1,
		Data:       data,
		Message:	err.ABCILog(),
	}
}

func WriteErrorRes(w http.ResponseWriter, err sdk.Error) {
	w.Header().Set("Content-Type", "application/json")
	nerr := NewErrorRes(err)
	resp, _ := json.Marshal(nerr)
	_, _ = w.Write(resp)
}

func PostProcessResponseBare(w http.ResponseWriter, ctx context.Context, body interface{}) {
	var res Response
	dataJson, _ := json.Marshal(body)
	switch body.(type) {
	case sdk.TxResponse:
		b := body.(sdk.TxResponse)
		if b.Code == 0 {
			res = Response{
				Ret:     1,
				Data:    dataJson,
			}
		} else {
			res = Response{
				Ret:     -1,
				Data:    dataJson,
				Message: b.Log,
			}
		}


	default:
		res = Response{
			Ret:     1,
			Data:    dataJson,
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

func MiddleHandler(ctx context.Context, f func(clictx context.Context, w http.ResponseWriter, r *http.Request), codeSpace sdk.CodespaceType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		async, err := util.CheckBool(r.FormValue("async"))
		if err != nil {
			WriteErrorRes(w, client.ErrParseParam(codeSpace, errors.New("error async")))
			return
		}
		from_str := r.FormValue("from")
		from, err := helper.StrToAddress(from_str)
		if err != nil {
			WriteErrorRes(w, client.ErrParseParam(codeSpace, errors.New("error from")))
			return
		}
		ctx = ctx.WithBlocked(async)
		ctx = ctx.WithFrom(from)

		f(ctx, w, r)
	}
}

func GetNecessaryParams(cliCtx context.Context, request *http.Request, cdc *codec.Codec, broadcast bool) (key string, from sdk.AccAddress, nonce, gas uint64, err error) {
	key = request.FormValue("privateKey")
	from = cliCtx.GetFromAddresses()
	if !broadcast {
		nonce = 0
		gas = 0
		return
	}
	Gas, err2 := strconv.ParseUint(request.FormValue("gas"), 10, 64)
	if err2 != nil || Gas < 0 {
		err = errors.New("gas error")
		return
	}
	gas = Gas
	userNonce := request.FormValue("nonce")
	if userNonce != "" {
		Nonce, err2 := strconv.ParseInt(userNonce, 10, 64)
		if err2 != nil || Nonce < 0 {
			err = errors.New("nonce err")
			return
		}
		nonce = uint64(Nonce)
		return
	}else {
		ctx, err2 := client.NewClientContextFromViper(cdc)
		if err2 != nil {
			err = errors.New("new client context error")
			return
		}
		nonce, _, err2 = ctx.GetNonceByAddress(from, false)
		if err2 != nil {
			err = errors.New("get nonce error")
			return
		}
	}
	return
}