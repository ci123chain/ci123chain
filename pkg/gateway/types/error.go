package types

import (
	"github.com/gorilla/websocket"
)

type ClientError struct {
	Error       interface{}    `json:"error"`
	Connect     *websocket.Conn  `json:"connect"`
}


func NewServerError(r interface{}, c *websocket.Conn) ClientError {
	return ClientError{
		Error:   r,
		Connect: c,
	}
}