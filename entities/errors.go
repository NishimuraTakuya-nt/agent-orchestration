package entities

import "errors"

var (
	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserNameRequired  = errors.New("user name is required")
	ErrUserEmailRequired = errors.New("user email is required")
	ErrUserAlreadyExists = errors.New("user already exists")

	// General errors
	ErrInvalidID         = errors.New("invalid ID")
	ErrInternalServer    = errors.New("internal server error")
)