package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type EnvConfig struct {
	botToken string
}

var config = EnvConfig{}

func init() {
	if err:= godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	config.botToken = os.Getenv("bot_token")
}