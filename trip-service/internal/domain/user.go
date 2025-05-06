package domain

type StartTrip struct {
	ID        string
	UserID    string
	Start     MapPoint
	End       MapPoint
	Cost      int64
	CreatedAt int64
}
