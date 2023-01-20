package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	logrus "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
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
	origin := "http://localhost/"
	url := wsUrl
	c, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	//c.Conn.EnableWriteCompression(true)
	//c.Conn.SetCompressionLevel(1)

	return nil
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

	var DiscordGateway Client
	if err = json.Unmarshal(respBody, &DiscordGateway); err != nil {
		return "", fmt.Errorf("error retrieving gateway URL json: %s", err)
	}

	logrus.WithFields(logrus.Fields{"gateway": DiscordGateway.Url}).Info("successfully received gateway url from Discord server")

	return DiscordGateway.Url, nil
}
