package types

import (
	"net/http"
)


type Proxy interface {

	Handle(w http.ResponseWriter,r *http.Request, backends []Instance)
}
