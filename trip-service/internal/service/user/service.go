package user

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, trip domain.StartTrip) (domain.Trip, error)
	UpdateToCanceledTrip(ctx context.Context, userID string) (domain.Trip, error)
	FindTrips(ctx context.Context, userID string) ([]domain.Trip, error)
}

type DriverAdapter interface {
	CancelTrip(ctx context.Context, driverID string) error
	FindDriver(ctx context.Context, req domain.FindDriverReq) error
}

type UserAdapter interface {
	Trips(ctx context.Context, trips []domain.Trip) error
}

type Service struct {
	repo          Repo
	driverAdapter DriverAdapter
	userAdapter   UserAdapter
}

func New(repo Repo, driverAdapter DriverAdapter, userAdapter UserAdapter) *Service {
	service := &Service{
		repo:          repo,
		driverAdapter: driverAdapter,
		userAdapter:   userAdapter,
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

	return s.driverAdapter.FindDriver(ctx, req)
}

func (s *Service) CancelTrip(ctx context.Context, userID string) error {
	trip, err := s.repo.UpdateToCanceledTrip(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo.UpdateToCanceledTrip: %w", err)
	}

	return s.driverAdapter.CancelTrip(ctx, trip.DriverID)
}

func (s *Service) Trips(ctx context.Context, userID string) ([]domain.Trip, error) {
	trips, err := s.repo.FindTrips(ctx, userID)
	if err != nil {
		return []domain.Trip{}, fmt.Errorf("repo.FindTrips: %w", err)
	}

	return trips, nil
}
