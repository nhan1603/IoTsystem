package auth

import "errors"

var (
	// ErrUserNotFound means user not found
	ErrUserNotFound = errors.New("user not found")
)
