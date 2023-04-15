package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	logrus "github.com/sirupsen/logrus"

	opcodes "github.com/maknop/disc-go/types"
)

type Client struct {
	Url  string `json:"url"`
	Conn *websocket.Conn
	Send chan []byte
}

func Connect(ctx context.Context, authToken string) error {
	wsUrl, err := getGatewayUrl()
	if err != nil {
		logrus.WithFields(logrus.Fields{"op_code": 10}).Info("error retrieving gateway url")
	}

	logrus.WithFields(logrus.Fields{"op_code": 10}).Info("sending initial request to server")
	connection, _, err := websocket.DefaultDialer.Dial(wsUrl, nil)
	if err != nil {
		return fmt.Errorf("An error occurred dialing the websocket server: %w", err)
	}
	defer connection.Close()

	// channel := make(chan []byte)

	// c := Client{
	// 	Url:  wsUrl,
	// 	Conn: connection,
	// 	Send: channel,
	// }

	_, p, err := connection.ReadMessage()
	if err != nil {
		return fmt.Errorf("There was an error reading the message: %w", err)
	}

	var HelloEvent opcodes.OP_10_Hello
	if err = json.Unmarshal(p, &HelloEvent); err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	fmt.Printf("Message received: %v", HelloEvent.OP)

	return nil
}

func getGatewayUrl() (string, error) {
	url := "https://discordapp.com/api/gateway"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not retrieve gateway url: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not retrieve response body from server: %s", err)
	}

	var DiscordGateway Client
	if err = json.Unmarshal(respBody, &DiscordGateway); err != nil {
		return "", fmt.Errorf("error retrieving gateway URL json: %w", err)
	}

	logrus.WithFields(logrus.Fields{"gateway": DiscordGateway.Url}).Info("successfully received gateway url from Discord server")

	return DiscordGateway.Url, nil
}

// func startHeartbeatInterval(heatbeatValue int) {
// 	while (true) {

// 	}
// }
