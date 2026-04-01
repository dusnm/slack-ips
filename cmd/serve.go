package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes"
)

func Serve(c *container.Container) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	server := httpserver.New(
		ctx,
		c,
		c.GetLogger().
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
