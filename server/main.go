package main

import (
    "fmt"
	"log"
    "math/rand"
    "net"
	"strconv"
    "strings"
    "sync"
	"time"
)

type PlayerState struct {
    mu   sync.Mutex
    up   bool
    down bool
}

var cn = make(chan *Game)

func handleConnection(conn net.Conn) {
    defer conn.Close()
	if QuickWrite(conn, "Joining lobby...") != nil {
		return
	}

	var game, youPlayer = CreateOrJoinGame()

	if QuickWrite(conn, fmt.Sprintf("Joined lobby!")) != nil {
		return
	}

	// wait for game to start
	for {
		if QuickWrite(conn, "Waiting for lobby to fill...") != nil {
			return
		}

		if game.val.CanStart() {
			break
		}

		time.Sleep(1 * time.Second)
	}

    fmt.Println("Attempting to start game...")

    game.val.LaunchGameLoopIfNotRunning(cn)

	if QuickWrite(conn, "Game starting !") != nil {
		return
	}

	State.mu.Lock()

    var otherPlayer *Node[Player]
	for e := game.val.players.First(); e != nil; e = e.Next() {
        if e.val.uuid != youPlayer.val.uuid {
            otherPlayer = e
        }
    }

	State.mu.Unlock()

    go handleUserInput(conn, youPlayer)

	for {
		State.mu.Lock()
		if QuickWrite(
			conn,
			fmt.Sprintf(
				"you=%s,other=%s,pipeX=%s,obstacleY1=%s,obstacleY2=%s",
				strconv.Itoa(youPlayer.val.posY),
				strconv.Itoa(otherPlayer.val.posY),
                strconv.Itoa(game.val.obstacleX),
                strconv.Itoa(game.val.obstacleY1),
                strconv.Itoa(game.val.obstacleY2),
			),
		) != nil {
			return
		}
		State.mu.Unlock()
		time.Sleep(30 * time.Millisecond)
	}
}

func handleUserInput(conn net.Conn, player *Node[Player]) {
    defer conn.Close()
    for {
        var message string
        if QuickRead(conn, &message) == nil {
            return
        }

        State.mu.Lock()
        var values = strings.Split(message, ",")
        for _, e := range values {
            var data = strings.Split(e, "=")
            if data[0] == "up" {
                player.val.up, _ = strconv.ParseBool(data[1])
            } else if data[1] == "down" {
                player.val.down, _ = strconv.ParseBool(data[1])
            }
        }
        State.mu.Unlock()
    }
}

func handleGameloop() {
    var game = <- cn
    fmt.Println("Starting game loop !")
    for {
        State.mu.Lock()

        // edit game
        game.obstacleX -= 10

        if game.obstacleX <= 0 {
            game.obstacleX = 1000
            game.obstacleY1 = rand.Intn(800)
            game.obstacleY2 = game.obstacleY1 + 100
        }

        // edit players
        for player := game.players.First(); player != nil; player = player.Next() {
            if player.val.up {
                player.val.posY += 10
            } else if player.val.down {
                player.val.posY += 10
            }
        }

        State.mu.Unlock()
        time.Sleep(50 * time.Millisecond)
    }
}

func main() {
    if true {
        ln, err := net.Listen("tcp", ":8080")
        if err != nil {
            log.Fatal(err)
        }
        go handleGameloop()
        for {
            conn, err := ln.Accept()
            if err != nil {
                // handle error
            }

            go handleConnection(conn)
        }
    }
}
