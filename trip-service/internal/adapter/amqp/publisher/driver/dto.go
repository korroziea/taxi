package driver

import (
	"github.com/korroziea/taxi/trip-service/internal/domain"
)

type mapPoint struct {
	Lon  float64 `json:"lon"`
	Lan  float64 `json:"lan"`
	Type string  `json:"type"`
}

type findDriverBody struct {
	UserID string   `json:"user_id"`
	Start  mapPoint `json:"start"`
	End    mapPoint `json:"end"`
}

func toFindDriverBody(req domain.FindDriverReq) findDriverBody {
	body := findDriverBody{
		UserID: req.UserID,
		Start:  mapPoint(req.Start),
		End:    mapPoint(req.End),
	}

	return body
}
