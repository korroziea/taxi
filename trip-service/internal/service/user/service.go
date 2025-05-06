package user

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, trip domain.StartTrip) (domain.Trip, error)
}

type Adapter interface {
	FindDriver(ctx context.Context, req domain.FindDriverReq) error
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

func (s *Service) StartTrip(ctx context.Context, trip domain.StartTrip) error {
	tripID, err := domain.GenTripID()
	if err != nil {
		return fmt.Errorf("domain.GenTripID: %w", err)
	}
	trip.ID = tripID

	_, err = s.repo.Create(ctx, trip)
	if err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}

	req := domain.FindDriverReq{
		UserID: trip.UserID,
	}

	return s.adapter.FindDriver(ctx, req)
}
