package main

import (
	"context"
	"fmt"

	"github.com/disc-go/src/gateway"
)

func main() {
	ctx := context.Background()
	fmt.Println("disc-go!")

	gateway.EstablishConnection(ctx)
}
