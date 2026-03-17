package httpserver

import (
	"encoding/json"
	"net/http"
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

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	return err
}
