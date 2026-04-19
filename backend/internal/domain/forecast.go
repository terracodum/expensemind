package domain

import "time"

type Point struct {
	T       int
	Balance float64
}

type Forecast struct {
	Points           []Point
	PredictedBalance float64
	Confidence       float64
}

type ForecastJob struct {
	ID        int
	Status    string // pending | running | done | failed
	Result    *Forecast
	CreatedAt time.Time
}
