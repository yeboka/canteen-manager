package sqlstore

import "github.com/yeboka/final-project/internal/app/model"

// CategoryRepository ...
type CategoryRepository struct {
	store *Store
}

func (r *CategoryRepository) Create(c *model.Category) error {

	if err := c.Validate(); err != nil {
		return err
	}

	var parentIDValue interface{}
	if c.ParentID != 0 {
		parentIDValue = c.ParentID
	} else {
		parentIDValue = nil
	}

	return r.store.db.QueryRow(
		"INSERT INTO categories (name, parent_id) VALUES ($1, $2) RETURNING id",
		c.Name,
		parentIDValue,
	).Scan(&c.ID)
}

func (r *CategoryRepository) Find(id int) (*model.Category, error) {
	c := &model.Category{}

	if err := r.store.db.QueryRow(
		"SELECT id, name, parent_id FROM categories WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.Name,
		&c.ParentID,
	); err != nil {
		return nil, err
	}

	return c, nil
}
