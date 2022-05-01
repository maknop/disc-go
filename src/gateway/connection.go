package gateway

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

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

	for {
		heartbeat_interval, _ := ReceiveMessage(connection)
		log.Printf("Value of heartbeat interval is: %d milliseconds (%d seconds)", heartbeat_interval, heartbeat_interval/1000)

		time.Sleep(time.Duration(heartbeat_interval))

		SendMessage(connection)
	}
}

func ReceiveMessage(connection *websocket.Conn) (int, error) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		panic(fmt.Sprintf("Error receiving message: %s", err))
	}

	log.Printf("Received: %s\n", msg)

	return strconv.Atoi(string(msg[53:58]))
}

func SendMessage(connection *websocket.Conn) {
	err := connection.WriteMessage(websocket.TextMessage, []byte("op: 1, d: 251"))
	if err != nil {
		log.Println(fmt.Sprintf("Error during writing to websocket: %s", err))
		panic("Shutting down...")
	} else {
		log.Printf("Message sent back to server.")
	}

}
