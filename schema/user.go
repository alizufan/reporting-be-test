package schema

import "time"

type User struct {
	ID        uint64    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Username  string    `json:"user_name" db:"user_name"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	CreatedBy uint64    `json:"created_by" db:"created_by"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	UpdatedBy uint64    `json:"updated_by" db:"updated_by"`
}

func (User) Table() string {
	return "Users"
}
