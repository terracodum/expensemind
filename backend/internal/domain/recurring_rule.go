package domain

import "time"

type RecurringRule struct {
	ID        int
	SourceID  string
	Type      string // income | expense
	Amount    float64
	Day       int
	StartDate time.Time
	Label     string
}
