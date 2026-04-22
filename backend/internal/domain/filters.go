package domain

import "time"

type Filters struct {
	From     time.Time
	To       time.Time
	Category string
}
