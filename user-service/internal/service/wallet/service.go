package wallet

import (
	"context"
	"fmt"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type Repo interface {
	Create(ctx context.Context, walletID string) (domain.ViewWallet, error)
	FindByUserID(ctx context.Context) ([]domain.ViewWallet, error)
	FindByUserAndWalletIDs(ctx context.Context, walletID string) (domain.ViewWallet, error)
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

func (s *Service) CreateWallet(ctx context.Context) (domain.ViewWallet, error) {
	walletID, err := domain.GenWalletID()
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("GenWalletID: %w", err)
	}

	wallet, err := s.repo.Create(ctx, walletID)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("repo.Create: %w", err)
	}

	return wallet, nil
}

func (s *Service) WalletList(ctx context.Context) ([]domain.ViewWallet, error) {
	wallets, err := s.repo.FindByUserID(ctx)
	if err != nil {
		return []domain.ViewWallet{}, fmt.Errorf("repo.FindByUserID: %w", err)
	}

	return wallets, nil
}

func (s *Service) Wallet(ctx context.Context, walletID string) (domain.ViewWallet, error) {
	wallet, err := s.repo.FindByUserAndWalletIDs(ctx, walletID)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("repo.FindByUserAndWalletIDs: %w", err)
	}

	return wallet, nil
}

func (s *Service) ChangeType(ctx context.Context, walletID string) (domain.ViewWallet, error) {
	return domain.ViewWallet{}, nil
}

func (s *Service) Refill(ctx context.Context, amount int64) (domain.ViewWallet, error) {
	return domain.ViewWallet{}, nil
}
