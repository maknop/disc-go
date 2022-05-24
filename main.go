package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

func TerminateGracefully() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("EXITING: Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
