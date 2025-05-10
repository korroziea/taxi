package trip

import (
	"time"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type tripsReq struct {
	UserID string `json:"user_id"`
}

func toTripsReq(userID string) tripsReq {
	req := tripsReq{
		UserID: userID,
	}

	return req
}

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

func toTripResp(trip tripResp) domain.Trip {
	resp := domain.Trip{
		ID:           trip.ID,
		Status:       domain.TripStatus(trip.Status),
		UserID:       trip.UserID,
		Cost:         trip.Cost,
		Start:        domain.MapPoint(trip.Start),
		End:          domain.MapPoint(trip.End),
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

func toDomains(trips []tripResp) []domain.Trip {
	resp := make([]domain.Trip, len(trips))

	for i := range resp {
		resp[i] = toTripResp(trips[i])
	}

	return resp
}
