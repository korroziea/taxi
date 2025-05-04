package user

import "github.com/korroziea/taxi/trip-service/internal/domain"

type mapPoint struct {
	Lon  float64 `json:"lon"`
	Lan  float64 `json:"lan"`
	Type string  `json:"type"`
}

type startTrip struct {
	UserID    string   `json:"user_id"`
	Start     mapPoint `json:"start"`
	End       mapPoint `json:"end"`
	Cost      int64    `json:"cost"`
	CreatedAt int64    `json:"created_at"`
}

func (t startTrip) toDomain() domain.StartTrip {
	trip := domain.StartTrip{
		UserID:    t.UserID,
		Start:     domain.MapPoint(t.Start),
		End:       domain.MapPoint(t.End),
		Cost:      t.Cost,
		CreatedAt: t.CreatedAt,
	}

	return trip
}
