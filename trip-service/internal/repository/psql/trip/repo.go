package trip

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/korroziea/taxi/trip-service/internal/domain"
)

const trips = "trips"

var (
	tripsColumns = []string{
		"id", "user_id", "cost", "start_point", "end_point", "distance", "duration", "driver_id", "driver_name", "driver_rate", "car_id", "car_number", "car_color", "waiting_time", "created_at", "updated_at",
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

func (r *Repo) Create(ctx context.Context, trip domain.StartTrip) (domain.Trip, error) {
	query, args, err := sq.
		Insert(trips).
		Columns(
			"id",
			"user_id",
			"cost",
			"start_point",
			"end_point",
			"created_at",
			"updated_at",
		).
		Values(
			trip.ID,
			trip.UserID,
			trip.Cost,
			trip.Start,
			trip.End,
			time.Now(),
			time.Now(),
		).
		Suffix(
			fmt.Sprintf(
				"RETURNING *",
			),
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Trip{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

func (r *Repo) UpdateDriverInfo(ctx context.Context, req domain.AcceptOrderReq) (domain.Trip, error) {
	query, args, err := sq.
		Update(trips).
		Set("status", "waiting").
		Set("driver_id", req.Driver.ID).
		Set("driver_name", req.Driver.FirstName).
		Set("driver_rate", req.Driver.Rate).
		Set("car_id", req.Car.ID).
		Set("car_number", req.Car.Number).
		Set("car_color", req.Car.Color).
		Where(
			sq.Eq{
				"user_id": req.UserID,
			},
		).
		Suffix(
			"RETURNING *",
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return domain.Trip{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	return r.doQueryRow(ctx, query, args...)
}

const queryTimeout = 5 * time.Second

func (r *Repo) doQueryRow(ctx context.Context, query string, args ...any) (domain.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var (
		trip domain.Trip

		distance     sql.NullInt32
		duration     sql.NullInt32
		driverID     sql.NullString
		driverName   sql.NullString
		driverRating sql.NullInt16
		carID        sql.NullString
		carNumber    sql.NullString
		carColor     sql.NullString
		waitingTime  sql.NullInt32
	)
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&trip.ID,
		&trip.Status,
		&trip.UserID,
		&trip.Cost,
		&trip.Start,
		&trip.End,
		&distance,
		&duration,
		&driverID,
		&driverName,
		&driverRating,
		&carID,
		&carNumber,
		&carColor,
		&waitingTime,
		&trip.CreatedAt,
		&trip.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Trip{}, fmt.Errorf("trip: %w", domain.ErrTripNotFound)
		}

		return domain.Trip{}, fmt.Errorf("%w: %w", domain.ErrInternal, err)
	}

	if distance.Valid {
		trip.Distance = distance.Int32
	}

	if duration.Valid {
		trip.Duration = duration.Int32
	}

	if driverID.Valid {
		trip.DriverID = driverID.String
	}

	if driverName.Valid {
		trip.DriverName = driverName.String
	}

	if driverRating.Valid {
		trip.DriverRating = driverRating.Int16
	}

	if carID.Valid {
		trip.CarID = carID.String
	}

	if carNumber.Valid {
		trip.CarNumber = carNumber.String
	}

	if carColor.Valid {
		trip.CarColor = carColor.String
	}

	if waitingTime.Valid {
		trip.WaitingTime = waitingTime.Int32
	}

	return trip, nil
}
