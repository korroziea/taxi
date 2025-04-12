package wallet

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, walletID string) (domain.Wallet, error)
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

func (s *Service) CreateWallet(ctx context.Context) (domain.Wallet, error) {
	walletID, err := domain.GenWalletID()
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("GenWalletID: %w", err)
	}

	wallet, err := s.repo.Create(ctx, walletID)
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("repo.Create: %w", err)
	}

	return wallet, nil
}

func (s *Service) ChangeType(ctx context.Context, walletID string) (domain.Wallet, error) {
	return domain.Wallet{}, nil
}

func (s *Service) Refill(ctx context.Context, amount int64) (int64, error) {
	return 0, nil
}
