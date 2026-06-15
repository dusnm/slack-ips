package main

import (
	"embed"
	"os"

	"github.com/dusnm/slack-ips/cmd"
	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/rs/zerolog"
)

var (
	//go:embed assets/*
	assetsFS embed.FS

	//go:embed templates/*
	templatesFS embed.FS
)

func main() {
	logger := zerolog.New(os.Stderr)
	di := container.New(
		assetsFS,
		templatesFS,
		logger.
			With().
			Str("component", "container").
			Logger(),
	)

	defer logger.Info().Msg("SHUTDOWN OK, GOODBYE")
	defer di.Close()

	cmd.Run(di)
}
