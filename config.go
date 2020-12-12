package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type EnvConfig struct {
	botToken    string
	postgresUrl string
}

var config = EnvConfig{}

func init() {
	// in case vars are provided with file and not in process.env
	godotenv.Load(".env")

	config.botToken = os.Getenv("bot_token")
	config.postgresUrl = os.Getenv("postgres_url")

	if config.botToken == "" || config.postgresUrl == "" {
		log.Fatal("Missing env variables")
	}
}
