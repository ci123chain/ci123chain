package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/gateway/types"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
	"time"
)

//type Description struct {
//	MessageType  string   	   `json:"message_type"`
//	Key          string   	   `json:"key"`
//	Connect      string   	   `json:"connect"`
//	Value        string        `json:"value"`
//	ValueType    string        `json:"value_type"`
//}

type MessageContent struct {
	Method    string   	`json:"method"`
	Args      []types.Description  `json:"args"`
}

// 消息结构
type Message struct {
	Time   string `json:"time"`   //time
	Content MessageContent `json:"content"` //content表示消息内容，这里简化别的消息
}

type Sender struct {
	conn *websocket.Conn
	send chan Message
}

var (
	count = 0
	ch = make(chan int, 1)
	sendch = make(chan int, 1)
)

func TestSubscribeNewBlock(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3030", Path: "/pubsub"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	sender := &Sender{
		conn: c,
		send: make(chan Message, 128),
	}
	//go sender.Message()
	go sender.loopSendMessage()
	go sender.Receive()
	defer func() {
		if r := recover(); r != nil {
			_ = c.Close()
			log.Fatal(fmt.Sprintf("something error: %v", r))
		}
	}()
	var args = types.Description{
		Key:         "tm.event",
		Connect:     "=",
		Value:       "NewBlock",
		ValueType:   "string",
	}

	subMsg := Message{
		Time:   time.Now().Format(time.RFC3339),
		Content:MessageContent{
			Method: "subscribe",
			Args:   []types.Description{args},
		},
	}
	sender.send <- subMsg


	select {
	case _ = <- ch:
		log.Println("exit")
	}
}

func TestSubscribeContract(t *testing.T) {

	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3030", Path: "/pubsub"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	sender := &Sender{
		conn: c,
		send: make(chan Message, 128),
	}
	//go sender.Message()
	go sender.loopSendMessage()
	go sender.Receive()
	defer func() {
		if r := recover(); r != nil {
			_ = c.Close()
			log.Fatal(fmt.Sprintf("something error: %v", r))
		}
	}()
	var arg1 = types.Description{
		Key:         "tm.event",
		Connect:     "=",
		Value:       "Tx",
		ValueType:   "string",
	}
	var arg2 = types.Description{
		Key:         "contract.operation",
		Connect:     "=",
		Value:       "init_contract",
		ValueType:   "string",
	}
	//var arg3 = Description{
	//	MessageType: "contract",
	//	Key:         "module",
	//	Connect:     "=",
	//	Value:       "wasm",
	//	ValueType:   "string",
	//}
	// tx.event = 'Tx' AND contract.operation = 'init_contract'
	subMsg := Message{
		Time:   time.Now().Format(time.RFC3339),
		Content:MessageContent{
			Method: "subscribe",
			Args:   []types.Description{arg1, arg2},
		},
	}
	sender.send <- subMsg


	select {
	case _ = <- ch:
		log.Println("exit")
	}
}

func TestSubscribeNewBlockAndUnsubscribeNewBlock(t *testing.T) {
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:3030", Path: "/pubsub"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	sender := &Sender{
		conn: c,
		send: make(chan Message, 128),
	}
	//go sender.Message()
	go sender.loopSendMessage()
	go sender.Receive()
	defer func() {
		if r := recover(); r != nil {
			_ = c.Close()
			log.Fatal(fmt.Sprintf("something error: %v", r))
		}
	}()
	var args = types.Description{
		Key:         "tm.event",
		Connect:     "=",
		Value:       "NewBlock",
		ValueType:   "string",
	}

	subMsg := Message{
		Time:   time.Now().Format(time.RFC3339),
		Content:MessageContent{
			Method: "subscribe",
			Args:   []types.Description{args},
		},
	}
	sender.send <- subMsg

	select {
	case _ = <- ch:
		subMsg := Message{
			Time:   time.Now().Format(time.RFC3339),
			Content:MessageContent{
				Method: "unsubscribe",
				Args:   []types.Description{args},
			},
		}
		sender.send <- subMsg
	case _ = <-sendch:
		log.Println("exit")
	}
}

func (sender *Sender) Message() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatal(fmt.Sprintf("something error: %v", r))
		}
	}()
	for {
		var str string
		var topic string
		fmt.Println("input:")
		_, err := fmt.Scanln(&str, &topic)
		if err != nil {
			panic(err)
		}
		var Topic []types.Description
		if topic == "" {
			Topic = nil
		}else {
			err = json.Unmarshal([]byte(topic), &Topic)
			if err != nil {
				panic(err)
			}
		}
		message := Message{
			Time:   time.Now().Format(time.RFC3339),
			Content:MessageContent{
				Method: str,
				Args:   Topic,
			},
		}
		sender.send <- message
	}
}

func (sender *Sender) Receive() {
	defer func() {
		if r := recover(); r != nil {
			_ = sender.conn.Close()
			log.Fatal(fmt.Sprintf("something error: %v", r))
		}
	}()
	for {
		_, data, err := sender.conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		fmt.Println("receive:")
		fmt.Println(string(data))
		count += 1
		fmt.Println(count)
		if count == 3 {
			ch <- 0
		}
	}
}

// 循环发送消息
func (sender *Sender) loopSendMessage() {
	for {
		m := <-sender.send
		if err := sender.conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息", m)
		if m.Content.Method == "unsubscribe" {
			time.Sleep(time.Second * 10)
			sendch <- 0
		}
	}
}
