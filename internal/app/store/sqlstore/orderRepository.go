package sqlstore

import (
	"database/sql"
	"github.com/yeboka/final-project/internal/app/model"
	"time"
)

// OrderRepository ...
type OrderRepository struct {
	store *Store
}

// Create ...
func (o *OrderRepository) Create(order *model.Order) error {
	order.CreatedAt = time.Now()

	return o.store.db.QueryRow(
		"INSERT INTO orders (user_id, createdAt, totalAmount) VALUES ($1, $2, $3) RETURNING id",
		order.UserId,
		order.CreatedAt,
		order.TotalAmount,
	).Scan(&order.ID)
}

// Delete ...
func (o *OrderRepository) Delete(orderId int) (sql.Result, error) {
	res, err := o.store.db.Exec("DELETE FROM orders WHERE id = $1", orderId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
