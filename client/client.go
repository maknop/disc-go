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

func AuthenticateUser(authToken string) error {
	ctx := context.Background()

	if err := gateway.EstablishConnection(ctx, authToken); err != nil {
		return fmt.Errorf("%s: there was an issue establishing gateway connection: %v", utils.GetCurrTimeUTC(), err)
	}

	return nil
}
