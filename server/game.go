package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type SafeState struct {
	mu    sync.Mutex
	games List[Game]
}

var State = SafeState{
	games: List[Game]{},
}

type Player struct {
	uuid  string
	posY  int
	up    bool
	down  bool
	alive bool
}

type Game struct {
	uuid            string
	obstacleX       int
	obstacleYTop    int
	obstacleYBottom int
	players         List[Player]
	running         bool
	over            bool
}

func CreateOrJoinGame() (gameJoined *Node[Game], playerCreated *Node[Player]) {
	var newPlayer = Player{
		uuid:  uuid.New().String(),
		posY:  450,
		alive: true,
		up:    false,
		down:  false,
	}

	State.mu.Lock()
	defer State.mu.Unlock()

	if State.games.Len() != 0 {
		for e := State.games.First(); e != nil; e = e.Next() {
			if e.val.players.Len() < Cfg.Game.MaxPlayers {
				var playerElem = e.val.players.PushBack(newPlayer)
				return e, playerElem
			}
		}
	}

	// no game found, create one
	var newGame = Game{
		uuid:      uuid.New().String(),
		obstacleX: 1000,
		players:   List[Player]{},
		running:   false,
		over:      false,
	}

	var playerElem = newGame.players.PushBack(newPlayer)

	var gameElem = State.games.PushBack(newGame)

	return gameElem, playerElem
}

func (game *Game) LaunchGameLoopIfNotRunning(cn chan *Game) {
	State.mu.Lock()

	if !game.running {
		game.running = true
		cn <- game
	}

	State.mu.Unlock()
}

func (game *Game) IsRunnning() bool {
	State.mu.Lock()
	defer State.mu.Unlock()
	return game.running
}

func (game *Game) CanStart() bool {
	State.mu.Lock()
	defer State.mu.Unlock()
	return game.players.Len() == Cfg.Game.MaxPlayers
}

func (game *Game) Stringify(currentPlayer *Node[Player]) []string {
	var result = make([]string, 0)
	/** Stringify players state */
	var i = 0
	for e := game.players.First(); e != nil; e = e.Next() {
		var pos = e.val.posY
		if !e.val.alive {
			pos = -1
		}
		if e == currentPlayer {
			result = append(result, fmt.Sprintf("you=%d", pos))
		} else {
			i++
			result = append(result, fmt.Sprintf("other%d=%d", i, pos))
		}
	}

	/** stringify game state */
	result = append(result, fmt.Sprintf("pipeX=%d", game.obstacleX))
	result = append(result, fmt.Sprintf("obstacleYTop=%d", game.obstacleYTop))
	result = append(result, fmt.Sprintf("obstacleYBottom=%d", game.obstacleYBottom))
	return result
}
