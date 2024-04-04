package rssclient

import (
	"fmt"
	"net/http"
)

type statusError struct {
	Code int
	Text string
}

func newStatusError(code int) *statusError {
	return &statusError{code, http.StatusText(code)}
}

func (e *statusError) Error() string {
	return fmt.Sprintf("incorrect response status: %d - %s", e.Code, e.Text)
}
