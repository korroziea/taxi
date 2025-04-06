package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, user domain.SignUpUser) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByPhoneAndPassword(ctx context.Context, user domain.SignInUser) (domain.User, error)
}

type Service struct {
	repo Repo
}

func New(repo Repo) *Service {
	service := &Service{
		repo: repo,
	}

	return service
}

func (s *Service) SignUp(ctx context.Context, user domain.SignUpUser) error {
	if err := s.doesUserExist(ctx, user.Phone); err != nil {
		return fmt.Errorf("doesUserExist: %w", err)
	}

	userID, err := domain.GenUserID()
	if err != nil {
		return fmt.Errorf("domain.GenUserID: %w", err)
	}

	user.ID = userID

	// todo: hash password

	_, err = s.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	return nil
}

func (s *Service) doesUserExist(ctx context.Context, phone string) error {
	_, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("repo.FindByPhone: %w", err)
	}

	return domain.ErrUserAlreadyExists
}

// todo: rework
func (s *Service) SignIn(ctx context.Context, user domain.SignInUser) error {
	_, err := s.repo.FindByPhoneAndPassword(ctx, user)
	if err != nil {
		return fmt.Errorf("repo.FindByPhoneAndPassword: %w", err)
	}

	return nil
}
