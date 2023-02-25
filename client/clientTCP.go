package main

import (
	"net"
	"os"
	"sync"
)

type SafeCounter struct {
	mu          sync.Mutex
	instruction []string
}

type ClientInfos struct {
	conn     *net.TCPConn
	servAddr string
}

var State = ClientInfos{
	conn: &net.TCPConn{},
}

var wg sync.WaitGroup

func (ci ClientInfos) Client() {

	ci.servAddr = "localhost:8080"
	tcpAddr, err := net.ResolveTCPAddr("tcp", ci.servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	State.conn, _ = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}

	wg.Add(1)
	go tcpRead()
	//go tcpWrite(conn) Don't need it here, old commit
	wg.Wait()
}

func SendButtonPressed(keyPressed string) {
	println("yooo i'm in, button pressed")
	message := keyPressed + "=true"
	go tcpWrite(message)
}

func tcpRead() {

	c := SafeCounter{}
	defer wg.Done()
	for {
		received := make([]byte, 1024)
		//println(State.conn)
		_, err := State.conn.Read(received)
		if err != nil {
			println("Read from server failed:", err.Error())
			os.Exit(1)
		}
		c.Lock(string(received))
	}
}

func tcpWrite(message string) {
	println(message)
	_, err := State.conn.Write([]byte(message))
	println("done")
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}
}
