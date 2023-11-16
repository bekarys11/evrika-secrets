package roles

import "time"

type Role struct {
	ID        int       `json:"id" db:"id" jsonapi:"primary,roles"`
	Name      string    `json:"name" jsonapi:"attr,name"`
	Alias     string    `json:"alias" jsonapi:"attr,alias"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
}
