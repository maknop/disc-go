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
		return fmt.Errorf("%s: connection could not be established: %s", utils.GetCurrTimeUTC(), err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	fmt.Printf("%s: [ OP CODE 10 ] sending initial request to Discord Gateway server", utils.GetCurrTimeUTC())
	heartbeatInterval, sequenceNum, err := ReceiveMessage(connection)
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 10 ] could not retrieve heartbeat interval: %s", utils.GetCurrTimeUTC(), err)
	}

	fmt.Printf("%s: [ OP CODE 10 ] value of heartbeat interval is: %d seconds", utils.GetCurrTimeUTC(), (heartbeatInterval / 1000))

	// OP 1 Heartbeat
	if err := SendMessage(connection, sequenceNum); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s: error during writing to websocket: %s", utils.GetCurrTimeUTC(), err))
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
		return 0, nil, fmt.Errorf("%s: [ OP CODE 10 ] error receiving message: %s", utils.GetCurrTimeUTC(), err)
	}

	fmt.Printf("%s: [ OP CODE 10 ] received: %s\n", utils.GetCurrTimeUTC(), msg)

	var OP10Hello opcodes.OP10Hello

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&OP10Hello); err != nil {
		return 0, nil, fmt.Errorf("%s: [ OP CODE 10 ] error parsing json data: %s", utils.GetCurrTimeUTC(), err)
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

	fmt.Printf("%s: [ OP CODE 01 ] sending the following payload: {op: %d, d: {Seq: %v}}", utils.GetCurrTimeUTC(), op1Heartbeat.OP, op1Heartbeat.D.Sequence)

	return nil
}

func ACK(connection *websocket.Conn) error {
	var op11ACK opcodes.OP11HeartbeatACK

	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 11 ] error receiving message: %s", utils.GetCurrTimeUTC(), err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op11ACK); err != nil {
		return fmt.Errorf("%s: [OP CODE 11 ] error parsing json data: %s", utils.GetCurrTimeUTC(), err)
	}

	fmt.Printf("%s: [ OP CODE 11 ] received: %s\n", utils.GetCurrTimeUTC(), msg)

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
		return fmt.Errorf("%s: error during writing to websocket: %s", utils.GetCurrTimeUTC(), err)
	}

	fmt.Printf("%s: [ OP CODE 02 ] sending Identity payload", utils.GetCurrTimeUTC())

	return nil
}

func Ready(connection *websocket.Conn) error {
	var ready opcodes.OP0Ready

	fmt.Printf("%s: [ OP CODE 0 ] reading message returned from server", utils.GetCurrTimeUTC())
	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 0 ] error receiving message: %s", utils.GetCurrTimeUTC(), err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&ready); err != nil {
		return fmt.Errorf("%s: [ OP CODE 0 ] error parsing json data: %s", utils.GetCurrTimeUTC(), err)
	}

	fmt.Printf("%s: [ OP CODE 0 ] Received: %s", msg, utils.GetCurrTimeUTC())

	return nil
}
