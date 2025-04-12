package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("entuty not found")

	ErrInvalidHashFormat  = errors.New("invalid hash format")
	ErrInvalidHashType    = errors.New("invalid hash type")
	ErrInvalidHashVersion = errors.New("invalid hash version")

	ErrParseBearerToken = errors.New("invalid token")
	ErrParseJWTToken    = errors.New("can't parse jwt token")
	ErrParseClaims      = errors.New("can't parse token claims")
	ErrFetchSub         = errors.New("can't fetch sub")
	ErrTokenNotFound    = errors.New("token not found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")

	ErrWalletNotFound   = errors.New("wallet not found")
	ErrChangeWalletType = errors.New("can't be converted to person type")
)
