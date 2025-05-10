package user

import (
	"time"

	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type mapPoint struct {
	Lon  float64 `json:"lon"`
	Lan  float64 `json:"lan"`
	Type string  `json:"type"`
}

type tripResp struct {
	ID           string            `json:"id"`
	Status       domain.TripStatus `json:"status"`
	UserID       string            `json:"user_id"`
	Cost         int64             `json:"cost"`
	Start        mapPoint          `json:"start"`
	End          mapPoint          `json:"end"`
	Distance     int32             `json:"distance"`
	Duration     int32             `json:"duration"`
	DriverID     string            `json:"driver_id"`
	DriverName   string            `json:"driver_name"`
	DriverRating int16             `json:"driver_rating"`
	CarID        string            `json:"car_id"`
	CarNumber    string            `json:"car_number"`
	CarColor     string            `json:"car_color"`
	WaitingTime  int32             `json:"waiting_time"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func toTripResp(trip domain.Trip) tripResp {
	resp := tripResp{
		ID:           trip.ID,
		Status:       domain.TripStatus(trip.Status),
		UserID:       trip.UserID,
		Cost:         trip.Cost,
		Start:        mapPoint(trip.Start),
		End:          mapPoint(trip.End),
		Distance:     trip.Distance,
		Duration:     trip.Duration,
		DriverID:     trip.DriverID,
		DriverName:   trip.DriverName,
		DriverRating: trip.DriverRating,
		CarID:        trip.CarID,
		CarNumber:    trip.CarNumber,
		CarColor:     trip.CarColor,
		WaitingTime:  trip.WaitingTime,
		CreatedAt:    trip.CreatedAt,
		UpdatedAt:    trip.UpdatedAt,
	}

	return resp
}

func toTripsResp(trips []domain.Trip) []tripResp {
	resp := make([]tripResp, len(trips))

	for i := range resp {
		resp[i] = toTripResp(trips[i])
	}

	return resp
}
