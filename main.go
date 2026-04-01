package main

import (
	"os"

	"github.com/dusnm/slack-ips/cmd"
	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stderr)
	di := container.New(
		logger.
			With().
			Str("component", "container").
			Logger(),
	)

	defer logger.Info().Msg("SHUTDOWN OK, GOODBYE")
	defer di.Close()

	cmd.Run(di)
}
