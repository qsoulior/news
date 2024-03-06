package httprequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func JSON(data any) (io.Reader, error) {
	buf := new(bytes.Buffer)
	e := json.NewEncoder(buf)

	err := e.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("e.Encode: %w", err)
	}

	return buf, nil
}

func Bytes(data []byte) io.Reader {
	return bytes.NewReader(data)
}

func String(data string) io.Reader {
	return strings.NewReader(data)
}
