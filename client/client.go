package client

import (
	"context"
	"fmt"

	gateway "github.com/maknop/disc-go/gateway"
	utils "github.com/maknop/disc-go/utils"
)

type User struct {
	auth_token string
}

func AuthenticateUser() error {
	ctx := context.Background()

	if err := gateway.EstablishConnection(ctx); err != nil {
		return fmt.Errorf("%s: there was an issue establishing gateway connection: %v", utils.GetCurrTimeUTC(), err)
	}

	fmt.Println("Did the thing")

	return nil
}
