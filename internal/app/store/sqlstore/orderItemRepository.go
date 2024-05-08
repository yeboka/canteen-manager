package sqlstore

import "github.com/yeboka/final-project/internal/app/model"

// OrderItemRepository ...
type OrderItemRepository struct {
	s *Store
}

// Create ...
func (i *OrderItemRepository) Create(item *model.OrderItem) error {
	return i.s.db.QueryRow(
		"INSERT INTO orderitem (order_id, menu_item_id, quantity) VALUES ($1, $2, $3) RETURNING id",
		item.OrderId,
		item.MenuItemId,
		item.Quantity).Scan(&item.ID)
}

func (i *OrderItemRepository) Delete(id int) error {
	_, err := i.s.db.Exec("DELETE FROM orderitem WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}

func (i *OrderItemRepository) Update(menuItemId int, quantity int) error {
	_, err := i.s.db.Exec("UPDATE orderitem SET quantity = $1 WHERE menu_item_id = $2", quantity, menuItemId)
	if err != nil {
		return err
	}

	return nil
}

func (i *OrderItemRepository) DeleteAllOrder(orderId int) error {
	_, err := i.s.db.Exec("DELETE FROM orderitem WHERE order_id = $1", orderId)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderItems ...
func (i *OrderItemRepository) GetOrderItems(orderId int) ([]*model.OrderItem, error) {
	var orderItems []*model.OrderItem

	rows, err := i.s.db.Query("SELECT id, order_id, menu_item_id, quantity FROM orderitem WHERE order_id = $1", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var oi model.OrderItem
		if err := rows.Scan(&oi.ID, &oi.MenuItemId, &oi.OrderId, &oi.Quantity); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, &oi)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderItems, nil
}
