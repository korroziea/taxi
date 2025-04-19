package domain

import "time"

type Trip struct {
	Distance     int32
	Duration     int32
	DriverName   string
	DriverRating int8
	CarNumber    string
	CarColor     string
	WaitingTime  int32
	CreatedAt    time.Time
}

type MapPoint struct {
	Lon  float64
	Lan  float64
	Type string
}

type StartTrip struct {
	Start     MapPoint
	End       MapPoint
	Cost      int64
	CreatedAt int64
}
