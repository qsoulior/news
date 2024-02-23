package service

import (
	"errors"
	"fmt"
)

var (
	ErrRateLimit      = errors.New("rate limit exceeded")
	ErrRequestInvalid = errors.New("request is invalid")
	ErrInternalServer = errors.New("API internal error")
	ErrUnexpectedCode = errors.New("http code is unexcepted")
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
