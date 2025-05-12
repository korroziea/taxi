package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/korroziea/taxi/driver-service/internal/domain"
	"github.com/korroziea/taxi/driver-service/internal/handler/driver"
)

type Hasher interface {
	Generate(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type Repo interface {
	Create(ctx context.Context, driver domain.SignUpDriver) (domain.Driver, error)
	UpdateStatus(ctx context.Context, driverID string, status domain.WorkStatus) (domain.Driver, error)
	FindByID(ctx context.Context, driverID string) (domain.Driver, error)
	FindByPhone(ctx context.Context, phone string) (domain.Driver, error)
	FindByPhoneAndPassword(ctx context.Context, driver domain.SignInDriver) (domain.Driver, error)
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

func (s *Service) Shift(ctx context.Context) error {
	driverID := driver.FromContext(ctx)
	driver, err := s.repo.FindByID(ctx, driverID)
	if err != nil {
		return fmt.Errorf("repo.FindByID: %w", err)
	}

	workStatus := domain.OffShift
	if driver.Status == domain.OffShift {
		workStatus = domain.Free
	}

	_, err = s.repo.UpdateStatus(ctx, driverID, workStatus)
	if err != nil {
		return fmt.Errorf("repo.UpdateStatus: %w", err)
	}

	return nil
}
