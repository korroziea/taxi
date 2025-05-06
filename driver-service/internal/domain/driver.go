package domain

import (
	"time"

	"github.com/korroziea/taxi/driver-service/pkg/utils"
)

type CarType string

const (
	Economy  CarType = "economy"
	Comfort  CarType = "comfort"
	Business CarType = "business"
)

type Car struct {
	ID        string
	Number    string
	Color     string
	Type      CarType
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenCarID() (string, error) {
	id, err := utils.GenID()
	if err != nil {
		return "", err
	}

	return "car_" + id, nil
}

type WorkStatus string

const (
	Free     WorkStatus = "free"
	Busy     WorkStatus = "busy"
	OffShift WorkStatus = "off-shift"
)

type Driver struct {
	ID        string
	FirstName string
	Phone     string
	Email     string
	Password  string
	Rate      int8
	Status    WorkStatus
	CarID     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GenDriverID() (string, error) {
	id, err := utils.GenID()
	if err != nil {
		return "", err
	}

	return "driver_" + id, nil
}

type SignUpDriver struct {
	ID        string
	FirstName string
	Phone     string
	Email     string
	Password  string
}

type SignInDriver struct {
	Phone    string
	Password string
}

type acceptOrderDriver struct {
	ID        string
	FirstName string
	Rate      int16
}

type acceptOrderCar struct {
	ID     string
	Number string
	Color  string
}

type AcceptOrderResp struct {
	UserID string
	Driver acceptOrderDriver
	Car    acceptOrderCar
}
