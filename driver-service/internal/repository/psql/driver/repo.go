package driver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korroziea/taxi/driver-service/internal/domain"
)

const drivers = "drivers"

var (
	driverColumns = []string{
		"id", "first_name", "email", "phone", "password", "rate", "status", "car_id", "created_at", "updated_at",
	}

	driverColumnsWithoutPassword = []string{
		"id", "first_name", "email", "phone", "rate", "status", "car_id", "created_at", "updated_at",
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

func (r *Repo) Create(ctx context.Context, driver domain.SignUpDriver) (domain.Driver, error) {
	query, args, err := sq.
		Insert(drivers).
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
			driver.ID,
			driver.FirstName,
			driver.Email,
			driver.Phone,
			driver.Password,
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING id, first_name, email, phone, password, rate, status, car_id, created_at, updated_at",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Driver{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) UpdateStatus(ctx context.Context, driverID string, status domain.WorkStatus) (domain.Driver, error) {
	query, args, err := sq.
		Update(drivers).
		Set("status", status).
		Where(
			sq.Eq{
				"id": driverID,
			},
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING *",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Driver{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) FindByID(ctx context.Context, driverID string) (domain.Driver, error) {
	query, args, err := sq.
		Select(driverColumns...).
		From(drivers).
		Where(
			sq.Eq{
				"id": driverID,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Driver{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) FindByPhone(ctx context.Context, phone string) (domain.Driver, error) {
	query, args, err := sq.
		Select(driverColumns...).
		From(drivers).
		Where(
			sq.Eq{
				"phone": phone,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Driver{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) FindByPhoneAndPassword(ctx context.Context, user domain.SignInDriver) (domain.Driver, error) {
	query, args, err := sq.
		Select(driverColumnsWithoutPassword...).
		From(drivers).
		Where(
			sq.Eq{
				"phone":    user.Phone,
				"password": user.Password,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Driver{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

const queryTimeout = 5 * time.Second

func (r *Repo) doQueryRow(ctx context.Context, query string, args ...any) (domain.Driver, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var (
		driver domain.Driver
		carID  sql.NullString
	)
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&driver.ID,
		&driver.FirstName,
		&driver.Email,
		&driver.Phone,
		&driver.Password,
		&driver.Rate,
		&driver.Status,
		&carID,
		&driver.CreatedAt,
		&driver.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Driver{}, fmt.Errorf("user: %w", domain.ErrDriverNotFound)
		}

		return domain.Driver{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	if carID.Valid {
		driver.CarID = &carID.String
	}

	return driver, nil
}
