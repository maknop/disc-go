package client

import (
	"fmt"

	utils "github.com/maknop/disc-go/utils"
)

func AuthenticateUser() error {
	if err := utils.LoadEnvVars(); err != nil {
		return fmt.Errorf("%s: there was an error loading .env file: %v", utils.GetCurrTimeUTC(), err)
	}

	return nil
}
