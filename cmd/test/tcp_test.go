package main

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestTcp(t *testing.T)  {
	_, err := net.DialTimeout("tcp", "localhost:26656", time.Second * 10)
	if err != nil {
		fmt.Println(err)
	}
	//conn.Write([]byte("aaaa"))
	select {
	}
}
