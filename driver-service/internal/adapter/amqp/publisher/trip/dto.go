package trip

import "github.com/korroziea/taxi/driver-service/internal/domain"

type acceptOrderDriver struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	Rate      int16  `json:"rate"`
}

type acceptOrderCar struct {
	ID     string `json:"id"`
	Number string `json:"number"`
	Color  string `json:"color"`
}

type acceptTripResp struct {
	UserID string            `json:"user_id"`
	Driver acceptOrderDriver `json:"driver"`
	Car    acceptOrderCar    `json:"car"`
}

func toAcceptTripResp(r domain.AcceptOrderResp) acceptTripResp {
	resp := acceptTripResp{
		UserID: r.UserID,
		Driver: acceptOrderDriver{
			ID:        r.Driver.ID,
			FirstName: r.Driver.FirstName,
			Rate:      r.Driver.Rate,
		},
		Car: acceptOrderCar{
			ID:     r.Car.ID,
			Number: r.Car.Number,
			Color:  r.Car.Color,
		},
	}

	return resp
}
