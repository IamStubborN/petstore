//go:generate easyjson -all order.go

package models

// easyjson:json
type Order struct {
	ID       int64  `json:"id,omitempty," db:"id" validate:"nonzero"`
	PetID    int64  `json:"pet_id" db:"pet_id" validate:"nonzero"`
	UserID   int64  `json:"user_id" db:"user_id" validate:"nonzero"`
	Quantity int32  `json:"quantity" db:"quantity" validate:"nonzero"`
	ShipDate string `json:"ship_date" db:"ship_date"`
	Status   string `json:"status" db:"order_status" validate:"nonzero"`
	Complete bool   `json:"complete" db:"complete"`
}
