package types

import (
	"net/http"
)


type Proxy interface {

	Handle(r *http.Request, backends []Instance, reqBody []byte)

	Response() *chan []byte
}
