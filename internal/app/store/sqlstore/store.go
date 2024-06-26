package sqlstore

import (
	"database/sql"
	"github.com/yeboka/final-project/internal/app/store"

	_ "github.com/lib/pq" // ...
)

// Store ...
type Store struct {
	db                  *sql.DB
	UserRepository      *UserRepository
	CategoryRepository  *CategoryRepository
	MenuItemRepository  *MenuItemRepository
	OrderRepository     *OrderRepository
	OrderItemRepository *OrderItemRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{store: s}

	return s.UserRepository
}

// Order ...
func (s *Store) Order() store.OrderRepository {
	if s.OrderRepository != nil {
		return s.OrderRepository
	}

	s.OrderRepository = &OrderRepository{store: s}

	return s.OrderRepository
}

func (s *Store) Category() store.CategoryRepository {
	if s.CategoryRepository != nil {
		return s.CategoryRepository
	}

	s.CategoryRepository = &CategoryRepository{store: s}

	return s.CategoryRepository
}

func (s *Store) MenuItem() store.MenuItemRepository {
	if s.MenuItemRepository != nil {
		return s.MenuItemRepository
	}

	s.MenuItemRepository = &MenuItemRepository{store: s}

	return s.MenuItemRepository
}

func (s *Store) OrderItem() store.OrderItemRepository {
	if s.OrderItemRepository != nil {
		return s.OrderItemRepository
	}

	s.OrderItemRepository = &OrderItemRepository{s: s}

	return s.OrderItemRepository
}
