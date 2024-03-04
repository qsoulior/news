package handler

import "net/http"

type JSONError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func ErrorJSON(w http.ResponseWriter, error string, code int) {
	EncodeJSON(w, &JSONError{
		Status: http.StatusText(code),
		Error:  error,
	}, code)
}
