//go:generate easyjson -all category.go

package models

// easyjson:json
type Category struct {
	ID    int64   `json:"id" db:"id"`
	Name  string  `json:"name" db:"name" validate:"nonzero"`
	Price float64 `json:"price,omitempty" db:"price" validate:"min=0"`
}
