package user

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, trip domain.StartTrip) (domain.Trip, error)
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

	return nil
}
