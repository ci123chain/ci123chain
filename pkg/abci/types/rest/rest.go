package rest

import (
	"CI123Chain/pkg/abci/codec"
	"CI123Chain/pkg/client/context"
	"net/http"
	"strconv"
)

// ErrorResponse defines the attributes of a JSON error response.
type ErrorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error"`
}


// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, err string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err}
}


// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(codec.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
}

func PostProcessResponseBare(w http.ResponseWriter, ctx context.Context, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch body.(type) {
	case []byte:
		resp = body.([]byte)
	default:
		resp, err = ctx.Cdc.MarshalJSONIndent(body, "", "  ")
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

func ParseQueryHeightOrReturnBadRequest(w http.ResponseWriter, cliCtx context.Context, r *http.Request) (context.Context, bool) {
	heightStr := r.FormValue("height")
	if heightStr != "" {
		height, err := strconv.ParseInt(heightStr, 10, 64)
		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		if height < 0 {
			WriteErrorResponse(w, http.StatusBadRequest, "height must be equal or greater than zero")
			return cliCtx, false
		}
		if height > 0 {
			cliCtx = cliCtx.WithHeight(height)
		}
	} else {
		cliCtx = cliCtx.WithHeight(0)
	}
	return cliCtx, true
}