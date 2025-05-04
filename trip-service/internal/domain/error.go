package domain

import "errors"

var (
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("entuty not found")

	ErrTripNotFound = errors.New("trip not found")
)
