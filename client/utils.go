package main

import (
	"strconv"
	"strings"
)

func (c *SafeCounter) Lock(received string) {
	c.mu.Lock()
	Parser(received)
	c.mu.Unlock()
}

func Parser(message string) {
	print("msg ", message)
	var values = strings.Split(message, ",")
	if len(values) > 1 {
		GameState.myposy, _ = strconv.Atoi(strings.Split(values[0], "=")[1])
		GameState.otherpos, _ = strconv.Atoi(strings.Split(values[1], "=")[1])
		GameState.pipeX, _ = strconv.Atoi(strings.Split(values[2], "=")[1])
		GameState.obstacleY1, _ = strconv.Atoi(strings.Split(values[3], "=")[1])
		GameState.obstacleY2 = GameState.obstacleY1 + 200 // bug using strconv, reading server message
		/*for _, e := range values {
			data := strings.Split(e, "=")
			println(data[1])
		}*/
	}
}
