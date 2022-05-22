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
		log.Fatalf("%s: [ OP CODE 10 ] Could not retrieve heartbeat interval: %s", curr_time, err)
	} else {
		log.Println("Environment variables loaded!")
	}
}
