package client

import (
	"context"
	"fmt"

	gateway "github.com/maknop/disc-go/gateway"
)

type User struct {
	auth_token string
}

func AuthenticateUser(authToken string) error {
	ctx := context.Background()

	if err := gateway.EstablishConnection(ctx, authToken); err != nil {
		return fmt.Errorf("there was an issue establishing gateway connection: %v", err)
	}

	return nil
}
