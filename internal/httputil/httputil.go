package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(payload)
}

func WriteError(writer http.ResponseWriter, status int, code, message string, details ...map[string]any) {
	body := map[string]any{"ok": false, "code": code, "message": message}
	if len(details) > 0 {
		for key, value := range details[0] {
			body[key] = value
		}
	}
	WriteJSON(writer, status, body)
}

func DecodeJSON(writer http.ResponseWriter, request *http.Request, target any) error {
	if err := json.NewDecoder(request.Body).Decode(target); err != nil {
		WriteError(writer, http.StatusBadRequest, "invalid_body", err.Error())
		return errors.New("invalid body")
	}
	return nil
}
