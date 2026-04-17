package domain

import "time"

type Transaction struct {
	ID          int
	Date        time.Time
	Amount      float64
	Description string
	Category    string
}
