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

const wallets = "wallets"

var (
	walletColumns = []string{
		"id", "owner_id", "type", "balance", "created_at", "updated_at",
	}

	userWalletsColumn = []string{
		"user_id", "wallet_id", "role",
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

func (r *Repo) Create(ctx context.Context, walletID string) (domain.Wallet, error) {
	query, args, err := sq.
		Insert(wallets).
		Columns(
			"id",
			"owner_id",
			"created_at",
			"updated_at",
		).
		Values(
			walletID,
			wallethndl.FromContext(ctx),
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING id, owner_id, type, balance, created_at, updated_at",
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
		&wallet.OwnerID,
		&wallet.Type,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wallet{}, fmt.Errorf("wallet: %w", domain.ErrUserNotFound)
		}

		return domain.Wallet{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return wallet, nil
}
