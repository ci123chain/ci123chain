package gateway

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("err: %s", err)
		res, _ := json.Marshal(types.ErrorResponse{
			Ret: -1,
			Message:  fmt.Sprintf("invalid request you have sent to server, err: %s", err.Error()),
		})
		_, _ = w.Write(res)
		return
	}
	//根据订阅的topic来建立新的map.
	// map [topic] -> conn
	go pubsubRoom.Receive(conn)
}

func checkBackend() {
	t := time.NewTicker(time.Second * 17)
	for {
		select {
		case <-t.C:
			spByte, _ := json.Marshal(serverPool.backends)
			spHash := makeHash(spByte)
			prByte, _ := json.Marshal(pubsubRoom.GetBackends())
			prHash := makeHash(prByte)
			if !bytes.Equal(spHash, prHash) {
				pubsubRoom.Mutex.Lock()
				pubsubRoom.SetBackends(serverPool.backends)
				pubsubRoom.Mutex.Unlock()
				if pubsubRoom.HasClientConnect() {
					pubsubRoom.AddShard()
				}
			}
		}
	}
}

func makeHash(code []byte) []byte {
	//get hash
	Md5Inst := md5.New()
	Md5Inst.Write(code)
	Result := Md5Inst.Sum([]byte(""))
	return Result
}