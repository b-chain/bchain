package main

import (
	"net"
	"log"
	"fmt"
)

func main() {
	port := "9090"
	Start(port)
}


func Start(port string) {
	host := ":" + port

	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		log.Printf("resolve tcp addr failed: %v\n", err)
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Printf("listen tcp port failed: %v\n", err)
		return
	}

	conns := make(map[string]net.Conn)

	messageChan := make(chan string, 10)

	go BroadMessages(&conns, messageChan)

	for {
		fmt.Printf("listening port %s ...\n", port)
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("Accept failed:%v\n", err)
			continue
		}

		conns[conn.RemoteAddr().String()] = conn
		fmt.Println(conns)

		go Handler(conn, &conns, messageChan)
	}
}

func BroadMessages(conns *map[string]net.Conn, messages chan string) {
	for {
		msg := <-messages
		fmt.Println(msg)

		for key, conn := range *conns {
			fmt.Println("connection is connected from ", key)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Printf("broad message to %s failed: %v\n", key, err)
				delete(*conns, key)
			}
		}
	}
}

func Handler(conn net.Conn, conns *map[string]net.Conn, messages chan string) {
	fmt.Println("connect from client ", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("read client message failed:%v\n", err)
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}

		recvStr := string(buf[0:length])
		messages <- recvStr
	}
}
