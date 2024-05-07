package store

// Store ...
type Store interface {
	User() UserRepository
	Order() OrderRepository
	Category() CategoryRepository
	MenuItem() MenuItemRepository
	OrderItem() OrderItemRepository
}
