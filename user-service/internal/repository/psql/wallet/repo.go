package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korroziea/taxi/user-service/internal/domain"
	wallethndl "github.com/korroziea/taxi/user-service/internal/handler/wallet"
)

const (
	wallets     = "wallets"
	userWallets = "user_wallets"
)

var (
	walletColumns = []string{
		"id", "type", "balance", "created_at", "updated_at",
	}

	userWalletsColumn = []string{
		"user_id", "wallet_id", "role", "created_at", "updated_at",
	}
)

type Repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repo {
	repo := &Repo{
		db: db,
	}

	return repo
}

func (r *Repo) Create(ctx context.Context, walletID string) (domain.ViewWallet, error) {
	wallet, err := r.createWallet(ctx, walletID)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("createWallet: %w", err)
	}

	userWallet, err := r.createUserWallet(ctx, walletID, domain.Owner)
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("createUserWallet: %w", err)
	}

	viewWallet := domain.ViewWallet{
		ID:      wallet.ID,
		Type:    wallet.Type,
		Role:    userWallet.Role,
		Balance: wallet.Balance,
	}

	return viewWallet, nil
}

func (r *Repo) createWallet(ctx context.Context, walletID string) (domain.Wallet, error) {
	query, args, err := sq.
		Insert(wallets).
		Columns(
			"id",
			"created_at",
			"updated_at",
		).
		Values(
			walletID,
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING id, type, balance, created_at, updated_at",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) createUserWallet(ctx context.Context, walletID string, role domain.UserWalletRole) (domain.UserWallet, error) {
	query, args, err := sq.
		Insert(userWallets).
		Columns(
			"user_id",
			"wallet_id",
			"role",
			"created_at",
			"updated_at",
		).
		Values(
			wallethndl.FromContext(ctx),
			walletID,
			role,
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING user_id, wallet_id, role, created_at, updated_at",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.UserWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doUserWalletQueryRow(ctx, query, args...)
}

func (r *Repo) FindByUserAndWalletIDs(ctx context.Context, walletID string) (domain.ViewWallet, error) {
	query, args, err := sq.
		Select(
			"wallets.id, wallets.type, user_wallets.role, wallets.balance",
		).
		From(userWallets).
		Join(
			wallets + " ON " + userWallets + ".wallet_id" + " = " + wallets + ".id",
		).
		Where(
			sq.Eq{
				"user_id":   wallethndl.FromContext(ctx),
				"wallet_id": walletID,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.ViewWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doViewWalletQueryRow(ctx, query, args...)
}

func (r *Repo) FindByUserID(ctx context.Context) ([]domain.ViewWallet, error) {
	query, args, err := sq.
		Select(
			"wallets.id, wallets.type, user_wallets.role, wallets.balance",
		).
		From(userWallets).
		Join(
			wallets + " ON " + userWallets + ".wallet_id" + " = " + wallets + ".id",
		).
		Where(
			sq.Eq{
				"user_id": wallethndl.FromContext(ctx),
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return []domain.ViewWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doViewWalletQueryRows(ctx, query, args...)
}

func (r *Repo) UpdateType(ctx context.Context, walletID string) (domain.Wallet, error) {
	query, args, err := sq.
		Update(wallets).
		Set("type", domain.Family).
		Where(
			sq.Eq{
				"id": walletID,
			},
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING id, type, balance, created_at, updated_at",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) UpdateBalance(ctx context.Context, walletID string, amount int64) (domain.Wallet, error) {
	query, args, err := sq.
		Update(wallets).
		Set("balance", amount). // todo: sum with previous
		Where(
			sq.Eq{
				"id": walletID,
			},
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING id, type, balance, created_at, updated_at",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Wallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

const queryTimeout = 5 * time.Second

func (r *Repo) doQueryRow(ctx context.Context, query string, args ...any) (domain.Wallet, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var wallet domain.Wallet
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&wallet.ID,
		&wallet.Type,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wallet{}, fmt.Errorf("wallet: %w", domain.ErrWalletNotFound)
		}

		return domain.Wallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return wallet, nil
}

func (r *Repo) doUserWalletQueryRow(ctx context.Context, query string, args ...any) (domain.UserWallet, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var userWallet domain.UserWallet
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&userWallet.UserID,
		&userWallet.WalletID,
		&userWallet.Role,
		&userWallet.CreatedAt,
		&userWallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserWallet{}, fmt.Errorf("wallet: %w", domain.ErrWalletNotFound)
		}

		return domain.UserWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return userWallet, nil
}

func (r *Repo) doViewWalletQueryRow(ctx context.Context, query string, args ...any) (domain.ViewWallet, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var viewWallet domain.ViewWallet
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&viewWallet.ID,
		&viewWallet.Type,
		&viewWallet.Role,
		&viewWallet.Balance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ViewWallet{}, fmt.Errorf("wallet: %w", domain.ErrWalletNotFound)
		}

		return domain.ViewWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return viewWallet, nil
}

func (r *Repo) doViewWalletQueryRows(ctx context.Context, query string, args ...any) ([]domain.ViewWallet, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	// .Scan(
	// 	&viewWallet.ID,
	// 	&viewWallet.Type,
	// 	&viewWallet.Role,
	// 	&viewWallet.Balance,
	// )

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain.ViewWallet{}, fmt.Errorf("wallet: %w", domain.ErrWalletNotFound)
		}

		return []domain.ViewWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}
	defer rows.Close()

	var viewWallets []domain.ViewWallet
	for rows.Next() {
		var vw domain.ViewWallet
		err := rows.Scan(
			&vw.ID,
			&vw.Type,
			&vw.Role,
			&vw.Balance,
		)
		if err != nil {
			return []domain.ViewWallet{}, fmt.Errorf("rows.Scan: %w", err)
		}

		viewWallets = append(viewWallets, vw)
	}

	if rows.Err() != nil {
		return []domain.ViewWallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, rows.Err())
	}

	return viewWallets, nil
}
