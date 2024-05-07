package model

import "time"

// Order ...
type Order struct {
	ID          int       `json:"id"`
	UserId      int       `json:"user_id"`
	CreatedAt   time.Time `json:"-"`
	TotalAmount int       `json:"-"`
}
