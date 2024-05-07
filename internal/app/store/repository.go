package store

import (
	"database/sql"
	"github.com/yeboka/final-project/internal/app/model"
)

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
}

// OrderRepository ...
type OrderRepository interface {
	Create(order *model.Order) error
	Delete(orderId int) (sql.Result, error)
}

type CategoryRepository interface {
	Create(category *model.Category) error
	Find(id int) (*model.Category, error)
}

type MenuItemRepository interface {
	Create(menuItem *model.MenuItem) error
	//FindByName(id int) (*model.MenuItem, error)
}
