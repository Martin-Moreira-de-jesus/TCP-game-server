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
	defer State.mu.Unlock()

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

	go handleUserInput(conn, youPlayer, game)

	for {
		State.mu.Lock()

		if !game.val.running || game == nil {
			if game == nil {
				game = nil
			}
			return
		}

		var youPlayerPos = youPlayer.val.posY
		if !youPlayer.val.alive {
			youPlayerPos = -1
		}

		var otherPlayerPos = otherPlayer.val.posY
		if !otherPlayer.val.alive {
			otherPlayerPos = -1
		}

		if QuickWrite(
			conn,
			fmt.Sprintf(
				"you=%d,other=%d,pipeX=%d,obstacleYTop=%d,obstacleYBottom=%d",
				youPlayerPos,
				otherPlayerPos,
				game.val.obstacleX,
				game.val.obstacleYTop,
				game.val.obstacleYBottom,
			),
		) != nil {
			return
		}
		State.mu.Unlock()
		time.Sleep(30 * time.Millisecond)
	}
}

func handleUserInput(conn net.Conn, player *Node[Player], game *Node[Game]) {
	defer conn.Close()
	for {
		var message string

		if QuickRead(conn, &message) != nil {
			return
		}
		if game == nil || !game.val.running {
			return
		}

		State.mu.Lock()
		var values = strings.Split(message, ",")
		for _, e := range values {
			var data = strings.Split(e, "=")
			if len(data) > 2 {
				fmt.Println("bad input")
				break
			}

			if data[0] == "up" {
				player.val.up, _ = strconv.ParseBool(data[1])
				//println("up : ", player.val.up)
			} else if data[0] == "down" {
				player.val.down, _ = strconv.ParseBool(data[1])
				//println("down : ", player.val.down)
			}
		}
		State.mu.Unlock()
	}
}

func handleGameLoop(game *Game) {
	var deadPlayers = 0
	var speed = 10
	var score = 0
	defer State.mu.Unlock()

	game.obstacleX = 1200
	game.obstacleYTop = rand.Intn(600)
	game.obstacleYBottom = game.obstacleYTop + 200

	for {
		State.mu.Lock()

		if deadPlayers == 2 {
			println("Your score is : ", score)
			println("Congrats !!!, EZ win")
			game.running = false
			return
		}

		// edit game
		game.obstacleX -= speed

		if game.obstacleX <= 0 {
			game.obstacleX = 1200
			game.obstacleYTop = rand.Intn(600)
			game.obstacleYBottom = game.obstacleYTop + 200
			speed = speed + 1
			score = score + 1
		}

		// edit players
		// we condider that the player is a square of 50 by
		for player := game.players.First(); player != nil; player = player.Next() {
			if !player.val.alive {
				continue
			}

			if game.obstacleX >= 50 && game.obstacleX <= 300 {
				if (player.val.posY+30 >= game.obstacleYBottom) || (player.val.posY-30 <= game.obstacleYTop) {
					player.val.alive = false
					deadPlayers += 1
					continue
				}
			}

			if player.val.up && player.val.posY > 0 {
				player.val.posY -= 10
			} else if player.val.down && player.val.posY < 800 {
				player.val.posY += 10
			}
		}

		State.mu.Unlock()
		time.Sleep(50 * time.Millisecond)
	}
}

func gameLoopListener() {
	for {
		var game = <-cn
		go handleGameLoop(game)
	}
}

func debug() {
	for {
		State.mu.Lock()

		for e := State.games.First(); e != nil; e = e.Next() {
			fmt.Println(e)
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
		go gameLoopListener()
		// go debug()
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
			}

			go handleConnection(conn)
		}
	}
}
