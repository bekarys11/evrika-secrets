package roles

import "time"

type Role struct {
	ID        int       `json:"id" db:"id" jsonapi:"primary,roles" example:"12"`
	Name      string    `json:"name" jsonapi:"attr,name" example:"Пользователь"`
	Alias     string    `json:"alias" jsonapi:"attr,alias" example:"user"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601" example:"2023-11-20T11:15:37Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601" example:"2023-11-20T11:15:37Z"`
}

type RoleSwaggerData struct {
	ID         string `json:"id" example:"11"`
	Type       string `json:"type" example:"users"`
	Attributes Role   `json:"attributes"`
}

type RoleSwaggerJson struct {
	Data []RoleSwaggerData `json:"data"`
}
