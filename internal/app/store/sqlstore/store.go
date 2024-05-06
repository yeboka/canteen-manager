package sqlstore

import (
	"database/sql"
	"github.com/yeboka/final-project/internal/app/store"

	_ "github.com/lib/pq" // ...
)

// Store ...
type Store struct {
	db                 *sql.DB
	UserRepository     *UserRepository
	CategoryRepository *CategoryRepository
	MenuItemRepository *MenuItemRepository
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
