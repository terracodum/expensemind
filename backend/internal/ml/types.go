package ml

type TimePoint struct {
	T                  int     `json:"t"`
	Balance            float64 `json:"balance"`
	DayOfWeek          int     `json:"day_of_week"`
	IsWeekend          bool    `json:"is_weekend"`
	FoodTotal          float64 `json:"food_total"`
	TransportTotal     float64 `json:"transport_total"`
	EntertainmentTotal float64 `json:"entertainment_total"`
	AvgTransactionSize float64 `json:"avg_transaction_size"`
	TransactionCount   int     `json:"transaction_count"`
}

type IncomeEvent struct {
	T      int     `json:"t"`
	Amount float64 `json:"amount"`
	Label  string  `json:"label"`
}

type Features struct {
	AvgDailyExpense float64       `json:"avg_daily_expense"`
	IncomeEvents    []IncomeEvent `json:"income_events"`
}

type PredictRequest struct {
	Timeseries []TimePoint `json:"timeseries"`
	Horizon    int         `json:"horizon"`
	Features   Features    `json:"features"`
}

type forecastPoint struct {
	T       int     `json:"t"`
	Balance float64 `json:"balance"`
}

type predictResponse struct {
	Forecast         []forecastPoint `json:"forecast"`
	PredictedBalance float64         `json:"predicted_balance"`
	Confidence       float64         `json:"confidence"`
}
