package sqlstore

import (
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

	err := o.store.db.QueryRow(
		"INSERT INTO orders (user_id, createdAt, totalamount) VALUES ($1, $2, $3) RETURNING id",
		order.UserId,
		order.CreatedAt,
		order.TotalAmount,
	).Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete ...
func (o *OrderRepository) Delete(id int) error {
	_, err := o.store.db.Exec("DELETE FROM orders WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

// GetOrders ...
func (o *OrderRepository) GetOrders(userId int) ([]*model.Order, error) {
	var orders []*model.Order

	rows, err := o.store.db.Query("SELECT id, user_id, createdat, totalamount FROM orders WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o model.Order
		if err := rows.Scan(&o.ID, &o.UserId, &o.CreatedAt, &o.TotalAmount); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
