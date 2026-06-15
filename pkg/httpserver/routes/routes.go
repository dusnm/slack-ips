package routes

import (
	"net/http"

	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes/image"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes/index"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes/settings"
)

func Register(server *httpserver.Server) {
	// Special case for static assets
	http.Handle("GET /assets/", http.FileServer(http.FS(server.DI.AssetsFS)))

	server.Route("get_index", "GET /", index.GET)
	server.Route("post_index", "POST /", index.POST)
	server.Route("get_image", "GET /image", image.GET)
	server.Route("get_settings", "GET /settings", settings.GET)
	server.Route("post_settings", "POST /settings", settings.POST)
	server.Register()
}
