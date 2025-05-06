package trip

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/driver-service/internal/domain"
)

type Repo interface {
	FindByFreeStatus(ctx context.Context) (domain.AcceptOrderResp, error)
}

type Adapter interface {
	AcceptOrder(ctx context.Context, resp domain.AcceptOrderResp) error
}

type Service struct {
	repo    Repo
	adapter Adapter
}

func New(
	repo Repo,
	adapter Adapter,
) *Service {
	service := &Service{
		repo:    repo,
		adapter: adapter,
	}

	return service
}

func (s *Service) AcceptOrder(ctx context.Context, userID string) error {
	resp, err := s.repo.FindByFreeStatus(ctx)
	if err != nil {
		return fmt.Errorf("repo.FindByFreeStatus: %w", err)
	}
	resp.UserID = userID
	
	return s.adapter.AcceptOrder(ctx, resp)
}
