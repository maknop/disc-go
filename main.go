package main

import (
	"disc-go/gateway"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("disc-go!")
	establishConnection()
}

func establishConnection() {
	gateway.Connect()
}
