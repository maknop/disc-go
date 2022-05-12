package main

import (
	"context"

	"github.com/maknop/disc-go/src/gateway"
)

func main() {
	ctx := context.Background()

	gateway.EstablishConnection(ctx)
}
