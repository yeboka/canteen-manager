package store

// Store ...
type Store interface {
	User() UserRepository
	Category() CategoryRepository
	MenuItem() MenuItemRepository
}
