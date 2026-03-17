package index

import (
	"context"
	"net/http"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/rs/zerolog"
)

func GET(
	_ context.Context,
	_ *container.Container,
	_ zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if r.URL.Path != "/" {
		return httpserver.Err(http.StatusNotFound, w, httpserver.ErrNotFound)
	}

	return httpserver.JSON(http.StatusOK, w, map[string]string{"status": http.StatusText(http.StatusOK)})
}
