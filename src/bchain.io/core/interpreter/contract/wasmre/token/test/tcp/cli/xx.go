package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	Start(os.Args[1])
	time.Sleep(time.Second)
}

func Start(tcpAddrStr string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddrStr)
	if err != nil {
		log.Printf("Resolve tcp addr failed: %v\n", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial to server failed: %v\n", err)
		return
	}

	go SendMsg(conn)

	buf := make([]byte, 1024)
	for {
		log.Printf("recv start:\n")
		//conn.SetReadDeadline(time.Now().Add(10*time.Second))
		length, err := conn.Read(buf)
		if err != nil {
			log.Printf("recv server msg failed: %v\n", err)
			conn.Close()
			break
		}
		fmt.Println(string(buf[0:length]))
	}
	fmt.Println("exit!")
}


func SendMsg(conn net.Conn) {
	username := conn.LocalAddr().String()
	for {
		var input string
		fmt.Scanln(&input)

		if input == "/q" || input == "/quit" {
			fmt.Println("Byebye ...")
			conn.Close()
			time.Sleep(time.Second)
			os.Exit(0)
		}

		if len(input) > 0 {
			msg := username + " say:" + input
			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				break
			}
		}
	}
}
