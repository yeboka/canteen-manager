package sqlstore

import "github.com/yeboka/final-project/internal/app/model"

// MenuItemRepository ...
type MenuItemRepository struct {
	store *Store
}

func (r *MenuItemRepository) Create(m *model.MenuItem) error {
	return r.store.db.QueryRow(
		"INSERT INTO menuitem (name, category_id, price, description) VALUES ($1, $2, $3, $4) RETURNING id",
		m.Name,
		m.CategoryID,
		m.Price,
		m.Description,
	).Scan(&m.ID)
}

func (r *MenuItemRepository) FindByCategoryId(categoryId int) ([]*model.MenuItem, error) {
	rows, err := r.store.db.Query(
		"SELECT id, category_id, name, price, description FROM menuitem WHERE category_id = $1",
		categoryId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []*model.MenuItem
	for rows.Next() {
		menuItem := &model.MenuItem{}
		if err := rows.Scan(
			&menuItem.ID,
			&menuItem.CategoryID,
			&menuItem.Name,
			&menuItem.Price,
			&menuItem.Description,
		); err != nil {
			return nil, err
		}
		menuItems = append(menuItems, menuItem)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return menuItems, nil
}

//func (r *MenuItemRepository) Find(id int) (*model.MenuItem, error) {
//	m := &model.MenuItem{}
//
//	if err := r.store.db.QueryRow(
//		"SELECT id, name, category_id, price, description FROM menuitem WHERE id = $1",
//		id,
//	).Scan(
//		&m.ID,
//		&m.Name,
//		&m.CategoryID,
//		&m.Price,
//		&m.Description,
//	); err != nil {
//		return nil, err
//	}
//
//	return m, nil
//}
