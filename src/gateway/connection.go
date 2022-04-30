package gateway

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	socketUrl = "wss://gateway.discord.gg/?v=9&encoding=json"
	interrupt chan os.Signal
)

func EstablishConnection() {
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	connection, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		fmt.Println("Connection could not be established: ", err)
	}

	ReceiveMessage(connection)

	defer connection.Close()
}

func ReceiveMessage(connection *websocket.Conn) {
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error receiving message: ", err)
			return
		}

		log.Printf("Received: %s\n", msg)
	}
}
