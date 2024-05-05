package teststore

import (
	"github.com/yeboka/final-project/internal/app/model"
	"github.com/yeboka/final-project/internal/app/store"
)

// Store ...
type Store struct {
	UserRepository *UserRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
		users: make(map[string]*model.User),
	}

	return s.UserRepository
}
