package users

import (
	"github.com/bekarys11/evrika-secrets/internal/roles"
	"time"
)

type User struct {
	ID        int       `json:"id" jsonapi:"primary,users"`
	Name      string    `json:"name" jsonapi:"attr,name"`
	Email     string    `json:"email" jsonapi:"attr,email"`
	Password  string    `json:"-"`
	IsActive  bool      `json:"is_active" db:"is_active" jsonapi:"attr,is_active"`
	RoleId    int       `json:"role_id" db:"role_id" jsonapi:"attr,role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
}

type UserResp struct {
	ID        int         `json:"id" jsonapi:"primary,users"`
	Name      string      `json:"name" validate:"required" jsonapi:"attr,name"`
	Email     string      `json:"email" validate:"required,email" jsonapi:"attr,email"`
	Password  string      `json:"-" validate:"required"`
	IsActive  bool        `json:"is_active" db:"is_active" jsonapi:"attr,is_active"`
	RoleId    int         `json:"role_id" db:"role_id" jsonapi:"attr,role_id"`
	CreatedAt time.Time   `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
	Role      *roles.Role `json:"role" jsonapi:"relation,role"`
}

type UserRequest struct {
	ID        int       `json:"id" jsonapi:"primary,users"`
	Name      string    `json:"name" validate:"required" jsonapi:"attr,name"`
	Email     string    `json:"email" validate:"required,email" jsonapi:"attr,email"`
	Password  string    `json:"password" validate:"required,gte=7,lte=20"`
	IsActive  bool      `json:"is_active" db:"is_active" jsonapi:"attr,is_active"`
	RoleId    int       `json:"role_id" db:"role_id" validate:"required,oneof=1 2" jsonapi:"attr,role_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601"`
}
