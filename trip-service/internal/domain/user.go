package domain

type MapPoint struct {
	Lon  float64
	Lan  float64
	Type string
}

type StartTrip struct {
	ID        string
	UserID    string
	Start     MapPoint
	End       MapPoint
	Cost      int64
	CreatedAt int64
}
