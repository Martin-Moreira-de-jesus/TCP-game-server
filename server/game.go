package main

import (
    "github.com/google/uuid"
    "sync"
)

type SafeState struct {
    mu sync.Mutex
    games []Game
}

var state = SafeState{
    games: make([]Game, 0),
}

type Player struct {
    uuid string
    posY int
    alive bool
}

type Game struct {
    uuid string
    cactusX int
    players []Player
}

func CreateOrJoinGame() (gameUUID string, playerUUID string) {
    var newPlayer = Player{
        uuid:  uuid.New().String(),
        posY:  0,
        alive: true,
    }

    for _, game := range state.games {
        if len(game.players) <= 1 {
            state.mu.Lock()

            game.players = append(game.players, newPlayer)

            state.mu.Unlock()

            return game.uuid, newPlayer.uuid
        }
    }

    var newGame = Game{
        uuid: uuid.New().String(),
        cactusX: 0,
        players: make([]Player, 1),
    }
    newGame.players[0] = newPlayer

    state.mu.Lock()

    state.games = append(state.games, newGame)

    state.mu.Unlock()

    return newGame.uuid, newPlayer.uuid
}

func GameStarted(gameUUID string) bool {
    for _, game := range state.games {
        if game.uuid == gameUUID {
            if len(game.players) >= 2 {
                return true
            }
        }
    }
    return false
}
