package driver

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type Repo interface {
	UpdateDriverInfo(ctx context.Context, req domain.AcceptOrderReq) (domain.Trip, error)
}

type Adapter interface {
	AcceptTrip(ctx context.Context, trip domain.Trip) error
}

type Service struct {
	repo    Repo
	adapter Adapter
}

func New(repo Repo, adapter Adapter) *Service {
	service := &Service{
		repo:    repo,
		adapter: adapter,
	}

	return service
}

func (s *Service) AcceptOrder(ctx context.Context, req domain.AcceptOrderReq) error {
	trip, err := s.repo.UpdateDriverInfo(ctx, req)
	if err != nil {
		return fmt.Errorf("repo.UpdateDriverInfo: %w", err)
	}

	return s.adapter.AcceptTrip(ctx, trip)
}
