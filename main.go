package main

import (
	"context"

	"github.com/maknop/disc-go/gateway"
)

func init() {
	TerminateGracefully()
	LoadEnvVars()
}

func main() {
	ctx := context.Background()

	gateway.EstablishConnection(ctx)
}
