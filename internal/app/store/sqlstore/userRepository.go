package sqlstore

import (
	"errors"
	"github.com/yeboka/final-project/internal/app/model"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (email, encrypted_password, username, role) VALUES ($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.EncryptedPassword,
		u.Username,
		u.Role,
	).Scan(&u.ID)
}

// Find ...
func (r *UserRepository) Find(id int) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, username, email, role FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Role,
	); err != nil {
		return nil, err
	}

	return u, nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}

	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users WHERE email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		return nil, err
	}

	return u, nil
}

// Update ...
func (r *UserRepository) Update(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	_, err := r.store.db.Exec(
		"UPDATE users SET username = $1, email = $2 WHERE id = $3",
		user.Username, user.Email, user.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRole ...
func (r *UserRepository) UpdateRole(userID int, newRole string) error {
	if newRole != "user" && newRole != "admin" {
		return errors.New("invalid role")
	}

	_, err := r.store.db.Exec(
		"UPDATE users SET role = $1 WHERE id = $2",
		newRole, userID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete ...
func (r *UserRepository) Delete(id int) error {
	_, err := r.store.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
