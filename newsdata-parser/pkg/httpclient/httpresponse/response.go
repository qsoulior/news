package httpresponse

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func JSON[T any](r *http.Response) (*T, error) {
	defer r.Body.Close()
	data := new(T)
	d := json.NewDecoder(r.Body)

	err := d.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("d.Decode: %w", err)
	}

	return data, nil
}

func Bytes(r *http.Response) ([]byte, error) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	return data, nil
}

func String(r *http.Response) (string, error) {
	data, err := Bytes(r)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
