package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Category struct {
	ID        int        `json:"id"`
	ParentID  int        `json:"parent_id,omitempty"`
	Name      string     `json:"name"`
	MenuItems []MenuItem `json:"menu_items"`
}

// Validate ...
func (c *Category) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Name, validation.NilOrNotEmpty, validation.Length(1, 45)),
	)
}
