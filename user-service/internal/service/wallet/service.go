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
	UpdateType(ctx context.Context, walletID string) (domain.Wallet, error)
	UpdateBalance(ctx context.Context, walletID string, amount int64) (domain.Wallet, error)
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
	return s.repo.FindByUserID(ctx)
}

func (s *Service) Wallet(ctx context.Context, walletID string) (domain.ViewWallet, error) {
	return s.repo.FindByUserAndWalletIDs(ctx, walletID)
}

func (s *Service) ChangeType(ctx context.Context, walletID string) (domain.ViewWallet, error) {
	wallet, err := s.repo.FindByUserAndWalletIDs(ctx, walletID)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("repo.FindByUserAndWalletIDs: %w", err)
	}

	if wallet.Type == domain.Family {
		return domain.ViewWallet{}, fmt.Errorf("wallet type already equals to family type: %w", domain.ErrChangeWalletType)
	}

	_, err = s.repo.UpdateType(ctx, walletID)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("repo.UpdateType: %w", err)
	}

	return s.repo.FindByUserAndWalletIDs(ctx, walletID)
}

func (s *Service) Refill(ctx context.Context, walletID string, amount int64) (domain.ViewWallet, error) {
	_, err := s.repo.UpdateBalance(ctx, walletID, amount)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("repo.UpdateBalance: %w", err)
	}

	return s.repo.FindByUserAndWalletIDs(ctx, walletID)
}
