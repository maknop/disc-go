package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	opcodes "github.com/maknop/disc-go/types"
	utils "github.com/maknop/disc-go/utils"

	"github.com/gorilla/websocket"
	logrus "github.com/sirupsen/logrus"
)

type DiscordGateway struct {
	Url string `json:"url"`
}

func getGatewayUrl() (string, error) {
	url := "https://discordapp.com/api/gateway"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not retrieve gateway url: %s", err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not retrieve response body from server: %s", err)
	}

	var DiscordGateway DiscordGateway
	if err = json.Unmarshal(respBody, &DiscordGateway); err != nil {
		return "", fmt.Errorf("error retrieving gateway URL json: %s", err)
	}

	logrus.WithFields(logrus.Fields{
		"gateway": DiscordGateway.Url,
	}).Info("successfully received gateway url from Discord server")

	return DiscordGateway.Url, nil
}

func Connect(ctx context.Context, authToken string) error {
	gatewayUrl, err := getGatewayUrl()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"op_code": 10,
		}).Info("error retrieving gateway url")
	}

	connection, _, err := websocket.DefaultDialer.DialContext(ctx, gatewayUrl, nil)
	if err != nil {
		return fmt.Errorf("%s: connection could not be established: %s", utils.GetCurrTimeUTC(), err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	logrus.WithFields(logrus.Fields{
		"op_code": 10,
	}).Info("sending initial request to gateway")

	heartbeat_interval, sequence_num, err := ReceiveHelloEvent(connection)
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 10 ] could not retrieve heartbeat interval: %s", utils.GetCurrTimeUTC(), err)
	}

	logrus.WithFields(logrus.Fields{
		"op_code": 10,
	}).Infof("received Hello event from gateway")

	go HeatbeatInterval(connection, heartbeat_interval, sequence_num)

	// OP 2 Identity
	if err := Identity(connection, authToken); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s: error during Identity (OP CODE 02): %s", utils.GetCurrTimeUTC(), err))

	}

	if err := Ready(connection); err != nil {
		return fmt.Errorf(fmt.Sprintf("%s: error during Ready: %s", utils.GetCurrTimeUTC(), err))
	}

	return nil
}

func HeatbeatInterval(connection *websocket.Conn, heartbeat_interval int, sequence_num *int) {
	for {
		time.Sleep(time.Duration(heartbeat_interval))

		logrus.WithFields(logrus.Fields{
			"op_code": 1,
		}).Info("Sending heartbeat event to gateway")

		if err := SendHeartbeatEvent(connection, sequence_num); err != nil {
			fmt.Errorf(fmt.Sprintf("%s: error occurred sending heartbeat event: %s", utils.GetCurrTimeUTC(), err))
		}

		logrus.WithFields(logrus.Fields{
			"op_code": 1,
		}).Info("loop")

		if err := ReceiveHeartbeatACKEvent(connection); err != nil {
			fmt.Errorf(fmt.Sprintf("%s: error occurred receiving heartbeat ACK event: %s", utils.GetCurrTimeUTC(), err))
		}

		logrus.WithFields(logrus.Fields{
			"op_code": 1,
		}).Info("Received heartbeat ACK event from gateway")
	}
}

func ReceiveHelloEvent(connection *websocket.Conn) (int, *int, error) {
	_, msg, err := connection.ReadMessage()
	if err != nil {
		return 0, nil, fmt.Errorf("%s: [ OP CODE 10 ] error receiving message: %s", utils.GetCurrTimeUTC(), err)
	}

	var op_10_hello opcodes.OP_10_Hello
	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_10_hello); err != nil {
		return 0, nil, fmt.Errorf("%s: [ OP CODE 10 ] error parsing json data: %s", utils.GetCurrTimeUTC(), err)
	}

	return op_10_hello.D.Heartbeat_Interval, op_10_hello.S, nil
}

func SendHeartbeatEvent(connection *websocket.Conn, sequence_num *int) error {
	op_1_heartbeat := opcodes.OP_1_Heartbeat{
		OP: 1,
		D:  opcodes.OP_1_Heartbeat_Data{Sequence: sequence_num},
	}

	err := connection.WriteJSON(op_1_heartbeat)
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 01 ] error sending message", utils.GetCurrTimeUTC())
	}

	return nil
}

func ReceiveHeartbeatACKEvent(connection *websocket.Conn) error {
	var op_11_ack opcodes.OP_11_Heartbeat_ACK

	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 11 ] error receiving message: %s", utils.GetCurrTimeUTC(), err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&op_11_ack); err != nil {
		return fmt.Errorf("%s: [OP CODE 11 ] error parsing json data: %s", utils.GetCurrTimeUTC(), err)
	}

	return nil
}

func Identity(connection *websocket.Conn, auth_token string) error {
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
		return fmt.Errorf("%s: error during writing to websocket: %s", utils.GetCurrTimeUTC(), err)
	}

	logrus.WithFields(logrus.Fields{
		"op_code": 2,
	}).Info("sending Identity payload")

	return nil
}

func Ready(connection *websocket.Conn) error {
	var ready opcodes.OP_0_Ready

	logrus.WithFields(logrus.Fields{
		"op_code": 0,
	}).Info("reading message returned from server")

	_, msg, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("%s: [ OP CODE 0 ] error receiving message: %s", utils.GetCurrTimeUTC(), err)
	}

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&ready); err != nil {
		return fmt.Errorf("%s: [ OP CODE 0 ] error parsing json data: %s", utils.GetCurrTimeUTC(), err)
	}

	return nil
}
