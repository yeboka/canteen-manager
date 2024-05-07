package model

import "time"

// Order ...
type Order struct {
	ID          int       `json:"id"`
	UserId      int       `json:"user_id"`
	CreatedAt   time.Time `json:"-"`
	TotalAmount int       `json:"totalAmount"`
}

// NewOrder ...
func NewOrder(id int, userId int, createdAt time.Time, totalAmount int) *Order {
	return &Order{
		ID:          id,
		UserId:      userId,
		CreatedAt:   createdAt,
		TotalAmount: totalAmount,
	}
}
