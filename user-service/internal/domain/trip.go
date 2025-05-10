package domain

import "time"

type TripStatus string

const (
	Processing TripStatus = "processing"
	Waiting    TripStatus = "waiting"
	Executing  TripStatus = "executing"
	Finished   TripStatus = "finished"
	Canceled   TripStatus = "canceled"
)

type Trip struct {
	ID           string
	Status       TripStatus
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

type MapPoint struct {
	Lon  float64
	Lan  float64
	Type string
}

type StartTrip struct {
	UserID    string
	Start     MapPoint
	End       MapPoint
	Cost      int64
	CreatedAt int64
}
