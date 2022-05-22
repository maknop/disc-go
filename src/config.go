package main

import (
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var (
	curr_time = time.Now().Format(time.Kitchen)
)

func LoadEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: Could not load environment file: %s", curr_time, err)
	} else {
		log.Printf("%s: .env file successfully loaded", curr_time)
	}
}
