package main

import (
    "fmt"
    "net"
    "time"
)

func QuickWrite(conn net.Conn, message string) bool {
    var _, err = conn.Write([]byte(message))
    if err != nil {
        return false
    }
    return true
}

func QuickRead(conn net.Conn, message *string) bool {
    var buf = make([]byte, 1024)
    var _, err = conn.Read(buf)
    *message = string(buf)
    if err != nil {
        return false
    }
    return true
}

func LogMessage(conn net.Conn, message string) bool {
    var now = time.Now()
    var log = fmt.Sprintf("[%s]: %s", now.String(), message)
    if QuickWrite(conn, now.String() + log) {
        return true
    }
    return false
}
