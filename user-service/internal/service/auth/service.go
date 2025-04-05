package auth

import (
	"context"
)

type Service struct{}

func New() *Service {
	service := &Service{}

	return service
}

func (s *Service) SignUp(ctx context.Context) error {
	return nil
}

func (s *Service) SignIn(ctx context.Context) error {
	return nil
}
