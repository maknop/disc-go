package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	opcodes "github.com/maknop/disc-go/types"
	utils "github.com/maknop/disc-go/utils"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	socketURL = "wss://gateway.discord.gg/?v=9&encoding=json&compress?=true"
)

func EstablishConnection(ctx context.Context) error {
	log.Info("Establishing connection...")
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, socketURL, nil)
	if err != nil {
		return fmt.Errorf("connection could not be established: %s", err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	fmt.Printf("[ OP CODE 10 ] sending initial request to Discord Gateway server", utils.GetCurrTimeUTC())
	heartbeatInterval, sequenceNum, err := ReceiveMessage(connection)
	if err != nil {
		return fmt.Errorf("[ OP CODE 10 ] could not retrieve heartbeat interval: %s", err)
	}

	fmt.Printf("[ OP CODE 10 ] value of heartbeat interval is: %d seconds", (heartbeatInterval / 1000))

	// OP 1 Heartbeat
	if err := SendMessage(connection, sequenceNum); err != nil {
		return fmt.Errorf(fmt.Sprintf("error during writing to websocket: %s", err))
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
		return 0, nil, fmt.Errorf("[ OP CODE 10 ] error receiving message: %s", err)
	}

	fmt.Printf("[ OP CODE 10 ] received: %s\n", msg)

	var OP10Hello opcodes.OP10Hello

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&OP10Hello); err != nil {
		return 0, nil, fmt.Errorf("[ OP CODE 10 ] error parsing json data: %s", err)
	}

	return OP10Hello.D.HeartbeatInterval, OP10Hello.S, nil
}

func SendMessage(connection *websocket.Conn, sequenceNum *int) error {
	op1Heartbeat := opcodes.OP1Heartbeat{
		OP: 1,
		D:  opcodes.OP1HeartbeatData{Sequence: sequenceNum},
	}

	err := connection.WriteJSON(op1Heartbeat)
	if err != nil {
		return err
	}

	fmt.Printf("[ OP CODE 01 ] sending the following payload: {op: %d, d: {Seq: %v}}", op1Heartbeat.OP, op1Heartbeat.D.Sequence)

	return nil
}

func ACK(connection *websocket.Conn) error {
	var op11ACK opcodes.OP11HeartbeatACK

	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("[ OP CODE 11 ] error receiving message: %s", err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op11ACK); err != nil {
		return fmt.Errorf("[OP CODE 11 ] error parsing json data: %s", err)
	}

	fmt.Printf("[ OP CODE 11 ] received: %s\n", msg)

	return nil
}

func Identity(connection *websocket.Conn) error {
	authToken := os.Getenv("AUTH_TOKEN")

	op2Identity := opcodes.OP2Identity{
		OP: 2,
		D: opcodes.OP2IdentityData{
			Token:   authToken,
			Intents: 513,
			Properties: opcodes.OP2IdentityProperties{
				OS:      "Linux",
				Browser: "disc-go",
				Device:  "disc-go",
			},
		},
	}

	err := connection.WriteJSON(op2Identity)
	if err != nil {
		return fmt.Errorf("error during writing to websocket: %s", err)
	}

	fmt.Printf("[ OP CODE 02 ] sending Identity payload", utils.GetCurrTimeUTC())

	return nil
}

func Ready(connection *websocket.Conn) error {
	var ready opcodes.OP0Ready

	fmt.Printf("[ OP CODE 0 ] reading message returned from server", utils.GetCurrTimeUTC())
	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("[ OP CODE 0 ] error receiving message: %s", err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&ready); err != nil {
		return fmt.Errorf("[ OP CODE 0 ] error parsing json data: %s", err)
	}

	fmt.Printf("[ OP CODE 0 ] Received: %s", msg)

	return nil
}
