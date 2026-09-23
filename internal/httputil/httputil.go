package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
)

const maxJSONBodyBytes = 4 * 1024 * 1024

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
	request.Body = http.MaxBytesReader(writer, request.Body, maxJSONBodyBytes)
	if err := json.NewDecoder(request.Body).Decode(target); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			WriteError(writer, http.StatusRequestEntityTooLarge, "body_too_large", "JSON request body is too large")
			return errors.New("body too large")
		}
		WriteError(writer, http.StatusBadRequest, "invalid_body", err.Error())
		return errors.New("invalid body")
	}
	return nil
}
