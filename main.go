package main

import (
	"fmt"

	client "github.com/maknop/disc-go/client"

	log "github.com/sirupsen/logrus"
)

func Start() {
	fmt.Print("attempting to authenticate user")

	if err := client.AuthenticateUser(); err != nil {
		log.Fatal("failed to authenticate user")
	}
}
