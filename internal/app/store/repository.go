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

type CategoryRepository interface {
	Create(category *model.Category) error
	Find(id int) (*model.Category, error)
	GetAllCategories() ([]*model.Category, error)
}

type MenuItemRepository interface {
	Create(menuItem *model.MenuItem) error
	FindByCategoryId(categoryId int) ([]*model.MenuItem, error)
}
