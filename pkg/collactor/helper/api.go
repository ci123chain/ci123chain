package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SuccessJSONResponse prepares data and writes a HTTP success
func SuccessJSONResponse(status int, v interface{}, w http.ResponseWriter) {
	WriteSuccessResponse(status, v, w)
}

type RespStruct struct {
	Ret  int `json:"ret"`
	Data interface{} `json:"data, omitempty"`
	Err  string `json:"err, omitempty"`
}

func success(data interface{}) []byte {
	rs := RespStruct{
		Ret: 1,
		Data: data,
	}
	bz, err := json.Marshal(rs)

	if err != nil {
		panic(err)
	}
	return bz
}

func errorResp(data string) []byte {
	bz, _ := json.Marshal(RespStruct{
		Ret: -1,
		Err: data,
	})
	return bz
}

// WriteSuccessResponse writes a HTTP success given a status code and data
func WriteSuccessResponse(statusCode int, data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	if _, err := w.Write(success(data)); err != nil {
		fmt.Printf("Write failed: %v", err)
	}
}

// WriteErrorResponse writes a HTTP error given a status code and an error message
func WriteErrorResponse(statusCode int, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	if _, err := w.Write(errorResp(err.Error())); err != nil {
		fmt.Printf("Write failed: %v", err)
	}
}
