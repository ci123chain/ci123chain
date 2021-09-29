package main

import (
	"crypto/tls"
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

func TestTcpTls(t *testing.T)  {
	var cert tls.Certificate
	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		ServerName:         "weelinknode1c.gw002.oneitfarm.com",
		InsecureSkipVerify: true,
	}

	//remote, err := tls.Dial("tcp", remoteServer, &config)

	remoteServer := "weelinknode1c.gw002.oneitfarm.com:7443"

	remote, err := tls.DialWithDialer(&net.Dialer{
		Timeout: time.Second * time.Duration(10),
	}, "tcp", remoteServer, &config)

	if err != nil {
		fmt.Printf("remote tls dial fail.  %s", err.Error())
		return
	}

	defer remote.Close()
	select {
	}
}
