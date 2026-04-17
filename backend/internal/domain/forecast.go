package domain

type Point struct {
	T       int
	Balance float64
}

type Forecast struct {
	Points           []Point
	PredictedBalance float64
	Confidence       float64
}
