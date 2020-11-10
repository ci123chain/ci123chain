package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"time"
)


const (
	DefaultSubscribeMethod = "subscribe"
	DefaultUnsubscribeMethod = "unsubscribe"
	DefaultUnsubscribeAllMethod = "unsubscribe_all"

	//subscribe clinet
	subscribeClient = "test-client"

	//time out
	timeOut = time.Second * 10

	DefaultRPCAddress = "tcp://0.0.0.0:26657"

	DefaultWSEndpoint = "/websocket"

	TMEVENT = "tm.event = "
	ANDCONNECT = " AND "
)

type PubSubRoom struct {
	ConnMap    map[string][]*websocket.Conn

	TMConnection  *rpcclient.HTTP
	HasCreatedConn map[string]bool
}

func (r *PubSubRoom)GetPubSubRoom(rpcAddress string) {
	r.ConnMap = make(map[string][]*websocket.Conn)
	r.TMConnection = GetConnection(rpcAddress)
	r.HasCreatedConn = make(map[string]bool)
}


func (r *PubSubRoom)Receive(c *websocket.Conn) {
	defer func() {
		err := recover()
		switch rt := err.(type) {
		case ClientError:
			r.HandleUnsubscribeAll(rt.Connect)
			fmt.Println(r)
		default:
			fmt.Println("here")
			fmt.Println(r)
		}
	}()
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			panic(NewServerError(err, c))
		}
		fmt.Println("receive:")
		fmt.Println(string(data))
		var m ReceiveMessage
		ok := IsValidMessage(data)
		if !ok {
			panic(errors.New("invalid message from client"))
		}
		_ = json.Unmarshal(data, &m)
		topic := m.Content.Args[0]
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
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	query := "tm.event = " + "'" + topic + "'"
	respones, err := r.TMConnection.Subscribe(ctx, subscribeClient, query)
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range respones {
			res, _ := json.Marshal(e.Data)
			response := SendMessage{
				Time:   time.Now().Format(time.RFC3339),
				Content: string(res),
			}
			Notify(r, topic, response)
		}
	}()
}

func (r *PubSubRoom) Unsubscribe(topic string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	query := "tm.event = " + "'" + topic + "'"
	err := r.TMConnection.Unsubscribe(ctx, subscribeClient, query)
	if err != nil {
		panic(err)
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
	Args      []string  `json:"args"`
}

func (msg MessageContent) GetTopic() string {

	if len(msg.Args) == 0 {
		panic(errors.New("unexpected number of args"))
	}
	var topic = TMEVENT + "'"
	for k, v := range msg.Args {
		if k == 0 {
			topic += v
			topic += "'"
		}else {
			topic += ANDCONNECT
			topic += v
		}
	}
	return topic
}


func (msg MessageContent) String() string {
	var str string
	for _, v := range msg.Args {
		str += v
		if len(msg.Args) > 1 {
			str += " "
		}
	}
	return fmt.Sprintf("Method: %s, Args: %s", msg.Method, str)
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
		if len(msg.Args) == 1 && msg.Args[0] == "all" {
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