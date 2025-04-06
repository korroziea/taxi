package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("entuty not found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)
