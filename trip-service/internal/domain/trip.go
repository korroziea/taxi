package domain

import (
	"time"

	"github.com/korroziea/taxi/trip-service/pkg/utils"
)

func GenTripID() (string, error) {
	id, err := utils.GenID()
	if err != nil {
		return "", err
	}

	return "trip_" + id, nil
}

type Trip struct {
	ID           string
	UserID       string
	Cost         int64
	Start        MapPoint
	End          MapPoint
	Distance     int32
	Duration     int32
	DriverID     string
	DriverName   string
	DriverRating int16
	CarID        string
	CarNumber    string
	CarColor     string
	WaitingTime  int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
