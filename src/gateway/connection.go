package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
		//jitter := rand.Float32()
		log.Printf("%s: Waiting for response from server.", curr_time)
		heartbeat_interval, sequence_num := ReceiveMessage(connection)
		if err != nil {
			log.Fatalf("%s: Could not retrieve heartbeat interval: %s", curr_time, err)
		} else {
			log.Printf("%s: SUCCESS: Value of heartbeat interval is: %d milliseconds", curr_time, heartbeat_interval)
		}

		//time.Sleep(time.Duration(float32(heartbeat_interval) * jitter))
		time.Sleep(time.Duration(*heartbeat_interval))
		SendMessage(connection, sequence_num)
	}
}

func ReceiveMessage(connection *websocket.Conn) (*int, *int) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		log.Fatalf("%s: Error receiving message: %s", curr_time, err)
	} else {
		log.Printf("%s: Received: %s\n", curr_time, msg)
	}

	var op_10_hello Payload

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_10_hello); err != nil {
		log.Fatalf("%s: Error parsing json data: %s", curr_time, err)
	}

	return op_10_hello.D.Heartbeat_Interval, op_10_hello.S
}

func SendMessage(connection *websocket.Conn, sequence_num *int) {
	var op_code *int
	op_code = new(int)
	*op_code = 1

	op_1_heartbeat := Payload{
		OP: op_code,
		D:  Data{Sequence: sequence_num},
	}

	op_1_heartbeat_json, err := json.Marshal(op_1_heartbeat)
	println(strings.ToLower(string(op_1_heartbeat_json)))
	if err != nil {
		log.Fatalf("%s: Error converting to json data: %s", curr_time, err)
	}

	err = connection.WriteJSON(op_1_heartbeat_json)
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s: Error during writing to websocket: %s", curr_time, err))
	} else {
		log.Printf("%s: Message sent back to server.", curr_time)
	}
}
