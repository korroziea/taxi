package trip

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korroziea/taxi/driver-service/internal/domain"
)

const (
	drivers = "drivers"
	cars    = "cars"
)

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

func (r *Repo) FindByFreeStatus(ctx context.Context) (domain.AcceptOrderResp, error) {
	query, args, err := sq.
		Select(
			"drivers.id, drivers.first_name, drivers.rate, cars.id, cars.number, cars.color",
		).
		From(drivers).
		Join(
			cars + " ON " + drivers + ".car_id" + " = " + cars + ".id",
		).
		Where(
			sq.Eq{
				"drivers.status": "free",
			},
		).
		Limit(1).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.AcceptOrderResp{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

const queryTimeout = 5 * time.Second

func (r *Repo) doQueryRow(ctx context.Context, query string, args ...any) (domain.AcceptOrderResp, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var resp domain.AcceptOrderResp
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&resp.Driver.ID,
		&resp.Driver.FirstName,
		&resp.Driver.Rate,
		&resp.Car.ID,
		&resp.Car.Number,
		&resp.Car.Color,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.AcceptOrderResp{}, fmt.Errorf("driver or car: %w", domain.ErrDriverCarNotFound)
		}

		return domain.AcceptOrderResp{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return resp, nil
}
