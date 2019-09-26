//go:generate easyjson -all invoice_item.go

package models

// easyjson:json
type InvoiceItem struct {
	ID       int64   `json:"id," db:"id"`
	User     string  `json:"user_name" db:"user_name"`
	Pet      string  `json:"pet" db:"pet"`
	Category string  `json:"category" db:"category"`
	ShipDate string  `json:"ship_date" db:"ship_date"`
	Quantity int32   `json:"quantity" db:"quantity"`
	Price    float64 `json:"price" db:"price"`
}
