package sqlstore

import (
	"database/sql"
	"fmt"
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
	fmt.Println(c)

	if c.ParentID <= 0 {
		return r.store.db.QueryRow(
			"INSERT INTO categories (name) VALUES ($1) RETURNING id",
			c.Name,
		).Scan(&c.ID)
	}

	return r.store.db.QueryRow(
		"INSERT INTO categories (name, parent_id) VALUES ($1, $2) RETURNING id",
		c.Name,
		c.ParentID,
	).Scan(&c.ID)
}

func (r *CategoryRepository) Find(id int) (*model.Category, error) {
	c := &model.Category{}

	var parentID sql.NullInt64

	if err := r.store.db.QueryRow(
		"SELECT id, name, parent_id FROM categories WHERE id = $1",
		id,
	).Scan(
		&c.ID,
		&c.Name,
		&parentID,
	); err != nil {
		return nil, err
	}

	if parentID.Valid {
		c.ParentID = int(parentID.Int64)
	} else {
		c.ParentID = -1
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
		var parentID sql.NullInt64
		c := &model.Category{}
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
			c.ParentID = -1
		}

		fmt.Println(c)
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}
