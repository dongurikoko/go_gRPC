package model

import "time"

type Todo struct {
	ID         int
	Title      string
	Created_at time.Time
	Updated_at time.Time
}
