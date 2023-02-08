package main

import (
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {

	servAddr := "localhost:3000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	wg.Add(1)
	go tcpRead(conn, err)
	//go tcpWrite(conn, err, message)
	wg.Wait()
}

func tcpRead(conn *net.TCPConn, err error) {
	defer wg.Done()
	for {
		received := make([]byte, 1024)
		_, err = conn.Read(received)
		if err != nil {
			println("Read from server failed:", err.Error())
			os.Exit(1)
		}
		println(received)
	}
}

func tcpWrite(conn *net.TCPConn, err error, message string) {
	_, err = conn.Write([]byte(message))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
}
