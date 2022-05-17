package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	socketUrl = "wss://gateway.discord.gg/?v=9&encoding=json"
	curr_time = time.Now().Format(time.Kitchen)
)

func EstablishConnection(ctx context.Context) {
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, socketUrl, nil)
	if err != nil {
		log.Fatalf("%s: Connection could not be established: %s", curr_time, err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	log.Printf("%s: [ OP CODE 10 ] Sending initial request to Discord Gateway server", curr_time)
	heartbeat_interval, sequence_num := ReceiveMessage(connection)
	if err != nil {
		log.Fatalf("%s: [ OP CODE 10 ] Could not retrieve heartbeat interval: %s", curr_time, err)
	} else {
		log.Printf("%s: [ OP CODE 10 ] Value of heartbeat interval is: %d seconds", curr_time, (heartbeat_interval / 1000))
	}

	time.Sleep(time.Duration((heartbeat_interval / 1000) * int(time.Second)))

	SendMessage(connection, sequence_num)

	time.Sleep(time.Duration((heartbeat_interval / 1000) * int(time.Second)))

	ACK(connection)
}

func ReceiveMessage(connection *websocket.Conn) (int, *int) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		log.Fatalf("%s: [ OP CODE 10 ] Error receiving message: %s", curr_time, err)
	} else {
		log.Printf("%s: [ OP CODE 10 ] Received: %s\n", curr_time, msg)
	}

	var op_10_hello OP_10_Hello

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_10_hello); err != nil {
		log.Fatalf("%s: Error parsing json data: %s", curr_time, err)
	}

	return op_10_hello.D.Heartbeat_Interval, op_10_hello.S
}

func SendMessage(connection *websocket.Conn, sequence_num *int) {
	op_1_heartbeat := OP_1_Heartbeat{
		OP: 1,
		D:  sequence_num,
	}

	op_1_heartbeat_json, err := json.Marshal(op_1_heartbeat)
	if err != nil {
		log.Fatalf("%s: Error converting to json data: %s", curr_time, err)
	}

	err = connection.WriteJSON(op_1_heartbeat_json)
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s: Error during writing to websocket: %s", curr_time, err))
	} else {
		log.Printf("%s: [ OP CODE 01 ] Sending the following payload: %s", curr_time, strings.ToLower(string(op_1_heartbeat_json)))
	}
}

func ACK(connection *websocket.Conn) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		log.Fatalf("%s: [ OP CODE 11 ] Error receiving message: %s", curr_time, err)
	} else {
		log.Printf("%s: [ OP CODE 11 ] Received: %s\n", curr_time, msg)
	}

	var op_11_ack OP_11_Heartbeat_ACK

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_11_ack); err != nil {
		log.Fatalf("%s: Error parsing json data: %s", curr_time, err)
	}

}
