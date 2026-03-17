package routes

import (
	"github.com/dusnm/slack-ips/pkg/httpserver"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes/image"
	"github.com/dusnm/slack-ips/pkg/httpserver/routes/index"
)

func Register(server *httpserver.Server) {
	server.Route("get_index", "GET /", index.GET)
	server.Route("post_index", "POST /", index.POST)
	server.Route("get_image", "GET /image", image.GET)
	server.Register()
}
