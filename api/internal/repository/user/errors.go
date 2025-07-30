package user

import "errors"

var (
	// ErrNotFound means the item was not found
	ErrNotFound = errors.New("not found")
)
