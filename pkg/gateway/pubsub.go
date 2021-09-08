package gateway

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
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

const (
	MaxConnection = 10
)

func PubSubHandle(w http.ResponseWriter, r *http.Request) {
	if len(pubsubRoom.RemoteClients) >= MaxConnection {
		res, _ := json.Marshal(types.ErrorResponse{
			Ret: -1,
			Message:  "max clinet has coonected, connection refused",
		})
		_, _ = w.Write(res)
		return
	}
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
	pubsubRoom.Mutex.Lock()
	pubsubRoom.RemoteClients = append(pubsubRoom.RemoteClients, conn)
	pubsubRoom.Mutex.Unlock()
	//根据订阅的topic来建立新的map.
	// map [topic] -> conn
	go pubsubRoom.Receive(conn)
}

func EthPubSubHandle(w http.ResponseWriter, r *http.Request) {
	if len(pubsubRoom.RemoteClients) >= MaxConnection {
		res, _ := json.Marshal(types.ErrorResponse{
			Ret: -1,
			Message:  "max clinet has coonected, connection refused",
		})
		_, _ = w.Write(res)
		return
	}
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
	pubsubRoom.Mutex.Lock()
	pubsubRoom.RemoteClients = append(pubsubRoom.RemoteClients, conn)
	pubsubRoom.Mutex.Unlock()
	go pubsubRoom.ReceiveEth(conn)
}

func checkEthBackend() {
	t := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-t.C:
			logger.Debug("Start eth backend check...")
			for _, v := range pubsubRoom.EthConnections {
				_ = v.WriteMessage(websocket.PingMessage, nil)
			}
		}
	}
}


func checkBackend() {
	t := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-t.C:
			logger.Debug("Start backend check...")
			var connectErr bool
			for k, v := range pubsubRoom.Connections {
				_, err := v.Health(context.Background())
				if err != nil {
					logger.Warn("lost connect on:", k)
					logger.Warn("health check error: ", err.Error())
					connectErr = true
					break
				}
			}
			if connectErr {
				logger.Warn("Lost connect, remove all subscribe")
				pubsubRoom.Mutex.Lock()
				pubsubRoom.RemoveAllTMConnections(errors.New("remote node is bad"))
				pubsubRoom.Mutex.Unlock()
			}
			spByte, _ := json.Marshal(serverPool.backends)
			spHash := makeHash(spByte)
			prByte, _ := json.Marshal(pubsubRoom.GetBackends())
			prHash := makeHash(prByte)
			if !bytes.Equal(spHash, prHash) {
				pubsubRoom.Mutex.Lock()
				logger.Info("set backends")
				pubsubRoom.SetBackends(serverPool.backends)
				pubsubRoom.Mutex.Unlock()
				if pubsubRoom.HasClientConnect() {
					logger.Info("add shard subscribe")
					pubsubRoom.AddShard()
				}
			}
			logger.Debug("Backend check completed")
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