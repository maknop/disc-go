package client

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/maknop/disc-go/config"
	"github.com/maknop/disc-go/gateway"
)

func Start() error {
	ctx := context.Background()

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true

	log.Info("attempting to authenticate user")
	if err := AuthenticateUser(ctx); err != nil {
		return fmt.Errorf("failed to authenticate user")
	}

	return nil
}

func AuthenticateUser(ctx context.Context) error {
	log.Info("Loading environment variables")
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("there was an error loading .env file: %v", err)
	}
	log.Info(".env file successfully loaded")

	if err := gateway.EstablishConnection(ctx); err != nil {
		return fmt.Errorf("there was an issue establishing gateway connection: %v", err)
	}

	log.Info("authentication was successful to the gateway")

	return nil
}
