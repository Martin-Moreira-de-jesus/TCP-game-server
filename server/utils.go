package main

import (
	"net"
)

func QuickWrite(conn net.Conn, message string) error {
	var _, err = conn.Write([]byte(message + "\n"))
	return err
}

func QuickRead(conn net.Conn, message *string) error {
	var buf = make([]byte, 1024)
	var _, err = conn.Read(buf)
	*message = string(buf)
	return err
}

/*
func LogMessage(conn net.Conn, message string) bool {
    var now = time.Now()
    var log = fmt.Sprintf("[%s]: %s", now.String(), message)
    if QuickWrite(conn, now.String() + log) {
        return true
    }
    return false
}
*/
