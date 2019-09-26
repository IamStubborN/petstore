//go:generate easyjson -all user.go

package models

// easyjson:json
type UserList []*User

// easyjson:json
type User struct {
	ID         int64  `json:"id,omitempty" db:"id"`
	Email      string `json:"email" db:"email" validate:"regexp=^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+.[a-zA-Z0-9-.]+$"`
	Username   string `json:"user_name" db:"user_name" validate:"regexp=^[a-zA-Z0-9]{5\\,50}$"`
	Password   string `json:"password" db:"password" validate:"regexp=^[a-zA-Z0-9]{5\\,50}$"`
	FirstName  string `json:"first_name,omitempty" db:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty" db:"last_name,omitempty"`
	Phone      string `json:"phone" db:"phone" validate:"regexp=^\\+[0-9]{9\\,15}$"`
	UserStatus int64  `json:"user_status_id" db:"user_status_id" validate:"nonzero"`
}
