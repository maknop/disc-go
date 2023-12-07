package client

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/maknop/disc-go/gateway"
	"github.com/maknop/disc-go/utils"
)

func Start() {
	fmt.Print("attempting to authenticate user")

	if err := AuthenticateUser(); err != nil {
		log.Fatal("failed to authenticate user")
	}
}

func AuthenticateUser() error {
	ctx := context.Background()

	if err := utils.LoadEnvVars(); err != nil {
		return fmt.Errorf("%s: there was an error loading .env file: %v", utils.GetCurrTimeUTC(), err)
	}

	if err := gateway.EstablishConnection(ctx); err != nil {
		return fmt.Errorf("%s: there was an issue establishing gateway connection: %v", utils.GetCurrTimeUTC(), err)
	}

	fmt.Println("Did the thing")

	return nil
}
