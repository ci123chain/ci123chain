package types

import (
	"context"
	"github.com/ci123chain/ci123chain/pkg/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	//"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	apptypes "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	vmtypes "github.com/ci123chain/ci123chain/pkg/vm/client/rest/websockets"
	"github.com/gorilla/websocket"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"strings"
	"sync"
	"time"
)


const (
	DefaultSubscribeMethod = "subscribe"
	DefaultUnsubscribeMethod = "unsubscribe"
	DefaultUnsubscribeAllMethod = "unsubscribe_all"
	//subscribe clinet
	subscribeClient = "abci-clients"
	//time out
	timeOut = time.Second * 10
	DefaultWSEndpoint = "/websocket"
)

var (
	EthPort   string
	ShardPort string
	TMPort    string
)
var cdc = apptypes.GetCodec()
func SetDefaultPort(tmport, shardport, ethport string) {
	TMPort = tmport
	EthPort = ethport
	ShardPort = shardport
}

type PubSubRoom struct {
	RemoteClients []*websocket.Conn
	ConnMap    map[string][]*websocket.Conn
	//HasCreatedConn map[string]bool

	backends 		[]Instance
	Connections    map[string]*rpcclient.HTTP
	EthConnections map[string]*websocket.Conn
	Mutex          sync.Mutex

	IDs            map[float64]*websocket.Conn
	Subs           map[string]*websocket.Conn
	MaxConnections  int
}

func (r *PubSubRoom) SetBackends(bs []Instance) {
	r.backends = bs
}

func (r PubSubRoom) GetBackends() []Instance {
	return r.backends
}

func (r *PubSubRoom)GetPubSubRoom() {
	r.RemoteClients = make([]*websocket.Conn, 0)
	r.ConnMap = make(map[string][]*websocket.Conn, 0)
	//r.HasCreatedConn = make(map[string]bool)
	r.Connections = make(map[string]*rpcclient.HTTP, 0)
	r.backends = make([]Instance, 0)
	r.EthConnections = make(map[string]*websocket.Conn, 0)
	r.IDs = make(map[float64]*websocket.Conn, 0)
	r.Subs = make(map[string]*websocket.Conn, 0)
}

func (r *PubSubRoom) HasClientConnect() bool {
	var has bool
	for _, v := range r.ConnMap {
		if len(v) != 0 {
			has = true
			continue
		}
	}
	return has
}

func (r *PubSubRoom)HasTMConnections() bool {
	return len(r.Connections) != 0
}

func (r *PubSubRoom) HasEthConnections() bool {
	return len(r.EthConnections) != 0
}

func (r *PubSubRoom) SetTMConnections() (err error) {
	for _, v := range r.backends {
		err, info := GetURL(v.URL().Host)
		if err != nil {
			logger.Error(fmt.Sprintf("SetTMConnections: get remote node domain info: %s, failed: %v", v.URL().Host, err))
			r.RemoveAllTMConnections(errors.New(fmt.Sprintf("SetTMConnections: get remote node domain info: %s, failed", v.URL().Host)))
			break
		}
		addr := rpcAddress(info.Host26657)
		if r.Connections[addr] == nil {
			conn, err := GetConnection(addr)
			if err != nil {
				Err := errors.New(fmt.Sprintf("connect remote addr: %s, failed, err: %v", addr, err))
				logger.Error(Err.Error())
				r.RemoveAllTMConnections(err)
				break
			}
			r.Connections[addr] = conn
		}
	}
	return err
}

func (r *PubSubRoom) SetEthConnections() (err error) {
	for _, v := range r.backends {
		var info *util.DomainInfo
		err, info = GetURL(v.URL().Host)
		if err != nil {
			logger.Error(fmt.Sprintf("SetEthConnections: get remote node domain info: %s, failed: %v", v.URL().Host, err))
			r.RemoveAllEthConnections()
			break
		}
		addr := info.Host8546
		if r.EthConnections[addr] == nil {
			conn, err := GetEthConnection(addr)
			if err != nil {
				Err := errors.New(fmt.Sprintf("connect remote addr: %s, failed, error: %v", addr, err))
				logger.Error(Err.Error())
				r.RemoveAllEthConnections()
				break
			}
			r.EthConnections[addr] = conn
		}
	}
	if err == nil {
		r.HandleEthSub()
	}

	return err
}

func (r *PubSubRoom) Receive(c *websocket.Conn) {

	defer func() {
		err := recover()
		switch rt := err.(type) {
		default:
			logger.Warn("recover unexpected error: %s", rt)
		}
	}()
	for {
		mt, data, err := c.ReadMessage()
		if err != nil {
			logger.Warn(fmt.Sprintf("got client message error: %v", err.Error()))
			r.HandleUnsubscribeAll(c)
			_ = c.Close()
			break
		}
		if mt == websocket.PingMessage {
			err = c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				logger.Error("client write message failed: %s", err.Error())
				r.HandleUnsubscribeAll(c)
				_ = c.Close()
				break
			}
		}
		logger.Info("receive clients message: %s", string(data))
		var m ReceiveMessage
		ok := IsValidMessage(data)
		if !ok {
			logger.Info("received invalid message from clients")
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("invalid message: %s, you have sent to server", string(data)),
			}
			_ = c.WriteJSON(res)
			continue
		}
		if len(r.backends) != 0 && !r.HasTMConnections() {
			r.Mutex.Lock()
			err := r.SetTMConnections()
			r.Mutex.Unlock()
			if err != nil {
				res := SendMessage{
					Time:    time.Now().Format(time.RFC3339),
					Content: fmt.Sprintf("connect to tendermint failed, err: %s", err.Error()),
				}
				_ = c.WriteJSON(res)
				r.HandleUnsubscribeAll(c)
				_ = c.Close()
				break
			}
		}else if len(r.backends) == 0 {
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("there has no backends in server pool"),
			}
			_ = c.WriteJSON(res)
			r.HandleUnsubscribeAll(c)
			_ = c.Close()
			break
		}
		_ = json.Unmarshal(data, &m)
		topic := m.Content.GetTopic()
		//query check
		_, err = query.New(topic)
		if err != nil {
			logger.Info("invalid topic from clients")
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("invalid topic: %s, you have sent to server", topic),
			}
			_ = c.WriteJSON(res)
			continue
		}
		switch m.Content.CommandType() {
		case DefaultSubscribeMethod:
			r.HandleSubscribe(topic, c)
		case DefaultUnsubscribeMethod:
			r.HandleUnsubscribe(topic, c)
		case DefaultUnsubscribeAllMethod:
			r.HandleUnsubscribeAll(c)
			_ = c.Close()
		}
	}
}

func (r *PubSubRoom) ReceiveEth(c *websocket.Conn){
	defer func() {
		err := recover()
		switch rt := err.(type) {
		default:
			logger.Warn("recover unexpected error: %s", rt)
		}
	}()

	for {
		mt, data, err := c.ReadMessage()
		if err != nil {
			logger.Warn("got client error, client may down already")
			r.RemoveEthConnection(c, nil)
			_ = c.Close()
			break
		}
		if mt == websocket.PingMessage {
			err = c.WriteMessage(websocket.PongMessage, nil)
			if err != nil {
				logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
				r.RemoveEthConnection(c, nil)
				_ = c.Close()
				break
			}
		}
		var msg EthSubcribeMsg
		err = json.Unmarshal(data, &msg)
		if err != nil {
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("server got invalid message: %s", string(data)),
			}
			_ = c.WriteJSON(res)
			continue
		}
		logger.Info("receive clients message: %s", string(data))

		if len(r.backends) != 0 && !r.HasEthConnections() {
			r.Mutex.Lock()
			err := r.SetEthConnections()
			r.Mutex.Unlock()
			if err != nil {
				res := SendMessage{
					Time:    time.Now().Format(time.RFC3339),
					Content: fmt.Sprintf("connect to server failed, err: %s", err.Error()),
				}
				_ = c.WriteJSON(res)
				r.RemoveEthConnection(c, nil)
				_ = c.Close()
				break
			}
		}else if len(r.backends) == 0 {
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("there has no backends in server pool"),
			}
			_ = c.WriteJSON(res)
			r.RemoveEthConnection(c, nil)
			_ = c.Close()
			break
		}
		if r.IDs[msg.ID] != nil && r.IDs[msg.ID].RemoteAddr().String() != "" {
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("the id: %v,  you input has been used", msg.ID),
			}
			_ = c.WriteJSON(res)
			_ = c.Close()
			break
		}

		if msg.Method == "eth_unsubscribe" {
			r.Mutex.Lock()
			r.Subs = DeleteSubs(r.Subs, c)
			r.Mutex.Unlock()
		} else if msg.Method == "eth_subscribe" {
			r.Mutex.Lock()
			r.IDs[msg.ID] = c
			r.Mutex.Unlock()
		}else {
			continue
		}

		for _, con := range r.EthConnections {
			go func(remote *websocket.Conn, message EthSubcribeMsg) {
				var ok = make(chan bool)
				err := remote.WriteJSON(message)
				if err != nil {
					remoteErr := fmt.Sprintf("remote write message failed: %s", err.Error())
					logger.Warn(remoteErr)
					//r.RemoveEthConnection(client, errors.New(remoteErr))
					ok <- false
					return
				}
				//remote write, client read.
				//go func(ok chan bool) {
				//	for {
				//		mt, Data, err := remote.ReadMessage()
				//		if mt == websocket.PongMessage {
				//			continue
				//		}
				//		if err != nil {
				//			remoteErr := fmt.Sprintf("remote read message failed: %s", err.Error())
				//			logger.Warn(remoteErr)
				//			r.RemoveEthConnection(client, errors.New(remoteErr))
				//			ok <- false
				//			break
				//		}
				//		//err = client.WriteMessage(websocket.BinaryMessage, Data)
				//		err = client.WriteMessage(1, Data)
				//		//err = client.WriteJSON(string(Data))
				//		if err != nil {
				//			logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
				//			r.RemoveEthConnection(client, nil)
				//			break
				//		}
				//	}
				//}(ok)
				////client write, remote read.
				//go func(ok chan bool) {
				//	for {
				//		mt, Data, err := client.ReadMessage()
				//		if mt == websocket.PingMessage {
				//			_ = client.WriteMessage(websocket.PongMessage, nil)
				//			continue
				//		}
				//		if err != nil {
				//			logger.Warn(fmt.Sprintf("read remote message failed: %s", err.Error()))
				//			r.RemoveEthConnection(client, nil)
				//			ok <- false
				//			break
				//		}
				//		err = remote.WriteMessage(websocket.BinaryMessage, Data)
				//		if err != nil {
				//			remoteErr := fmt.Sprintf("client write message failed: %s", err.Error())
				//			logger.Warn(remoteErr)
				//			r.RemoveEthConnection(client, errors.New(remoteErr))
				//			break
				//		}
				//	}
				//}(ok)

				select {
				case _ = <- ok:
					r.Mutex.Lock()
					r.RemoveAllEthConnections()
					r.Mutex.Unlock()
					break
				}
			}(con, msg)
		}
	}
}

func (r *PubSubRoom) HandleEthSub() {

	for _, con := range r.EthConnections {
		go func(remote *websocket.Conn) {
			for {
				mt, Data, err := remote.ReadMessage()
				if mt == websocket.PongMessage || Data == nil {
					continue
				}
				var subres vmtypes.SubscriptionResponseJSON
				err = json.Unmarshal(Data, &subres)
				if err == nil {
					switch subres.Result.(type) {
					case string:
						c := r.IDs[subres.ID]
						if c != nil {
							err = c.WriteMessage(websocket.TextMessage, Data)
							if err != nil {
								logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
								r.RemoveEthConnection(c, nil)
								continue
							}
							r.Mutex.Lock()
							r.IDs = DeleteIDs(r.IDs, subres.ID)
							r.Subs[subres.Result.(string)] = c
							r.Mutex.Unlock()
							continue
						}
						logger.Warn(fmt.Sprintf("no client mathc id: %v", subres.ID))
						continue
					case bool:
						//r.Subs = DeleteSubs(r.Subs, c)
						continue
					default:
						//logger.Warn(fmt.Sprintf("subscription result got invalid type: %v", typ))
						//continue
					}
					//c := r.IDs[subres.ID]
					//if c != nil {
					//	err = c.WriteMessage(websocket.TextMessage, Data)
					//	if err != nil {
					//		logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
					//		r.RemoveEthConnection(c, nil)
					//		continue
					//	}
					//	r.Mutex.Lock()
					//	r.IDs = DeleteIDs(r.IDs, subres.ID)
					//	switch typ := subres.Result.(type) {
					//	case string:
					//		r.Subs[subres.Result.(string)] = c
					//	case bool:
					//		//r.Subs = DeleteSubs(r.Subs, c)
					//	default:
					//		logger.Warn(fmt.Sprintf("subscription result got invalid type: %v", typ))
					//	}
					//	r.Mutex.Unlock()
					//	continue
					//}
				}
				var result vmtypes.SubscriptionNotification
				err = json.Unmarshal(Data, &result)
				if err != nil {
					logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
					//r.Subs = DeleteSubs(r.Subs, client)
					//r.RemoveEthConnection(client, nil)
					continue
				}
				client := r.Subs[string(result.Params.Subscription)]
				if client == nil || client.RemoteAddr().String() == "" {
					logger.Warn(fmt.Sprintf("got empty client in subs with subid: %v", string(result.Params.Subscription)))
					continue
				}
				err = client.WriteMessage(websocket.TextMessage, Data)
				if err != nil {
					logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
					r.Subs = DeleteSubs(r.Subs, client)
					r.RemoveEthConnection(client, nil)
					continue
				}
			}
		}(con)
	}
}

func (r *PubSubRoom) HandleSubscribe(topic string, c *websocket.Conn) {
	//r.Mutex.Lock()
	if len(r.ConnMap[topic]) == 0 {
		r.Subscribe(topic)
		r.ConnMap[topic] = make([]*websocket.Conn, 0)
		r.ConnMap[topic] = append(r.ConnMap[topic], c)

	}else {
		r.ConnMap[topic] = append(r.ConnMap[topic], c)
	}
	//r.Mutex.Unlock()
}

func (r *PubSubRoom) HandleUnsubscribe(topic string, c *websocket.Conn) {
	r.Mutex.Lock()
	if len(r.ConnMap[topic]) == 0 {
		//do nothing
	}else {
		logger.Info("connection address: %s, unsubscribe topic:%s", c.RemoteAddr().String(), topic)
		r.ConnMap[topic] = DeleteSlice(r.ConnMap[topic], c)
		if len(r.ConnMap[topic]) == 0 {
			r.Unsubscribe(topic)
		}
	}
	r.Mutex.Unlock()
}

func (r *PubSubRoom) HandleUnsubscribeAll(c *websocket.Conn) {
	r.Mutex.Lock()
	logger.Info("connection address: %s, unsubscribe all topic", c.RemoteAddr().String())
	for k, v := range r.ConnMap {
		if r.ConnMap[k] != nil {
			r.ConnMap[k] = DeleteSlice(v, c)
			if len(r.ConnMap[k]) == 0 {
				r.Unsubscribe(k)
			}
		}
	}
	r.RemoteClients = DeleteSlice(r.RemoteClients, c)
	r.Mutex.Unlock()
}

func (r *PubSubRoom) Subscribe(topic string) {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Info("info: %s", err)
			}
		}()
		for _, conn := range r.Connections {
			responses, err := conn.Subscribe(ctx, subscribeClient, topic)
			if err != nil {
				logger.Error(fmt.Sprintf("subscribe topic: %s failed", topic))
				if _, err := conn.Health(ctx); err != nil {
					////connection failed
					r.Mutex.Lock()
					r.RemoveAllTMConnections(err)
					r.Mutex.Unlock()
					break
				}else {
					Err := fmt.Sprintf("subscribe topic: %s failed, maybe you should check your topic", topic)
					ErrRes := SendMessage{
						Time:   time.Now().Format(time.RFC3339),
						Content: Err,
					}
					Notify(r, topic, ErrRes)
					continue
				}
			}
			tmRes := TMResponse{
				response: responses,
				topic:    topic,
				addr:     conn.Remote(),
			}
			go func(res TMResponse) {
				var start = time.Now()
				var worked = false
				for e := range res.response {
					logger.Info("receive tendermint response: %v", e.Data)
					var v interface{}
					switch e.Data.(type) {
					case types.EventDataTx:
						tx := e.Data.(types.EventDataTx)
						var aa sdk.Tx
						err = cdc.UnmarshalBinaryBare(tx.Tx, &aa)
						if err != nil {
							logger.Error(fmt.Sprintf("got error tx %x", tx))
						}
						res, err := json.Marshal(aa)
						if err != nil {
							logger.Error(fmt.Sprintf("marshal error: %s", err.Error()))
						}
						tx.Tx = res
						v = tx
					default:
						v = e.Data
					}
					response := SendMessage{
						Time:   time.Now().Format(time.RFC3339),
						Content: v,//e.Data,
					}
					Notify(r, res.topic, response)
					worked = true
				}
				var end = time.Now()
				if !worked {
					if end.Sub(start).Seconds() > 15 {
						Err := fmt.Sprintf("the node: %s, is bad, no result return after 15 senconds", res.addr)
						response := SendMessage{
							Time:   time.Now().Format(time.RFC3339),
							Content: Err,
						}
						Notify(r, res.topic, response)
						r.Mutex.Lock()
						r.RemoveAllTMConnections(errors.New(Err))
						r.Mutex.Unlock()
						return
					}
				}
			}(tmRes)
		}
	}()
}

func (r *PubSubRoom) AddShard() {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	//defer cancel()
	for _, v := range r.backends {
		err, info := GetURL(v.URL().Host)
		if err != nil {
			logger.Error(fmt.Sprintf("AddShard: get remote node domain info: %s, failed: %v", v.URL().Host, err))
			r.Mutex.Lock()
			r.RemoveAllTMConnections(errors.New(fmt.Sprintf("AddShard:get remote node domain info: %s, failed", v.URL().Host)))
			r.Mutex.Unlock()
			break
		}
		addr := rpcAddress(info.Host26657)
		if r.Connections[addr] == nil {
			conn, err := GetConnection(addr)
			if err != nil {
				Err := errors.New(fmt.Sprintf("connect remote addr: %s, failed", addr))
				logger.Error(Err.Error())
				r.Mutex.Lock()
				r.RemoveAllTMConnections(Err)
				r.Mutex.Unlock()
				break
			}
			r.Mutex.Lock()
			r.Connections[addr] = conn
			r.Mutex.Unlock()
			if r.ConnMap != nil {
				for topic := range r.ConnMap {
					responses, err := conn.Subscribe(ctx, subscribeClient, topic)
					if err != nil {
						logger.Error(fmt.Sprintf("subscribe topic: %s failed", topic))
						////connection health check
						if _, err := conn.Health(ctx); err != nil {
							////connection failed
							//Err := fmt.Sprintf("sorry, remote connection disconnect, subscribe will stop later...")
							//ErrRes := SendMessage{
							//	Time:   time.Now().Format(time.RFC3339),
							//	Content: Err,
							//}
							//NotifyAll(r, ErrRes)
							r.Mutex.Lock()
							r.RemoveAllTMConnections(err)
							r.Mutex.Unlock()
							break
						}else {
							//delete(r.Connections, k)
							//_ = conn.Stop()
							Err := fmt.Sprintf("subscribe topic: %s failed, maybe you should check your topic", topic)
							ErrRes := SendMessage{
								Time:   time.Now().Format(time.RFC3339),
								Content: Err,
							}
							Notify(r, topic, ErrRes)
							continue
						}

					}
					tmRes := TMResponse{
						response: responses,
						topic:    topic,
						addr:     conn.Remote(),
					}

					go func(res TMResponse) {
						var start = time.Now()
						var worked = false
						for e := range res.response {
							//var v interface{}
							//switch e.Data.(types) {
							//case types.EventDataTx:
							//	ok, _ := regexp.MatchString("tm.event = 'Tx'", topic)
							//	if !ok {
							//		continue
							//	}
							//	tx := e.Data.(types.EventDataTx)
							//	hash := hex.EncodeToString(tx.Tx.Hash())
							//	v = NewGotTxResponse(tx.Height, tx.Index, hash, tx.Result)
							//	//var sender string
							//	//var operation string
							//	//if len(tx.Result.Events) == 0{
							//	//	logger.Error(fmt.Sprintf("got unexpected tx: %x", tx))
							//	//	continue
							//	//}else {
							//	//	//var success bool
							//	//	//var amount string
							//	//	//var isTransfer bool
							//	//	//for _, kv := range tx.Result.Events[0].Attributes {
							//	//	//	if string(kv.Key) == "sender" {
							//	//	//		sender = string(kv.Value)
							//	//	//	}
							//	//	//	if string(kv.Key) == "operation" {
							//	//	//		operation = string(kv.Value)
							//	//	//		if operation == "transfer" {
							//	//	//			isTransfer = true
							//	//	//		}
							//	//	//	}
							//	//	//}
							//	//	//if isTransfer {
							//	//	//	for _, kv := range tx.Result.Events[0].Attributes {
							//	//	//		if string(kv.Key) == "amount" {
							//	//	//			amount = string(kv.Value)
							//	//	//		}
							//	//	//	}
							//	//	//}
							//	//	//if tx.TxResult.Result.Code == 0 {
							//	//	//	success = true
							//	//	//}else {
							//	//	//	success = false
							//	//	//}
							//	//	//height := tx.Height
							//	//	//gasUsed := tx.Result.GasUsed
							//	//	//gasWanted := tx.Result.GasWanted
							//	//	//hash := hex.EncodeToString(value.Tx.Hash())
							//	//	//a := NewGotTxData(operation, sender , hash ,height, gasUsed, gasWanted, success, amount)
							//	//	//v = a
							//	//}
							//default:
							//	v = e.Data
							//}
							//response := SendMessage{
							//	Time:   time.Now().Format(time.RFC3339),
							//	Content: v,
							//}
							//response := SendMessage{
							//	Time:   time.Now().Format(time.RFC3339),
							//	Content: e.Data,
							//}
							var v interface{}
							switch e.Data.(type) {
							case types.EventDataTx:
								//ok, _ := regexp.MatchString("tm.event = 'Tx'", topic)
								//if !ok {
								//	continue
								//}
								tx := e.Data.(types.EventDataTx)
								var aa sdk.Tx
								err = cdc.UnmarshalBinaryBare(tx.Tx, &aa)
								if err != nil {
									logger.Error(fmt.Sprintf("got error tx %x", tx))
								}
								res, err := json.Marshal(aa)
								if err != nil {
									logger.Error(fmt.Sprintf("marshal error: %s", err.Error()))
								}
								tx.Tx = res
								v = tx
							default:
								v = e.Data
							}
							response := SendMessage{
								Time:   time.Now().Format(time.RFC3339),
								Content: v,//e.Data,
							}
							Notify(r, res.topic, response)
							worked = true
						}
						var end = time.Now()
						if !worked {
							if end.Sub(start).Seconds() > 15 {
								Err := fmt.Sprintf("the node: %s, is bad, no result return after 15 senconds", res.addr)
								response := SendMessage{
									Time:   time.Now().Format(time.RFC3339),
									Content: Err,
								}
								Notify(r, res.topic, response)
								r.Mutex.Lock()
								r.RemoveAllTMConnections(errors.New(Err))
								r.Mutex.Unlock()
								return
							}
						}
					}(tmRes)
				}
			}
		}

		if r.EthConnections[addr] == nil {
			remote, err := GetEthConnection(addr)
			if err != nil {
				logger.Warn(fmt.Sprintf("get connection wtih addr: %v, failed, error: %v", addr, err))
				r.RemoveAllEthConnections()
				break
			}
			go func() {
				for {
					mt, Data, err := remote.ReadMessage()
					if mt == websocket.PongMessage || Data == nil {
						continue
					}
					var result vmtypes.SubscriptionNotification
					err = json.Unmarshal(Data, &result)
					if err != nil {
						logger.Warn(fmt.Sprintf("receive invalid resposne from remote: %s", err.Error()))
						r.RemoveAllEthConnections()
						break
					}
					client := r.Subs[string(result.Params.Subscription)]
					if client == nil || client.RemoteAddr().String() == "" {
						logger.Warn(fmt.Sprintf("got empty client in subs with subid: %v", string(result.Params.Subscription)))
						continue
					}
					err = client.WriteMessage(websocket.BinaryMessage, Data)
					if err != nil {
						logger.Warn(fmt.Sprintf("client write message failed: %s", err.Error()))
						r.RemoveEthConnection(client, nil)
						continue
					}
				}
			}()
		}
	}
}

func (r *PubSubRoom) Unsubscribe(topic string) {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	//defer cancel()
	//query := topic
	//err := r.TMConnection.Unsubscribe(ctx, subscribeClient, query)
	//if err != nil {
	//	panic(err)
	//}
	logger.Info("unsubscribe topic %s", topic)
	for _, conn := range r.Connections {
		err := conn.Unsubscribe(ctx, subscribeClient, topic)
		if err != nil {
			//delete(r.Connections, k)
			//_ = conn.Stop()
			//logger.Error("unsubscribe error: %s", err)
			//continue
			r.RemoveAllTMConnections(err)
		}
	}
}

func (r *PubSubRoom) RemoveAllTMConnections(err error) {
	for _, v := range r.ConnMap {
		//r.HasCreatedConn[k] = false
		for _, val := range v {
			if val != nil {
				//Err := fmt.Sprintf("the node: %s, is bad, no result return after 15 senconds", res.addr)
				response := SendMessage{
					Time:   time.Now().Format(time.RFC3339),
					Content: err.Error(),
				}
				_ = val.WriteJSON(response)
				_ = val.Close()
			}
		}
		//delete(r.ConnMap, k)
	}
	for _, v := range r.Connections {
		_ = v.Stop()
	}
	//r.Clients = make([]*websocket.Conn, 0)
	r.RemoteClients = make([]*websocket.Conn, 0)
	r.backends = make([]Instance, 0)
	r.ConnMap = make(map[string][]*websocket.Conn, 0)
	//r.HasCreatedConn = make(map[string]bool)
	r.Connections = make(map[string]*rpcclient.HTTP, 0)
}

func (r *PubSubRoom) RemoveAllEthConnections() {
	for _, v := range r.EthConnections {
		_ = v.Close()
	}
	for _, v := range r.Subs {
		_ = v.Close()
	}
	for _, v := range r.IDs {
		_ = v.Close()
	}
	r.RemoteClients = make([]*websocket.Conn, 0)
	r.backends = make([]Instance, 0)
	r.ConnMap = make(map[string][]*websocket.Conn, 0)
	r.EthConnections = make(map[string]*websocket.Conn, 0)
	r.IDs = make(map[float64]*websocket.Conn, 0)
	r.Subs = make(map[string]*websocket.Conn, 0)
}

func (r *PubSubRoom) RemoveEthConnection(c *websocket.Conn, err error) {
	if len(r.EthConnections) == 0 {
		return
	}
	r.Mutex.Lock()
	r.RemoteClients = DeleteSlice(r.RemoteClients, c)
	r.Mutex.Unlock()
	if err != nil {
		errRes := SendMessage{
			Time:    time.Now().String(),
			Content: err.Error(),
		}
		_ = c.WriteJSON(errRes)
	}
	_ = c.Close()
}


func Notify(r *PubSubRoom, topic string, m SendMessage) {
	by, _ := json.Marshal(m)
	logger.Debug("publish message: %s, on topic : %s", string(by), topic)
	for i := 0; i < len(r.ConnMap[topic]); i++ {
		conn := r.ConnMap[topic][i]
		if err := conn.WriteJSON(m); err != nil {
			logger.Error("notify message error: %s", err.Error())
			r.HandleUnsubscribeAll(conn)
			r.Mutex.Lock()
			r.RemoteClients = DeleteSlice(r.RemoteClients, conn)
			r.Mutex.Unlock()
			_ = conn.Close()
			i--
		}
	}
}

func DeleteSubs(subs map[string]*websocket.Conn, c *websocket.Conn) map[string]*websocket.Conn{
	result := make(map[string]*websocket.Conn, 0)
	for key, val := range subs {
		fmt.Println(fmt.Sprintf("addr: %v", val.RemoteAddr().String()))
		if val.LocalAddr().String() == c.LocalAddr().String() && val.RemoteAddr().String() != c.RemoteAddr().String() {
			result[key] = val
		}
	}
	fmt.Println(fmt.Sprintf("result: %v", result))
	return result
}

func DeleteIDs(ids map[float64]*websocket.Conn, id float64) map[float64]*websocket.Conn {
	result := make(map[float64]*websocket.Conn, 0)
	for key, val := range ids {
		if key != id {
			result[key] = val
		}
	}
	return result
}

func DeleteSlice(res []*websocket.Conn, s *websocket.Conn) []*websocket.Conn {
	result := make([]*websocket.Conn, 0)
	for _, val := range res {
		if val.RemoteAddr().String() != s.RemoteAddr().String() {
			result = append(result, val)
		}
	}
	return result
}


// Server端消息结构
type SendMessage struct {
	Time   string  `json:"time"`   //time
	Content interface{} `json:"content"` //content表示消息内容，这里简化别的消息
}

//client端消息结构
type ReceiveMessage struct {
	Time   string 			  `json:"time"`   //time
	Content MessageContent    `json:"content"` //content表示消息内容，这里简化别的消息
}

func IsValidMessage(msg []byte) bool {
	var rm ReceiveMessage
	var trueOperation bool
	var trueValueType bool
	err := json.Unmarshal(msg, &rm)
	if err != nil || len(rm.Content.Args) == 0 {
		return false
	}
	for _, v := range rm.Content.Args {
		switch v.ValueType {
		case "string":
			trueValueType = true
		case "number":
			trueValueType = true
		case "date":
			trueValueType = true
		case "time":
			trueValueType = true
		default:
			trueValueType = false
		}

		switch v.Connect {
		case "=":
			trueOperation = true
		case "<":
			trueOperation = true
		case "<=":
			trueOperation = true
		case ">":
			trueOperation = true
		case ">=":
			trueOperation = true
		case "CONTAINS":
			trueOperation = true
		case "EXISTS":
			trueOperation = true
		default:
			trueOperation = false
		}
	}
	if trueOperation && trueValueType {
		return true
	}
	return false
}

func (msg ReceiveMessage) String() string {
	res, _ := json.Marshal(msg)
	return string(res)
}


type MessageContent struct {
	Method    string   	`json:"method"`
	Args      []Description  `json:"args"`
}

type Description struct {
	Key          string   	   `json:"key"`
	Connect      string   	   `json:"connect"`
	Value        string        `json:"value"`
	ValueType    string        `json:"value_type"`
}

func (msg MessageContent) GetTopic() string {
	var topic string
	for k, v := range msg.Args {
		if k != 0 {
			topic += " AND "
		}
		str := fmt.Sprintf("%s %s ", v.Key, v.Connect)
		switch v.ValueType {
		case "string":
			str = str + "'" + v.Value + "'"
		case "time":
			str = str + " TIME " + v.Value
		case "date":
			str = str + " DATE " + v.Value
		default:
			str = str + v.Value
		}
		topic += str
	}
	return topic
}

func (msg MessageContent) CommandType() string {
	switch msg.Method {
	case DefaultSubscribeMethod:
		return DefaultSubscribeMethod
	case DefaultUnsubscribeMethod:
		if msg.IsUnsubscribeAll() {
			return DefaultUnsubscribeAllMethod
		}
		return DefaultUnsubscribeMethod
	default:
		panic(errors.New("unexpected method"))
	}
}

func (msg MessageContent) IsUnsubscribeAll() bool {
	if msg.Method == DefaultUnsubscribeMethod {
		if len(msg.Args) == 0 {
			return true
		}
	}
	return false
}

func GetConnection(addr string) (*rpcclient.HTTP, error){
	client, err := rpcclient.New(addr, DefaultWSEndpoint)
	if err != nil {
		logger.Error("new WSEvent client failed: %s", err)
		return nil, err
	}
	err = client.Start()
	if err != nil {
		logger.Error("connect error: %s", err)
		return nil, err
	}
	return client, nil
}

func GetEthConnection(addr string) (*websocket.Conn, error) {
	str := strings.Split(addr, "//")
	var link string
	if len(str) >= 2 {
		link = strings.Split(str[1], ":")[0] + ":" + EthPort
	}else {
		link = strings.Split(str[0], ":")[0] + ":" + EthPort
	}
	u := url.URL{Scheme: util.DefaultWS, Host: link,  Path: "/"}
	if os.Getenv(util.IDG_APPID) != "" {
		u = url.URL{Scheme: util.DefaultWSS, Host: link,  Path: "/"}
	}
	dialer := websocket.DefaultDialer
	c, _, err := dialer.Dial(u.String(), nil)
	return c, err
}

func rpcAddress(host string) string {
	prefix := util.SchemaPrefix()
	//if os.Getenv(util.IDG_APPID) != "" {
	//	res = util.DefaultHTTPS
	//}
	str := strings.Split(host, ":")
	res := prefix + str[0] + ":" + TMPort
	return res
}

//types GotTxResponse struct {
//	Height int64                  `json:"height"`
//	Index  uint32                 `json:"index"`
//	Tx     string                     `json:"tx"`
//	Result abci.ResponseDeliverTx `json:"result"`
//}
//
//func NewGotTxResponse(height int64, index uint32, tx string, result abci.ResponseDeliverTx) GotTxResponse {
//	return GotTxResponse{
//		Height: height,
//		Index:  index,
//		Tx:     tx,
//		Result: result,
//	}
//}

//types GotTxData struct {
//	TxType    string   `json:"tx_type"`
//	TxSender  string   `json:"tx_sender"`
//	TxHeight  int64   `json:"tx_height"`
//	TxHash    string   `json:"tx_hash"`
//	TxGasUsed  int64   `json:"tx_gas_used"`
//	TxGasWanted int64    `json:"tx_gas_wanted"`
//	TxSuccess  bool    `json:"tx_success"`
//	TxAmount   string   `json:"tx_amount"`
//}
//
//types Attribute struct {
//	Key   string   `json:"key"`
//	Value string   `json:"value"`
//}
//
//func NewGotTxData(ty, sender, hash string, height, gasUsed, gasWanted int64, success bool, amount string) GotTxData {
//	return GotTxData{
//		TxType:      ty,
//		TxSender:    sender,
//		TxHeight:    height,
//		TxHash:      hash,
//		TxGasUsed:   gasUsed,
//		TxGasWanted: gasWanted,
//		TxSuccess:   success,
//		TxAmount:    amount,
//	}
//}

func GetURL(host string) (error, *util.DomainInfo) {
	cli := &http.Client{
		Transport:&http.Transport{DisableKeepAlives:true},
	}
	prefix := util.SchemaPrefix()
	reqUrl := prefix + strings.Split(host, ":")[0] + ":" + ShardPort + "/info"

	req2, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil || req2 == nil {
		return err, nil
	}
	//not use one connection
	req2.Close = true
	rep2, err := cli.Do(req2)
	if err != nil {
		return err, nil
	}
	b, err := ioutil.ReadAll(rep2.Body)
	if err != nil {
		return err, nil
	}
	defer rep2.Body.Close()

	var info util.DomainInfo
	err = json.Unmarshal(b, &info)
	if err != nil {
		return err, nil
	}
	return nil, &info
}

type TMResponse struct {
	response <- chan ctypes.ResultEvent
	topic  string
	addr   string
}

type EthSubcribeMsg struct {
	JsonRpc  string   `json:"jsonrpc"`
	ID       float64   `json:"id"`
	Method   string   `json:"method"`
	Params   []interface{}  `json:"params"`
}