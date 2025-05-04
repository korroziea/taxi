package trip

import "github.com/korroziea/taxi/user-service/internal/domain"

type mapPoint struct {
	Lon  float64 `json:"lon"`
	Lan  float64 `json:"lan"`
	Type string  `json:"type"`
}

type startTripBody struct {
	UserID    string   `json:"user_id"`
	Start     mapPoint `json:"start"`
	End       mapPoint `json:"end"`
	Cost      int64    `json:"cost"`
	CreatedAt int64    `json:"created_at"`
}

func toStartTripBody(trip domain.StartTrip) startTripBody {
	body := startTripBody{
		UserID:    trip.UserID,
		Start:     mapPoint(trip.Start),
		End:       mapPoint(trip.End),
		Cost:      trip.Cost,
		CreatedAt: trip.CreatedAt,
	}

	return body
}
