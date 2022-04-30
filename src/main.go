package main

import (
	"fmt"

	"github.com/disc-go/src/gateway"
)

func main() {
	fmt.Println("disc-go!")

	gateway.EstablishConnection()
}
