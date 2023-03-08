package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/withmandala/go-log"
)

var Logger = log.New(os.Stdout).WithColor()

type PlayerState struct {
	mu   sync.Mutex
	up   bool
	down bool
}

var cn = make(chan *Game)

func CleanUp(game *Node[Game], player *Node[Player]) {
	State.mu.Lock()
	defer State.mu.Unlock()
	if game.val.running || game.val.over {
		return
	}
	if game.val.players.Len() == 1 {
		State.games.Remove(game)
	} else {
		game.val.players.Remove(player)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	if QuickWrite(conn, "Joining lobby...") != nil {
		return
	}

	var game, youPlayer = CreateOrJoinGame()

	if QuickWrite(conn, fmt.Sprintf("Joined lobby!")) != nil {
		return
	}
	// in case game crashes, clean it up
	defer CleanUp(game, youPlayer)

	// wait for game to start
	for {
		if QuickWrite(conn, "Waiting for lobby to fill...") != nil {
			return
		}

		if game.val.CanStart() {
			break
		}

		time.Sleep(300 * time.Millisecond)
	}

	game.val.LaunchGameLoopIfNotRunning(cn)

	if QuickWrite(conn, "Game starting !") != nil {
		return
	}

	go handleUserInput(conn, youPlayer, game)

	defer State.mu.Unlock()
	for {
		State.mu.Lock()
		if game == nil || !game.val.running {
			if game == nil {
				game = nil
			} else {
				State.games.Remove(game)
			}
			return
		}

		var positions = game.val.Stringify(youPlayer)

		if QuickWrite(conn, strings.Join(positions, ",")) != nil {
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
			} else if data[0] == "down" {
				player.val.down, _ = strconv.ParseBool(data[1])
			}
		}
		State.mu.Unlock()
	}
}

func handleGameLoop(game *Game) {
	var deadPlayers = 0
	var obstacleSpeed = 10
	var playersSpeed = 10
	defer State.mu.Unlock()

	game.obstacleX = 1200
	game.obstacleYTop = rand.Intn(600)
	game.obstacleYBottom = game.obstacleYTop + 200

	for {
		State.mu.Lock()

		if deadPlayers == Cfg.Game.MaxPlayers {
			Logger.Infof("Game %s ended", game.uuid)
			game.running = false
			game.over = true
			return
		}

		// edit game
		game.obstacleX -= obstacleSpeed

		if game.obstacleX <= 0 {
			game.obstacleX = 1200
			game.obstacleYTop = rand.Intn(600)
			game.obstacleYBottom = game.obstacleYTop + 200
			obstacleSpeed++
			playersSpeed++
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
				player.val.posY -= playersSpeed
			} else if player.val.down && player.val.posY < 800 {
				player.val.posY += playersSpeed
			}
		}

		if playersSpeed <= 20 {
			playersSpeed++
		}

		State.mu.Unlock()
		time.Sleep(50 * time.Millisecond)
	}
}

func gameLoopListener() {
	Logger.Info("Game loop listener started")
	for {
		Logger.Info("Waiting for a game to start")
		var game = <-cn
		Logger.Infof("Game %s started", game.uuid)
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
		InitConfig()
		ln, err := net.Listen("tcp", fmt.Sprintf(":%s", Cfg.Server.Port))
		Logger.Infof("Server started on port %s", Cfg.Server.Port)
		Logger.Infof("Game size set to %d", Cfg.Game.MaxPlayers)
		if err != nil {
			Logger.Error(err)
			os.Exit(2)
		}
		go gameLoopListener()

		for {
			conn, err := ln.Accept()
			if err != nil {
				Logger.Warn("User connection failed")
			} else {
				Logger.Info("User connected")
			}
			go handleConnection(conn)
		}
	}
}
