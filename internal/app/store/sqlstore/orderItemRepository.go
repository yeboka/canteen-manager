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
