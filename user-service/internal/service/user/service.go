package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type Hasher interface {
	Generate(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type Repo interface {
	Create(ctx context.Context, user domain.SignUpUser) (domain.User, error)
	UpdateProfile(ctx context.Context, user domain.ProfileUser) (domain.User, error)
	FindByID(ctx context.Context, id string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByPhoneAndPassword(ctx context.Context, user domain.SignInUser) (domain.User, error)
}

type Service struct {
	hasher Hasher
	repo   Repo
}

func New(hasher Hasher, repo Repo) *Service {
	service := &Service{
		hasher: hasher,
		repo:   repo,
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

	hash, err := s.hasher.Generate(user.Password)
	if err != nil {
		return fmt.Errorf("hasher.Generate: %w", err)
	}
	user.Password = hash

	_, err = s.repo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	return nil
}

func (s *Service) doesUserExist(ctx context.Context, phone string) error {
	_, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil
		}

		return fmt.Errorf("repo.FindByPhone: %w", err)
	}

	return domain.ErrUserAlreadyExists
}

func (s *Service) SignIn(ctx context.Context, user domain.SignInUser) (string, error) {
	foundUser, err := s.repo.FindByPhone(ctx, user.Phone)
	if err != nil {
		return "", fmt.Errorf("repo.FindByPhoneAndPassword: %w", err)
	}

	verified, err := s.hasher.Verify(user.Password, foundUser.Password)
	if err != nil {
		return "", fmt.Errorf("hasher.Verify: %w", err)
	}
	if verified != true {
		return "", domain.ErrWrongPassword
	}

	return foundUser.ID, nil
}

func (s *Service) Profile(ctx context.Context, user domain.ProfileUser) error {
	foundUser, err := s.repo.FindByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("repo.FindByID: %w", err)
	}

	verified, err := s.hasher.Verify(user.Password, foundUser.Password)
	if err != nil {
		return fmt.Errorf("hasher.Verify: %w", err)
	}
	if verified != true {
		return domain.ErrWrongPassword
	}

	_, err = s.repo.UpdateProfile(ctx, user)
	if err != nil {
		return fmt.Errorf("repo.UpdateProfile: %w", err)
	}

	return nil
}
