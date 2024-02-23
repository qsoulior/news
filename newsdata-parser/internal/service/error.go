package service

import (
	"errors"
	"fmt"

	"github.com/qsoulior/news/newsdata-parser/internal/repo"
)

var (
	ErrRateLimit      = errors.New("rate limit exceeded")
	ErrRequestInvalid = errors.New("request is invalid")
	ErrInternalServer = errors.New("API internal error")
	ErrUnexpectedCode = errors.New("http code is unexcepted")
	ErrNotExist       = repo.ErrNotExist
)

type ResponseError struct {
	Err  error
	Code string
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err, e.Code)
}

func (e *ResponseError) Unwrap() error {
	return e.Err
}
