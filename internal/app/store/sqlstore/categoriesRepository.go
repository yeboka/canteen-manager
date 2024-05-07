package sqlstore

import (
	"database/sql"
	"github.com/yeboka/final-project/internal/app/model"
)

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

func (r *CategoryRepository) GetAllCategories() ([]*model.Category, error) {
	rows, err := r.store.db.Query(
		"SELECT id, name, parent_id FROM categories",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		c := &model.Category{}
		var parentID sql.NullInt64
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&parentID,
		); err != nil {
			return nil, err
		}
		if parentID.Valid {
			c.ParentID = int(parentID.Int64)
		} else {
			c.ParentID = 0
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
