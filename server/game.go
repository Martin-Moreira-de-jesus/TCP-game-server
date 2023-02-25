package main

import (
	"github.com/google/uuid"
	"sync"
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
			if e.val.players.Len() <= 1 {
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
	return game.players.Len() == 2
}
