package main

import (
	"strconv"
	"strings"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Client struct {
		Address string `yaml:"address"`
	} `yaml:"client"`
}

func (c *SafeCounter) Lock(received string) {
	c.mu.Lock()
	Parser(received)
	c.mu.Unlock()
}

func Parser(message string) {
	var values = strings.Split(message, ",")
	if len(values) > 1 {
		GameState.players = make([]Player, 0)
		for _, value := range values {
			var operands = strings.Split(value, "=")
			var key = operands[0]
			var val, _ = strconv.Atoi(operands[1])
			if key == "you" {
				var player = Player{
					isMe: true,
					posY: val,
				}
				GameState.players = append(GameState.players, player)
			} else if strings.Contains(key, "other") {
				var player = Player{
					isMe: false,
					posY: val,
				}
				GameState.players = append(GameState.players, player)
			} else if key == "pipeX" {
				GameState.pipeX = val
			} else if key == "obstacleYTop" {
				GameState.obstacleY1 = val
			} else if key == "obstacleYBottom" {
				GameState.obstacleY2 = val
			}
		}
	}
}

var Cfg Config

func InitConfig() {
	ReadFile(&Cfg)
}

func ReadFile(cfg *Config) {
	f, err := os.Open("config.yaml")
	if err != nil {
		println("Error reading file")
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		println("Error reading file")
	}
}
