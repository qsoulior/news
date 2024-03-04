package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DecodeJSON[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()
	data := new(T)
	d := json.NewDecoder(r.Body)

	err := d.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("d.Decode: %w", err)
	}

	return data, nil
}

func EncodeJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := json.NewEncoder(w)
	e.Encode(data)
}
