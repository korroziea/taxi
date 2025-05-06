package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("entuty not found")

	ErrInvalidHashFormat  = errors.New("invalid hash format")
	ErrInvalidHashType    = errors.New("invalid hash type")
	ErrInvalidHashVersion = errors.New("invalid hash version")
	ErrWrongPassword      = errors.New("wrong password")

	ErrParseBearerToken = errors.New("invalid token")
	ErrParseJWTToken    = errors.New("can't parse jwt token")
	ErrParseClaims      = errors.New("can't parse token claims")
	ErrFetchSub         = errors.New("can't fetch sub")
	ErrTokenNotFound    = errors.New("token not found")

	ErrDriverAlreadyExists = errors.New("driver already exists")
	ErrDriverNotFound      = errors.New("driver not found")

	ErrDriverCarNotFound = errors.New("driver or car not found")
)
