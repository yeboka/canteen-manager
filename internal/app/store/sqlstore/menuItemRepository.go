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

func (r *MenuItemRepository) Update(mi *model.MenuItem) error {
	_, err := r.store.db.Query(
		"UPDATE menuitem SET name = $1, price = $2, description = $3 WHERE id = $4",
		mi.Name, mi.Price, mi.Description, mi.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *MenuItemRepository) Delete(id int) error {
	_, err := r.store.db.Exec("DELETE FROM menuitem WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
