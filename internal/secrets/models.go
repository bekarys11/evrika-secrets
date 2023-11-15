package secrets

import (
	"time"
)

type Secret struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title"`
	Key       string    `json:"key"`
	Data      string    `json:"data"`
	AuthorId  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
