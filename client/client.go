package client

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/maknop/disc-go/gateway"
	"github.com/maknop/disc-go/utils"
)

func Start() error {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	log.Info("attempting to authenticate user")
	if err := AuthenticateUser(); err != nil {
		return fmt.Errorf("failed to authenticate user")
	}

	return nil
}

func AuthenticateUser() error {
	ctx := context.Background()

	if err := utils.LoadEnvVars(); err != nil {
		log.Fatalf("there was an error loading .env file: %v", err)
	}

	if err := gateway.EstablishConnection(ctx); err != nil {
		return fmt.Errorf("there was an issue establishing gateway connection: %v", err)
	}

	log.Info("authentication was successful to the gateway")

	return nil
}
