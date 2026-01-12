package dto

import "session-23/internal/data/entity"

type PriceStats struct {
	Min float64
	Max float64
	Avg float64
}

type DashboardResponse struct {
	TotalCars int64
	Stats     PriceStats
	Cars      []entity.Car
}

type ResultCars struct {
	Data []entity.Car
	Err  error
}
type ResultTotal struct {
	Data int64
	Err  error
}
type ResultStats struct {
	Data PriceStats
	Err  error
}
