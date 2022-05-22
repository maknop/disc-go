package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	opcodes "github.com/maknop/disc-go/src/types"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	socketUrl = "wss://gateway.discord.gg/?v=9&encoding=json&compress?=true"
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

	// OP 1 Heartbeat
	SendMessage(connection, sequence_num)

	// OP 11 ACK
	ACK(connection)

	// OP 2 Identity
	Identity(connection)

	Ready(connection)
}

func ReceiveMessage(connection *websocket.Conn) (int, *int) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		log.Fatalf("%s: [ OP CODE 10 ] Error receiving message: %s", curr_time, err)
	} else {
		log.Printf("%s: [ OP CODE 10 ] Received: %s\n", curr_time, msg)
	}

	var op_10_hello opcodes.OP_10_Hello

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_10_hello); err != nil {
		log.Fatalf("%s: Error parsing json data: %s", curr_time, err)
	}

	return op_10_hello.D.Heartbeat_Interval, op_10_hello.S
}

func SendMessage(connection *websocket.Conn, sequence_num *int) {
	op_1_heartbeat := opcodes.OP_1_Heartbeat{
		OP: 1,
		D:  opcodes.OP_1_Heartbeat_Data{Sequence: sequence_num},
	}

	err := connection.WriteJSON(op_1_heartbeat)
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s: Error during writing to websocket: %s", curr_time, err))
	} else {
		log.Printf("%s: [ OP CODE 01 ] Sending the following payload: {op: %d, d: {Seq: %v}}", curr_time, op_1_heartbeat.OP, op_1_heartbeat.D.Sequence)
	}
}

func ACK(connection *websocket.Conn) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		log.Fatalf("%s: [ OP CODE 11 ] Error receiving message: %s", curr_time, err)
	} else {
		log.Printf("%s: [ OP CODE 11 ] Received: %s\n", curr_time, msg)
	}

	var op_11_ack opcodes.OP_11_Heartbeat_ACK

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_11_ack); err != nil {
		log.Fatalf("%s: Error parsing json data: %s", curr_time, err)
	}

}

func Identity(connection *websocket.Conn) {
	op_2_identity := opcodes.OP_2_Identity{
		OP: 2,
		D: opcodes.OP_2_Identity_Data{
			Token:   "",
			Intents: 7,
			Properties: opcodes.OP_2_Identity_Properties{
				OS:      "Linux",
				Browser: "disc-go",
				Device:  "disc-go",
			},
		},
	}

	err := connection.WriteJSON(op_2_identity)
	if err != nil {
		log.Fatalf(fmt.Sprintf("%s: Error during writing to websocket: %s", curr_time, err))
	} else {
		log.Printf("%s: [ OP CODE 02 ] Sending the following payload: ", curr_time)
	}
}

func Ready(connection *websocket.Conn) {
	ready := opcodes.Ready{}
}
