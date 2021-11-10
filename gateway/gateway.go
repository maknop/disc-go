package gateway

import (
	"net/url"

	websocket "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func Connect() {
	u, err := url.Parse("wss://gateway.discord.gg/?v=9&encoding=json")
	if err != nil {
		log.Error(err)
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("connection successful!")
	}

	// send message
	msg := conn.WriteMessage(websocket.TextMessage, []byte("Message sent successful!"))
	if err != nil {
		log.Error(err)
	} else {
		log.Info(msg)
	}

	// receive message
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Error(err)
	} else {
		log.Info(message)
	}
}
