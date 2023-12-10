package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func TerminateGracefully() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("EXITING: Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func GetCurrTimeUTC() string {
	return time.Now().Format(time.Kitchen)
}
