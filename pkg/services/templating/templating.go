package templating

import (
	"embed"
	"html/template"
	"io"
	"path/filepath"

	"github.com/dusnm/slack-ips/pkg/types"
)

type (
	Service struct {
		templateFS     embed.FS
		templateEngine *template.Template
	}
)

func New(
	templateFS embed.FS,
) *Service {
	templateEngine := template.Must(
		template.ParseFS(
			templateFS,
			filepath.Join("templates", "*.html"),
		),
	)

	return &Service{
		templateFS:     templateFS,
		templateEngine: templateEngine,
	}
}

func (s *Service) Render(w io.Writer, page types.Page, data any) error {
	name := page.String()
	t := template.Must(s.templateEngine.Clone())
	t = template.Must(
		t.ParseFS(
			s.templateFS,
			filepath.Join("templates", name),
		),
	)

	return t.ExecuteTemplate(w, name, data)
}
