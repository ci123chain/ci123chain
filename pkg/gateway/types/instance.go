package types

import (
	"net/http/httputil"
	"net/url"
)

type Instance interface {
	IsAlive() bool
	SetAlive(live bool)
	URL() *url.URL
	Proxy() *httputil.ReverseProxy
}