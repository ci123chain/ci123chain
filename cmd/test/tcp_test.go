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
		ServerName: 		"weelinknode1.gw106.oneitfarm.com",
		InsecureSkipVerify: true,
	}

	remoteServer := "weelinknode1.gw106.oneitfarm.com:7443"
	i := 0
	for  {
		fmt.Println("Beigin Connection times: ", i)
		i++
		conn, err := tls.DialWithDialer(&net.Dialer{
			Timeout: time.Second * time.Duration(10),
		}, "tcp", remoteServer, &config)

		if err != nil {
			fmt.Printf("remote tls dial fail.  %s", err.Error())
			continue
		}
		_, err = conn.Write([]byte("aaaa"))
		if err != nil {
			fmt.Println("error: ",err)
		}

		time.Sleep(time.Second * 3)
		//defer conn.Close()
	}

	select {
	}
}
