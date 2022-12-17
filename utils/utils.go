package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnvVars() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("%s: failed to load .env file: %s", GetCurrTimeUTC(), err)
	}

	fmt.Printf("%s: .env file successfully loaded", GetCurrTimeUTC())

	return nil
}

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
