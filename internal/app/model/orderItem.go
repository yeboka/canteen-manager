package model

// OrderItem ...
type OrderItem struct {
	ID         int `json:"id"`
	OrderId    int `json:"order_id"`
	MenuItemId int `json:"menu_item_id"`
	Quantity   int `json:"quantity"`
}
