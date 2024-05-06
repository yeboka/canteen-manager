package store

import (
	"github.com/yeboka/final-project/internal/app/model"
)

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}
