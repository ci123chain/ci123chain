package types

import (
	"context"
	//"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	apptypes "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/gorilla/websocket"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
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
	subscribeClient = "abci-client"

	//time out
	timeOut = time.Second * 10

	DefaultTCP = "tcp://"
	DefaultRPCPort = "80"

	DefaultPrefix = "tm."
	DefaultWSEndpoint = "/websocket"
)

var (
	DefaultPort = "80"
	cdc = apptypes.MakeCodec()
)
func SetDefaultPort(port string) {
	DefaultPort = port
}

type PubSubRoom struct {
	ConnMap    map[string][]*websocket.Conn
	HasCreatedConn map[string]bool

	backends 		[]Instance
	Connections    map[string]*rpcclient.HTTP
	Mutex          sync.Mutex
}

func (r *PubSubRoom) SetBackends(bs []Instance) {
	r.backends = bs
}

func (r PubSubRoom) GetBackends() []Instance {
	return r.backends
}

func (r *PubSubRoom)GetPubSubRoom() {
	r.ConnMap = make(map[string][]*websocket.Conn)
	r.HasCreatedConn = make(map[string]bool)
	r.Connections = make(map[string]*rpcclient.HTTP)
	r.backends = make([]Instance, 0)
}

func (r *PubSubRoom) HasClientConnect() bool {
	var has bool
	for _, v := range r.ConnMap {
		if len(v) != 0 {
			has = true
		}
	}
	return has
}

func (r *PubSubRoom)HasTMConnections() bool {
	return len(r.Connections) != 0
}

func (r *PubSubRoom)SetTMConnections() {
	for _, v := range r.backends {
		addr := rpcAddress(v.URL().Host)
		if r.Connections[addr] == nil {
			conn, ok := GetConnection(addr)
			if !ok {
				//continue
				logger.Error(fmt.Sprintf("connect remote addr: %s, failed", addr))
				r.Mutex.Lock()
				r.RemoveAllTMConnections()
				r.Mutex.Unlock()
				break
			}
			r.Connections[addr] = conn
		}
	}
}


func (r *PubSubRoom)Receive(c *websocket.Conn) {

	//query.MustParse("")
	defer func() {
		err := recover()
		switch rt := err.(type) {
		case ClientError:
			r.HandleUnsubscribeAll(rt.Connect)
		default:
			logger.Info("info: %s", rt)
		}
	}()
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			panic(NewServerError(err, c))
		}
		logger.Info("receive: %s", string(data))
		var m ReceiveMessage
		ok := IsValidMessage(data)
		if !ok {
			logger.Info("received invalid message from client")
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("invalid message: %s, you have sent to server", string(data)),
			}
			_ = c.WriteJSON(res)
			continue
		}
		if len(r.backends) != 0 && !r.HasTMConnections() {
			r.Mutex.Lock()
			r.SetTMConnections()
			r.Mutex.Unlock()
		}
		_ = json.Unmarshal(data, &m)
		topic := m.Content.GetTopic()
		//query check
		q, err := query.New(topic)
		if err != nil {
			logger.Info("invalid topic from client")
			res := SendMessage{
				Time:    time.Now().Format(time.RFC3339),
				Content: fmt.Sprintf("invalid topic: %s, you have sent to server", q.String()),
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
		}
	}
}


func (r *PubSubRoom) HandleSubscribe(topic string, c *websocket.Conn) {
	r.Mutex.Lock()
	if len(r.ConnMap[topic]) == 0 {
		r.ConnMap[topic] = make([]*websocket.Conn, 0)
		r.ConnMap[topic] = append(r.ConnMap[topic], c)
		if !r.HasCreatedConn[topic] {
			r.Subscribe(topic)
		}

	}else {
		r.ConnMap[topic] = append(r.ConnMap[topic], c)
	}
	r.Mutex.Unlock()
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
	r.Mutex.Unlock()
}

func (r *PubSubRoom) Subscribe(topic string) {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	//responses, err := r.TMConnection.Subscribe(ctx, subscribeClient, query)
	//if err != nil {
	//	panic(err)
	//}
	//
	//go func() {
	//	for e := range responses {
	//		res, _ := json.Marshal(e.Data)
	//		response := SendMessage{
	//			Time:   time.Now().Format(time.RFC3339),
	//			Content: string(res),
	//		}
	//		Notify(r, topic, response)
	//	}
	//}()
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
				////connection health check
				if _, err := conn.Health(); err != nil {
					////connection failed
					//Err := fmt.Sprintf("sorry, remote connection disconnect, subscribe will stop later...")
					//ErrRes := SendMessage{
					//	Time:   time.Now().Format(time.RFC3339),
					//	Content: Err,
					//}
					//NotifyAll(r, ErrRes)
					r.Mutex.Lock()
					r.RemoveAllTMConnections()
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
			go func() {
				for e := range responses {
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
						//hash := hex.EncodeToString(tx.Tx.Hash())
						//v = NewGotTxResponse(tx.Height, tx.Index, hash, tx.Result)
						v = tx
					default:
						v = e.Data
					}
					//response := SendMessage{
					//	Time:   time.Now().Format(time.RFC3339),
					//	Content: v,
					//}
					response := SendMessage{
						Time:   time.Now().Format(time.RFC3339),
						Content: v,//e.Data,
					}
					Notify(r, topic, response)
				}
			}()
		}
	}()
}

func (r *PubSubRoom) AddShard() {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	//defer cancel()
	for _, v := range r.backends {
		addr := rpcAddress(v.URL().Host)
		if r.Connections[addr] == nil {
			conn, ok := GetConnection(addr)
			if !ok {
				logger.Error(fmt.Sprintf("connect remote addr: %s, failed", addr))
				//Err := fmt.Sprintf("sorry, remote connection disconnect, subscribe will stop later...")
				//ErrRes := SendMessage{
				//	Time:   time.Now().Format(time.RFC3339),
				//	Content: Err,
				//}
				//NotifyAll(r, ErrRes)
				r.Mutex.Lock()
				r.RemoveAllTMConnections()
				r.Mutex.Unlock()
				break
			}
			r.Mutex.Lock()
			r.Connections[addr] = conn
			r.Mutex.Unlock()
			if r.ConnMap != nil {
				for topic := range r.ConnMap {
					responses, err := conn.Subscribe(ctx, subscribeClient, topic)
					//if err != nil {
					//	r.Mutex.Lock()
					//	logger.Error(fmt.Sprintf("subscribe topic: %s, failed", topic))
					//	delete(r.Connections, addr)
					//	r.Mutex.Unlock()
					//	_ = conn.Stop()
					//	continue
					//}
					if err != nil {
						logger.Error(fmt.Sprintf("subscribe topic: %s failed", topic))
						////connection health check
						if _, err := conn.Health(); err != nil {
							////connection failed
							//Err := fmt.Sprintf("sorry, remote connection disconnect, subscribe will stop later...")
							//ErrRes := SendMessage{
							//	Time:   time.Now().Format(time.RFC3339),
							//	Content: Err,
							//}
							//NotifyAll(r, ErrRes)
							r.Mutex.Lock()
							r.RemoveAllTMConnections()
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

					go func() {
						for e := range responses {
							//var v interface{}
							//switch e.Data.(type) {
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
								//hash := hex.EncodeToString(tx.Tx.Hash())
								//v = NewGotTxResponse(tx.Height, tx.Index, hash, tx.Result)
								v = tx
							default:
								v = e.Data
							}
							//response := SendMessage{
							//	Time:   time.Now().Format(time.RFC3339),
							//	Content: v,
							//}
							response := SendMessage{
								Time:   time.Now().Format(time.RFC3339),
								Content: v,//e.Data,
							}
							Notify(r, topic, response)
						}
					}()
				}
			}
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
	for k, conn := range r.Connections {
		err := conn.Unsubscribe(ctx, subscribeClient, topic)
		if err != nil {
			delete(r.Connections, k)
			_ = conn.Stop()
			logger.Error("unsubscribe error: %s", err)
			continue
		}
	}
}

func (r *PubSubRoom) RemoveAllTMConnections() {
	for _, v := range r.ConnMap {
		//r.HasCreatedConn[k] = false
		for _, val := range v {
			if val != nil {
				_ = val.Close()
			}
		}
		//delete(r.ConnMap, k)
	}
	for _, v := range r.Connections {
		_ = v.Stop()
	}
	//r.Clients = make([]*websocket.Conn, 0)
	r.ConnMap = make(map[string][]*websocket.Conn)
	r.HasCreatedConn = make(map[string]bool)
	r.Connections = make(map[string]*rpcclient.HTTP)
}

func Notify(r *PubSubRoom, topic string, m SendMessage) {
	by, _ := json.Marshal(m)
	logger.Debug("publish message: %s, on topic : %s", string(by), topic)
	for i := 0; i < len(r.ConnMap[topic]); i++ {
		conn := r.ConnMap[topic][i]
		if err := conn.WriteJSON(m); err != nil {
			r.HandleUnsubscribeAll(conn)
			_ = conn.Close()
			i--
		}
	}
}

func DeleteSlice(res []*websocket.Conn, s *websocket.Conn) []*websocket.Conn {
	j := 0
	for _, val := range res {
		if val != s {
			res[j] = val
			j++
		}
	}
	return res[:j]
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

func GetConnection(addr string) (*rpcclient.HTTP, bool){
	client := rpcclient.NewHTTP(addr, DefaultWSEndpoint)
	err := client.Start()
	if err != nil {
		logger.Error("connect error: %s", err)
		return nil, false
	}
	return client, true
}

func rpcAddress(host string) string {
	res := DefaultTCP
	str := strings.Split(host, ":")
	res = res + DefaultPrefix + str[0] + ":" + DefaultPort
	return res
}

//type GotTxResponse struct {
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

//type GotTxData struct {
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
//type Attribute struct {
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