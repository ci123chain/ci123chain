package gateway

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var (
	ug       = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func PubSubHandle(w http.ResponseWriter, r *http.Request) {
	logg := zap.S()
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		logg.Error(err)
		return
	}
	//根据订阅的topic来建立新的map.
	// map [topic] -> conn
	go pubsubRoom.Receive(conn)
}