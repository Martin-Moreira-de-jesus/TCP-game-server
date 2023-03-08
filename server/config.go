package main

import (
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Game struct {
		MaxPlayers int `yaml:"maxplayers"`
	} `yaml:"game"`
}

var Cfg Config

func processError(err error) {
	Logger.Error(err)
	os.Exit(2)
}

func InitConfig() {
	ReadFile(&Cfg)
}

func ReadFile(cfg *Config) {
	f, err := os.Open("config.yaml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}
