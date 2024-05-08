package store

import (
	"github.com/yeboka/final-project/internal/app/model"
)

// UserRepository ...
type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Update(user *model.User) error
	UpdateRole(id int, role string) error
	Delete(id int) error
}

// OrderRepository ...
type OrderRepository interface {
	Create(order *model.Order) error
	Delete(id int) error
	GetOrder(id int) (*model.Order, error)
}

type CategoryRepository interface {
	Create(category *model.Category) error
	Find(id int) (*model.Category, error)
	GetAllCategories() ([]*model.Category, error)
}

type MenuItemRepository interface {
	Create(menuItem *model.MenuItem) error
	GetPrice(id int) int
	//FindByName(id int) (*model.MenuItem, error)
	FindByCategoryId(categoryId int) ([]*model.MenuItem, error)
	Update(mi *model.MenuItem) error
	Delete(id int) error
}

type OrderItemRepository interface {
	Create(item *model.OrderItem) error
	Delete(id int) error
	DeleteAllOrder(orderId int) error
}
