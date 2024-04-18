package repo

import "errors"

var (
	ErrInvalidID = errors.New("invalid ObjectID")
	ErrNotFound  = errors.New("entity not found")
)
