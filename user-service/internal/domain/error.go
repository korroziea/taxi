package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("entuty not found")

	ErrInvalidHashFormat  = errors.New("invalid hash format")
	ErrInvalidHashType    = errors.New("invalid hash type")
	ErrInvalidHashVersion = errors.New("invalid hash version")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")
)
