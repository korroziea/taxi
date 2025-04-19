package trip

import (
	"time"

	"github.com/korroziea/taxi/user-service/internal/domain"
)

type mapPoint struct {
	Lon  float64 `json:"lon"`
	Lan  float64 `json:"lan"`
	Type string  `json:"type"`
}

func (p mapPoint) toDomain() domain.MapPoint {
	point := domain.MapPoint{
		Lon:  p.Lon,
		Lan:  p.Lan,
		Type: p.Type,
	}

	return point
}

type startTripReq struct {
	Start mapPoint `json:"start"`
	End   mapPoint `json:"end"`
}

func (r startTripReq) toDomain() domain.StartTrip {
	trip := domain.StartTrip{
		Start: r.Start.toDomain(),
		End:   r.End.toDomain(),
	}

	return trip
}

type tripResp struct {
	Distance     int32     `json:"distance"`
	Duration     int32     `json:"duration"`
	DriverName   string    `json:"driver_name"`
	DriverRating int8      `json:"driver_rating"`
	CarNumber    string    `json:"car_number"`
	CarColor     string    `json:"car_color"`
	WaitingTime  int32     `json:"waiting_time"`
	CreatedAt    time.Time `json:"created_at"`
}

func toView(trip domain.Trip) tripResp {
	resp := tripResp{
		Distance:     trip.Distance,
		Duration:     trip.Duration,
		DriverName:   trip.DriverName,
		DriverRating: trip.DriverRating,
		CarNumber:    trip.CarNumber,
		CarColor:     trip.CarColor,
		WaitingTime:  trip.WaitingTime,
		CreatedAt:    trip.CreatedAt,
	}

	return resp
}
