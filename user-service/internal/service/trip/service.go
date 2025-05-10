package trip

import (
	"context"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type Repo interface {
	CheckWalletBalance(ctx context.Context, cost int64) error
}

type HTTPAdapter interface {
	Trips(ctx context.Context, userID string) ([]domain.Trip, error)
}

type Adapter interface {
	StartTrip(ctx context.Context, trip domain.StartTrip) error
	CancelTrip(ctx context.Context) error
}

type Service struct {
	// repo    Repo
	httpAdapter HTTPAdapter
	adapter     Adapter
}

func New(adapter Adapter, httpAdapter HTTPAdapter) *Service {
	service := &Service{
		// repo:    repo,
		httpAdapter: httpAdapter,
		adapter:     adapter,
	}

	return service
}

func (s *Service) StartTrip(ctx context.Context, trip domain.StartTrip) error {
	// if err := s.repo.CheckWalletBalance(ctx, 0); err != nil { // todo: add cost
	// 	return fmt.Errorf("repo.CheckWalletBalance: %w", err)
	// }

	return s.adapter.StartTrip(ctx, trip)
}

func (s *Service) CancelTrip(ctx context.Context) error {
	return s.adapter.CancelTrip(ctx)
}

func (s *Service) Trips(ctx context.Context, userID string) ([]domain.Trip, error) {
	return s.httpAdapter.Trips(ctx, userID)
}

func (s *Service) Cost(ctx context.Context) (int64, error) {
	return 0, nil
}
