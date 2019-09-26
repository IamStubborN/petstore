//go:generate easyjson -all pet.go

package models

// easyjson:json
type PetList []*Pet

// easyjson:json
type Pet struct {
	ID        int64    `json:"id,omitempty" db:"id" validate:"nonzero"`
	Category  Category `json:"category" db:"category"`
	Name      string   `json:"name" db:"name" validate:"nonzero"`
	PhotoURLs []string `json:"photo_urls,omitempty" db:"photo_urls,omitempty"`
	Tags      []Tag    `json:"tags,omitempty" db:"tags,omitempty"`
	Status    string   `json:"status" db:"pet_status_id" validate:"nonzero"`
}
