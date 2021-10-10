package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func setupRoutes() {
	http.HandleFunc("/", gateway.homePage())
	http.HandleFunc("/ws", gateway.wsEndpoint())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error("Failed to set up server")
	}

	log.Info("App startup successful")
}

func main() {
	log.Info("disc-go!")

	handler := gateway.wsEndpoint()
}
