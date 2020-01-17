package types

import (
	"net/http"
)


type Proxy interface {

	//HaveResponse() chan int

	Handle(r *http.Request, backends []Instance, reqBody []byte) ([]byte, error)

	Response() chan []byte
}
