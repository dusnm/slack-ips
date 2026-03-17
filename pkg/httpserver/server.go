package httpserver

import (
	"context"
	"net/http"

	"github.com/dusnm/slack-ips/pkg/container"
	"github.com/rs/zerolog"
)

type (
	HandlerFunc func(
		context.Context,
		*container.Container,
		zerolog.Logger,
		http.ResponseWriter,
		*http.Request,
	) error

	Server struct {
		ctx      context.Context
		di       *container.Container
		logger   zerolog.Logger
		handlers []handler
	}

	handler struct {
		ctx     context.Context
		di      *container.Container
		pattern string
		f       HandlerFunc
		logger  zerolog.Logger
	}
)

func New(
	ctx context.Context,
	di *container.Container,
	logger zerolog.Logger,
) *Server {
	return &Server{
		ctx:    ctx,
		di:     di,
		logger: logger,
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.f(h.ctx, h.di, h.logger, w, r); err != nil {
		h.logger.
			Error().
			Err(err).
			Msg("")
	}
}

func (s *Server) Route(
	name string,
	pattern string,
	f HandlerFunc,
) {
	s.handlers = append(s.handlers, handler{
		ctx:     s.ctx,
		di:      s.di,
		pattern: pattern,
		f:       f,
		logger: s.logger.
			With().
			Str("component", "handler:"+name).
			Logger(),
	})
}

func (s *Server) Register() {
	for _, h := range s.handlers {
		http.Handle(h.pattern, h)
	}
}

func (s *Server) Serve() {
	socket := s.di.GetConfig().App.Socket()
	s.logger.Info().Msgf("starting HTTP server on: http://%s", socket)
	if err := http.ListenAndServe(socket, nil); err != nil {
		s.logger.Fatal().Err(err).Msg("failed to start server")
	}
}
