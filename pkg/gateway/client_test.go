package gateway

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"testing"
	"time"
)

type MessageContent struct {
	Method    string   	`json:"method"`
	Args      []string  `json:"args"`
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
	count int = 0
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

	subMsg := Message{
		Time:   time.Now().Format(time.RFC3339),
		Content:MessageContent{
			Method: "subscribe",
			Args:   []string{"NewBlock"},
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

	subMsg := Message{
		Time:   time.Now().Format(time.RFC3339),
		Content:MessageContent{
			Method: "subscribe",
			Args:   []string{"NewBlock"},
		},
	}
	sender.send <- subMsg


	select {
	case _ = <- ch:
		subMsg := Message{
			Time:   time.Now().Format(time.RFC3339),
			Content:MessageContent{
				Method: "unsubscribe",
				Args:   []string{"NewBlock"},
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
		message := Message{
			Time:   time.Now().Format(time.RFC3339),
			Content:MessageContent{
				Method: str,
				Args:   []string{topic},
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
