//go:generate easyjson -all tag.go

package models

// easyjson:json
type Tag struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"nonzero"`
}
