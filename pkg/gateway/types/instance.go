package types

import (
	"net/url"
)

type Instance interface {
	IsAlive() bool
	SetAlive(live bool)
	URL() *url.URL
	FailTime() int
}