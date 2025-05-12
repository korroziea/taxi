package trip

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/driver-service/internal/domain"
)

type TripRepo interface {
	FindByFreeStatus(ctx context.Context) (domain.AcceptOrderResp, error)
}

type DriverRepo interface {
	UpdateStatus(ctx context.Context, driverID string, status domain.WorkStatus) (domain.Driver, error)
}

type Adapter interface {
	AcceptOrder(ctx context.Context, resp domain.AcceptOrderResp) error
}

type Service struct {
	tripRepo   TripRepo
	driverRepo DriverRepo
	adapter    Adapter
}

func New(
	tripRepo TripRepo,
	driverRepo DriverRepo,
	adapter Adapter,
) *Service {
	service := &Service{
		tripRepo:   tripRepo,
		driverRepo: driverRepo,
		adapter:    adapter,
	}

	return service
}

func (s *Service) AcceptOrder(ctx context.Context, userID string) error {
	resp, err := s.tripRepo.FindByFreeStatus(ctx)
	if err != nil {
		return fmt.Errorf("repo.FindByFreeStatus: %w", err)
	}
	resp.UserID = userID

	_, err = s.driverRepo.UpdateStatus(ctx, resp.Driver.ID, domain.Busy)
	if err != nil {
		return fmt.Errorf("repo.UpdateStatus: %w", err)
	}

	return s.adapter.AcceptOrder(ctx, resp)
}

func (s *Service) CancelTrip(ctx context.Context, driverID string) error {
	_, err := s.driverRepo.UpdateStatus(ctx, driverID, domain.Free)
	if err != nil {
		return fmt.Errorf("repo.UpdateStatus: %w", err)
	}

	return nil
}
