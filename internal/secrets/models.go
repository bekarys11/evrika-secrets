package secrets

import (
	"time"
)

type Secret struct {
	ID        int       `json:"id" db:"id" jsonapi:"primary,secrets"`
	Title     string    `json:"title" jsonapi:"attr,title"`
	Key       string    `json:"key" jsonapi:"attr,key"`
	Data      string    `json:"data" jsonapi:"attr,data"`
	AuthorId  int       `json:"author_id" db:"author_id" jsonapi:"attr,author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
}

type UsersSecret struct {
	UserIds  []int `json:"user_ids" db:"user_id"`
	SecretId int   `json:"secret_id" db:"secret_id"`
}
