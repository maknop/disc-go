package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	opcodes "github.com/maknop/disc-go/types"
)

var (
	socketURL = "wss://gateway.discord.gg/?v=9&encoding=json&compress?=true"
)

type client struct {
	ctx               context.Context
	connection        *websocket.Conn
	heartbeatInterval int
	p                 int
}

func EstablishConnection(ctx context.Context) error {
	log.Info("Establishing connection...")
	connection, _, err := websocket.DefaultDialer.DialContext(ctx, socketURL, nil)
	if err != nil {
		return fmt.Errorf("connection could not be established: %s", err)
	}

	connection.EnableWriteCompression(true)
	connection.SetCompressionLevel(1)

	cl := client{ctx: ctx, connection: connection}

	heartbeatInterval, err := cl.Read(connection)
	if err != nil {
		return fmt.Errorf("[ OP CODE 10 ] could not retrieve heartbeat interval: %s", err)
	}

	cl.heartbeatInterval = heartbeatInterval

	go cl.HeartbeatEvent(connection, heartbeatInterval)
	return nil
}

func (cl *client) HeartbeatEvent(connection *websocket.Conn, heartbeat int) {
	jitter := rand.Float32()
	waitTime := jitter * float32(heartbeat)

	byteMessage, err := json.Marshal(opcodes.OP1Heartbeat{OP: 1, D: nil})
	if err != nil {
		fmt.Errorf("Error marshalling data to byte slice")
	}

	connection.WriteMessage(cl.p, byteMessage)

	for {
		byteMessage, err := json.Marshal(opcodes.OP1Heartbeat{OP: 1, D: nil})
		if err != nil {
			fmt.Errorf("Error marshalling data to byte slice")
		}

		time.Sleep(time.Duration(waitTime))

		connection.WriteMessage(cl.p, byteMessage)
	}
}

func (cl *client) Read(connection *websocket.Conn) (int, error) {
	p, msg, err := connection.ReadMessage()
	if err != nil {
		return 0, fmt.Errorf("[ OP CODE 10 ] error receiving message: %s", err)
	}

	cl.p = p

	log.Info(fmt.Sprintf("[ OP CODE 10 ] received: %s", msg))

	var OP10Hello opcodes.OP10Hello

	if err := json.NewDecoder(bytes.NewReader(msg)).Decode(&OP10Hello); err != nil {
		return 0, fmt.Errorf("[ OP CODE 10 ] error parsing json data: %s", err)
	}

	return OP10Hello.D.HeartbeatInterval, nil
}

func (cl *client) getGatewayURL() {

}
