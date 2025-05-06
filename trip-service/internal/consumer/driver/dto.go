package driver

import "github.com/korroziea/taxi/trip-service/internal/domain"

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

type acceptTripReq struct {
	UserID string            `json:"user_id"`
	Driver acceptOrderDriver `json:"driver"`
	Car    acceptOrderCar    `json:"car"`
}

func (r acceptTripReq) toDomain() domain.AcceptOrderReq {
	req := domain.AcceptOrderReq{
		UserID: r.UserID,
		Driver: domain.AcceptOrderDriver(r.Driver),
		Car:    domain.AcceptOrderCar(r.Car),
	}

	return req
}
