package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korroziea/taxi/user-service/internal/domain"
)

const users = "users"

var (
	userColumns = []string{
		"id", "first_name", "email", "phone", "password", "created_at", "updated_at",
	}

	userColumnsWithoutPassword = []string{
		"id", "first_name", "email", "phone", "created_at", "updated_at",
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

func (r *Repo) Create(ctx context.Context, user domain.SignUpUser) (domain.User, error) {
	query, args, err := sq.
		Insert(users).
		Columns(
			"id",
			"first_name",
			"email",
			"phone",
			"password",
			"created_at",
			"updated_at",
		).
		Values(
			user.ID,
			user.FirstName,
			user.Email,
			user.Phone,
			user.Password,
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING id, first_name, email, phone, password, created_at, updated_at",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) UpdateProfile(ctx context.Context, user domain.ProfileUser) (domain.User, error) {
	query, args, err := sq.
		Update(users).
		Set("email", user.Email).
		Where(
			sq.Eq{
				"id": user.ID,
			},
		).
		Suffix(
			"RETURNING *",
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) FindByID(ctx context.Context, id string) (domain.User, error) {
	query, args, err := sq.
		Select(userColumns...).
		From(users).
		Where(
			sq.Eq{
				"id": id,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	query, args, err := sq.
		Select(userColumns...).
		From(users).
		Where(
			sq.Eq{
				"phone": phone,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) FindByPhoneAndPassword(ctx context.Context, user domain.SignInUser) (domain.User, error) {
	query, args, err := sq.
		Select(userColumnsWithoutPassword...).
		From(users).
		Where(
			sq.Eq{
				"phone":    user.Phone,
				"password": user.Password,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

const queryTimeout = 5 * time.Second

func (r *Repo) doQueryRow(ctx context.Context, query string, args ...any) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var user domain.User
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("user: %w", domain.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return user, nil
}
