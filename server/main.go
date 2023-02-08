package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func handleConnection(conn net.Conn) {
    LogMessage(conn, "Joining lobby...")
    var gameId, playerId = CreateOrJoinGame()
    LogMessage(conn, fmt.Sprintf("Joined lobby [%s] as [%s]", gameId, playerId))

    // wait for game to start
    for {
        LogMessage(conn, "Waiting for lobby to fill...")
        if GameStarted(gameId) {
            break
        }
        time.Sleep(1 * time.Second)
    }

    LogMessage(conn, "Game starting !")

    var

    // game loop
    for {

    }
}

func main() {
    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            // handle error
        }
        go handleConnection(conn)
    }
}