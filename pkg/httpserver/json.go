package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/dusnm/slack-ips/pkg/services/templating"
	"github.com/dusnm/slack-ips/pkg/types"
)

func Err(status int, w http.ResponseWriter, err error) error {
	reason := "unknown error"
	if err != nil {
		reason = err.Error()
	}

	return JSON(status, w, map[string]string{"error": reason})
}

func JSON(status int, w http.ResponseWriter, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(body)
	return err
}

func HTML(
	status int,
	w http.ResponseWriter,
	templateService *templating.Service,
	page types.Page,
	data any,
) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	return templateService.Render(w, page, data)
}
