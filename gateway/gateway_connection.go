package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	opcodes "github.com/maknop/disc-go/types"

	"github.com/gorilla/websocket"
)

var (
	socketUrl = "wss://gateway.discord.gg/?v=9&encoding=json&compress?=true"
	curr_time = time.Now().Format(time.Kitchen)
)

func EstablishConnection(ctx context.Context) error {
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, socketUrl, nil)
	if err != nil {
		return fmt.Errorf("%s: connection could not be established: %s", curr_time, err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	fmt.Printf("%s: [ OP CODE 10 ] sending initial request to Discord Gateway server", curr_time)
	heartbeat_interval, sequence_num, err := ReceiveMessage(connection)
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 10 ] could not retrieve heartbeat interval: %s", curr_time, err)
	}

	fmt.Printf("%s: [ OP CODE 10 ] value of heartbeat interval is: %d seconds", curr_time, (heartbeat_interval / 1000))

	// OP 1 Heartbeat
	if err := SendMessage(connection, sequence_num); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s: error during writing to websocket: %s", curr_time, err))
	}

	// OP 11 ACK
	if err := ACK(connection); err != nil {
		return err
	}

	// OP 2 Identity
	if err := Identity(connection); err != nil {
		return err
	}

	if err := Ready(connection); err != nil {
		return err
	}

	return nil
}

func ReceiveMessage(connection *websocket.Conn) (int, *int, error) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		return 0, nil, fmt.Errorf("%s: [ OP CODE 10 ] error receiving message: %s", curr_time, err)
	}

	fmt.Printf("%s: [ OP CODE 10 ] received: %s\n", curr_time, msg)

	var op_10_hello opcodes.OP_10_Hello

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_10_hello); err != nil {
		return 0, nil, fmt.Errorf("%s: [ OP CODE 10 ] error parsing json data: %s", curr_time, err)
	}

	return op_10_hello.D.Heartbeat_Interval, op_10_hello.S, nil
}

func SendMessage(connection *websocket.Conn, sequence_num *int) error {
	op_1_heartbeat := opcodes.OP_1_Heartbeat{
		OP: 1,
		D:  opcodes.OP_1_Heartbeat_Data{Sequence: sequence_num},
	}

	err := connection.WriteJSON(op_1_heartbeat)
	if err != nil {
		return err
	}

	fmt.Printf("%s: [ OP CODE 01 ] sending the following payload: {op: %d, d: {Seq: %v}}", curr_time, op_1_heartbeat.OP, op_1_heartbeat.D.Sequence)

	return nil
}

func ACK(connection *websocket.Conn) error {
	var op_11_ack opcodes.OP_11_Heartbeat_ACK

	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 11 ] error receiving message: %s", curr_time, err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_11_ack); err != nil {
		return fmt.Errorf("%s: [OP CODE 11 ] error parsing json data: %s", curr_time, err)
	}

	fmt.Printf("%s: [ OP CODE 11 ] received: %s\n", curr_time, msg)

	return nil
}

func Identity(connection *websocket.Conn) error {
	auth_token := os.Getenv("AUTH_TOKEN")

	op_2_identity := opcodes.OP_2_Identity{
		OP: 2,
		D: opcodes.OP_2_Identity_Data{
			Token:   auth_token,
			Intents: 513,
			Properties: opcodes.OP_2_Identity_Properties{
				OS:      "Linux",
				Browser: "disc-go",
				Device:  "disc-go",
			},
		},
	}

	err := connection.WriteJSON(op_2_identity)
	if err != nil {
		return fmt.Errorf("%s: error during writing to websocket: %s", curr_time, err)
	}

	fmt.Printf("%s: [ OP CODE 02 ] sending Identity payload", curr_time)

	return nil
}

func Ready(connection *websocket.Conn) error {
	var ready opcodes.OP_0_Ready

	fmt.Printf("%s: [ OP CODE 0 ] reading message returned from server", curr_time)
	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 0 ] error receiving message: %s", curr_time, err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&ready); err != nil {
		return fmt.Errorf("%s: [ OP CODE 0 ] error parsing json data: %s", curr_time, err)
	}

	fmt.Printf("%s: [ OP CODE 0 ] Received: %s", msg, curr_time)

	return nil
}
