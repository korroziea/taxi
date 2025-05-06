package driver

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type Repo interface {
	UpdateDriverInfo(ctx context.Context, req domain.AcceptOrderReq) (domain.Trip, error)
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

func (s *Service) AcceptOrder(ctx context.Context, req domain.AcceptOrderReq) error {
	_, err := s.repo.UpdateDriverInfo(ctx, req)
	if err != nil {
		return fmt.Errorf("repo.UpdateDriverInfo: %w", err)
	}
	
	return nil
}
