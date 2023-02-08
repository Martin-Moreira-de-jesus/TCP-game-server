package main

import (
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	for {
		message := "Halo"
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

		println("write to server = ", message)
		wg.Add(1)
		go tcpRead(conn, err)
		wg.Add(1)
		go tcpWrite(conn, err, message)

		wg.Wait()
	}
}

func tcpRead(conn *net.TCPConn, err error) {
	defer wg.Done()
	received := make([]byte, 1024)
	_, err = conn.Read(received)
	if err != nil {
		println("Read from server failed:", err.Error())
		os.Exit(1)
	}
}

func tcpWrite(conn *net.TCPConn, err error, message string) {
	defer wg.Done()
	_, err = conn.Write([]byte(message))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
}
