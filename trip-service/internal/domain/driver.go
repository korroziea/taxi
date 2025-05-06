package domain

type FindDriverReq struct {
	UserID string
	Start  MapPoint
	End    MapPoint
}

type AcceptOrderDriver struct {
	ID        string
	FirstName string
	Rate      int16
}

type AcceptOrderCar struct {
	ID     string
	Number string
	Color  string
}

type AcceptOrderReq struct {
	UserID string
	Driver AcceptOrderDriver
	Car    AcceptOrderCar
}
