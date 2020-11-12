package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/logger"
	"github.com/gorilla/websocket"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"strings"
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
)

func SetDefaultPort(port string) {
	DefaultPort = port
}

type PubSubRoom struct {
	ConnMap    map[string][]*websocket.Conn

	TMConnection  *rpcclient.HTTP
	HasCreatedConn map[string]bool

	backends 		[]Instance
	Connections    map[string]*rpcclient.HTTP
}

func (r *PubSubRoom) SetBackends(bs []Instance) {
	r.backends = bs
}

func (r PubSubRoom) GetBackends() []Instance {
	return r.backends
}

func (r *PubSubRoom)GetPubSubRoom() {
	r.ConnMap = make(map[string][]*websocket.Conn)
	//r.TMConnection = GetConnection(rpcAddress)
	r.HasCreatedConn = make(map[string]bool)
	r.Connections = make(map[string]*rpcclient.HTTP)
	r.backends = make([]Instance, 0)
}


func (r *PubSubRoom)Receive(c *websocket.Conn) {
	defer func() {
		err := recover()
		switch rt := err.(type) {
		case ClientError:
			r.HandleUnsubscribeAll(rt.Connect)
		default:
			logger.Info("info: %s", r)
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
				Content: fmt.Sprintf("invalid %s message you have sent to server", string(data)),
			}
			_ = c.WriteJSON(res)
			continue
		}
		_ = json.Unmarshal(data, &m)
		topic := m.Content.GetTopic()
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
	if len(r.ConnMap[topic]) == 0 {
		r.ConnMap[topic] = make([]*websocket.Conn, 0)
		r.ConnMap[topic] = append(r.ConnMap[topic], c)

		if !r.HasCreatedConn[topic] {
			r.Subscribe(topic)
		}

	}else {
		r.ConnMap[topic] = append(r.ConnMap[topic], c)
	}
}

func (r *PubSubRoom) HandleUnsubscribe(topic string, c *websocket.Conn) {
	if len(r.ConnMap[topic]) == 0 {
		//do nothing
	}else {
		r.ConnMap[topic] = DeleteSlice(r.ConnMap[topic], c)
		if len(r.ConnMap[topic]) == 0 {
			r.Unsubscribe(topic)
		}
	}
}

func (r *PubSubRoom) HandleUnsubscribeAll(c *websocket.Conn) {
	for k, v := range r.ConnMap {
		if r.ConnMap[k] != nil {
			r.ConnMap[k] = DeleteSlice(v, c)
			if len(r.ConnMap[k]) == 0 {
				r.Unsubscribe(k)
			}
		}
	}
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
				continue
			}
			go func() {
				for e := range responses {
					res, _ := json.Marshal(e.Data)
					response := SendMessage{
						Time:   time.Now().Format(time.RFC3339),
						Content: string(res),
					}
					Notify(r, topic, response)
				}
			}()
		}
	}()
}

func (r *PubSubRoom) AddShard() {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	for _, v := range r.backends {
		addr := rpcAddress(v.URL().Host)
		if r.Connections[addr] == nil {
			conn := GetConnection(addr)
			r.Connections[addr] = conn
			if r.ConnMap != nil {
				for topic := range r.ConnMap {
					responses, err := conn.Subscribe(ctx, subscribeClient, topic)
					if err != nil {
						delete(r.Connections, addr)
						continue
					}

					go func() {
						for e := range responses {
							res, _ := json.Marshal(e.Data)
							response := SendMessage{
								Time:   time.Now().Format(time.RFC3339),
								Content: string(res),
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
	Content string `json:"content"` //content表示消息内容，这里简化别的消息
}

//client端消息结构
type ReceiveMessage struct {
	Time   string 			  `json:"time"`   //time
	Content MessageContent    `json:"content"` //content表示消息内容，这里简化别的消息
}

func IsValidMessage(msg []byte) bool {
	var rm ReceiveMessage
	err := json.Unmarshal(msg, &rm)
	if err != nil {
		return false
	}
	return true
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
	MessageType  string   	   `json:"message_type"`
	Key          string   	   `json:"key"`
	Connect      string   	   `json:"connect"`
	Value        string        `json:"value"`
	ValueType    string        `json:"value_type"`
}

func (msg MessageContent) GetTopic() string {

	if len(msg.Args) == 0 {
		panic(errors.New("unexpected number of args"))
	}
	var topic string
	for k, v := range msg.Args {
		if k != 0 {
			topic += " AND "
		}
		str := fmt.Sprintf("%s.%s %s ", v.MessageType, v.Key, v.Connect)
		switch v.ValueType {
		case "string":
			str = str + "'" + v.Value + "'"
		case "Time":
			str = str + " TIME " + v.Value
		case "Date":
			str = str + " DATE " + v.Value
		default:
			str = str + v.Value
		}
		topic += str
	}
	return topic
}

//
//func (msg MessageContent) String() string {
//	var str string
//	for _, v := range msg.Args {
//		str += v
//		if len(msg.Args) > 1 {
//			str += " "
//		}
//	}
//	return fmt.Sprintf("Method: %s, Args: %s", msg.Method, str)
//}

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

func GetConnection(addr string) *rpcclient.HTTP {
	client := rpcclient.NewHTTP(addr, DefaultWSEndpoint)
	err := client.Start()
	if err != nil {
		panic(err)
	}
	return client
}

func rpcAddress(host string) string {
	res := DefaultTCP
	str := strings.Split(host, ":")
	res = res + str[0] + ":" + DefaultPort
	return res
}