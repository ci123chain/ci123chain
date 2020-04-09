package types

import (
	"net/http"
)


type Proxy interface {

	Handle(r *http.Request, backends []Instance, RequestParams map[string]string) []byte

	Response() *chan []byte
}
