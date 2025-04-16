package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/korroziea/taxi/driver-service/internal/domain"
)

type Hasher interface {
	Generate(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type Repo interface {
	Create(ctx context.Context, user domain.SignUpDriver) (domain.Driver, error)
	FindByPhone(ctx context.Context, phone string) (domain.Driver, error)
	FindByPhoneAndPassword(ctx context.Context, user domain.SignInDriver) (domain.Driver, error)
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

func (s *Service) SignUp(ctx context.Context, driver domain.SignUpDriver) error {
	if err := s.doesDriverExist(ctx, driver.Phone); err != nil {
		return fmt.Errorf("doesDriverExist: %w", err)
	}

	driverID, err := domain.GenDriverID()
	if err != nil {
		return fmt.Errorf("domain.GenUserID: %w", err)
	}
	driver.ID = driverID

	hash, err := s.hasher.Generate(driver.Password)
	if err != nil {
		return fmt.Errorf("hasher.Generate: %w", err)
	}
	driver.Password = hash

	_, err = s.repo.Create(ctx, driver)
	if err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	return nil
}

func (s *Service) doesDriverExist(ctx context.Context, phone string) error {
	_, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		if errors.Is(err, domain.ErrDriverNotFound) {
			return nil
		}

		return fmt.Errorf("repo.FindByPhone: %w", err)
	}

	return domain.ErrDriverAlreadyExists
}

func (s *Service) SignIn(ctx context.Context, driver domain.SignInDriver) (string, error) {
	foundDriver, err := s.repo.FindByPhone(ctx, driver.Phone)
	if err != nil {
		return "", fmt.Errorf("repo.FindByPhoneAndPassword: %w", err)
	}

	verified, err := s.hasher.Verify(driver.Password, foundDriver.Password)
	if err != nil {
		return "", fmt.Errorf("hasher.Verify: %w", err)
	}
	if verified != true {
		return "", domain.ErrWrongPassword
	}

	return foundDriver.ID, nil
}
