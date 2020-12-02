package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	apptypes "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/gorilla/websocket"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
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
	DefaultRPCPort = "26657"

	DefaultWSEndpoint = "/websocket"
)

var (
	DefaultPort = "26657"
	cdc = amino.NewCodec()
)

func init() {
	cdc.RegisterConcrete(&apptypes.CommonTx{}, "commontx", nil)
}

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
				continue
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
		r.ConnMap[topic] = DeleteSlice(r.ConnMap[topic], c)
		if len(r.ConnMap[topic]) == 0 {
			r.Unsubscribe(topic)
		}
	}
	r.Mutex.Unlock()
}

func (r *PubSubRoom) HandleUnsubscribeAll(c *websocket.Conn) {
	r.Mutex.Lock()
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
		for k, conn := range r.Connections {
			responses, err := conn.Subscribe(ctx, subscribeClient, topic)
			if err != nil {
				delete(r.Connections, k)
				_ = conn.Stop()
				continue
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
					//default:
					//	v = e.Data
					//}
					//response := SendMessage{
					//	Time:   time.Now().Format(time.RFC3339),
					//	Content: v,
					//}
					response := SendMessage{
						Time:   time.Now().Format(time.RFC3339),
						Content: e.Data,
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
				continue
			}
			r.Mutex.Lock()
			r.Connections[addr] = conn
			r.Mutex.Unlock()
			if r.ConnMap != nil {
				for topic := range r.ConnMap {
					responses, err := conn.Subscribe(ctx, subscribeClient, topic)
					if err != nil {
						r.Mutex.Lock()
						delete(r.Connections, addr)
						r.Mutex.Unlock()
						_ = conn.Stop()
						continue
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
							response := SendMessage{
								Time:   time.Now().Format(time.RFC3339),
								Content: e.Data,
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

func Notify(r *PubSubRoom, topic string, m SendMessage) {
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
	res = res + str[0] + ":" + DefaultPort
	return res
}

type GotTxResponse struct {
	Height int64                  `json:"height"`
	Index  uint32                 `json:"index"`
	Tx     string                     `json:"tx"`
	Result abci.ResponseDeliverTx `json:"result"`
}

func NewGotTxResponse(height int64, index uint32, tx string, result abci.ResponseDeliverTx) GotTxResponse {
	return GotTxResponse{
		Height: height,
		Index:  index,
		Tx:     tx,
		Result: result,
	}
}

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