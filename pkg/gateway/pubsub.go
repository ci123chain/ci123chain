package gateway

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
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

func checkBackend() {

	//serverPool.SharedCheck()
	t := time.NewTicker(time.Second * 17)
	for {
		select {
		case <-t.C:
			spByte, _ := json.Marshal(serverPool.backends)
			//log.Println("serverpool backends:")
			//log.Println(serverPool.backends)
			spHash := makeHash(spByte)
			prByte, _ := json.Marshal(pubsubRoom.GetBackends())
			prHash := makeHash(prByte)
			if !bytes.Equal(spHash, prHash) {
				pubsubRoom.SetBackends(serverPool.backends)
				//log.Println("get backends:")
				//log.Println(pubsubRoom.GetBackends()[0].URL().Host)
				pubsubRoom.AddShard()
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