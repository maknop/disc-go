package gateway

import (
	"context"
	"encoding/json"
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
	curr_time = time.Now().Format(time.Kitchen)
	interrupt chan os.Signal
)

func EstablishConnection(ctx context.Context) {
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	connection, _, err := websocket.DefaultDialer.DialContext(ctx, socketUrl, nil)
	if err != nil {
		log.Fatalf("%s: Connection could not be established: %s", curr_time, err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	for {
		log.Printf("%s: Waiting for response from server.", curr_time)
		heartbeat_interval, err := ReceiveMessage(connection)
		if err != nil {
			log.Fatalf("%s: Could not retrieve heartbeat interval: %s", curr_time, err)
		} else {
			log.Printf("%s: SUCCESS: Value of heartbeat interval is: %d milliseconds (%d seconds)", curr_time, heartbeat_interval, heartbeat_interval/1000)
		}

		time.Sleep(time.Duration(heartbeat_interval))

		SendMessage(connection)
	}
}

func ReceiveMessage(connection *websocket.Conn) (int, error) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		log.Fatalf("%s: Error receiving message: %s", curr_time, err)
	} else {
		log.Printf("%s: Received: %s\n", curr_time, msg)
	}

	msgData := []string{}
	jsonMsg, _ := json.MarshalIndent(&msgData, "", " ")
	log.Printf("json data: %s", jsonMsg)
	return strconv.Atoi(string(msg[53:58]))
}

func SendMessage(connection *websocket.Conn) {
	err := connection.WriteMessage(websocket.TextMessage, []byte("op: 1, d: 251"))
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s: Error during writing to websocket: %s", curr_time, err))
	} else {
		log.Printf("%s: Message sent back to server.", curr_time)
	}

}
