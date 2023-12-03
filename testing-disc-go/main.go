package main

import (
	"fmt"

	discgo "github.com/maknop/disc-go"
)

func main() {
	fmt.Println("Testing the disc-go library")

	if err := discgo.Start(); err != nil {
		fmt.Print("error connecting to disc-go library")
	}
}
