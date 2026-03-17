package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger := zerolog.New(os.Stdout)
	di := container.New(
		logger.
			With().
			Str("component", "container").
			Logger(),
	)

	defer logger.Info().Msg("SHUTDOWN OK, GOODBYE")
	defer di.Close()

	server := httpserver.New(
		ctx,
		di,
		logger.
			With().
			Str("component", "httpserver").
			Logger(),
	)

	routes.Register(server)

	go func() {
		server.Serve()
	}()

	<-ctx.Done()
}
