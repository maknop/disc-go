package client

import (
	"context"
	"fmt"
	"os"

	gateway "github.com/maknop/disc-go/gateway"
	utils "github.com/maknop/disc-go/utils"
	"github.com/sirupsen/logrus"
	formatter "github.com/x-cray/logrus-prefixed-formatter"
)

type User struct {
	auth_token string
}

var logger = &logrus.Logger{
	Out:   os.Stderr,
	Level: logrus.DebugLevel,
	Formatter: &formatter.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		ForceFormatting: true,
	},
}

func AuthenticateUser(authToken string) error {
	ctx := context.Background()

	if err := gateway.EstablishConnection(ctx, authToken); err != nil {
		return fmt.Errorf("%s: there was an issue establishing gateway connection: %v", utils.GetCurrTimeUTC(), err)
	}

	return nil
}
