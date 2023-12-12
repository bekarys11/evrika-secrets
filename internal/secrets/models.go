package secrets

import (
	"time"
)

type Secret struct {
	ID        int       `json:"id" db:"id" jsonapi:"primary,secrets" example:"11"`
	Title     string    `json:"title" jsonapi:"attr,title" validate:"required" example:"адрес базы данных"`
	Key       string    `json:"key" jsonapi:"attr,key" validate:"required" example:"dbHost"`
	Data      string    `json:"data" jsonapi:"attr,data" validate:"required" example:"localhost:5432"`
	Type      string    `json:"type" jsonapi:"attr,type" validate:"oneof= auth ssh env" example:"env"`
	AuthorId  int       `json:"author_id" db:"author_id" jsonapi:"attr,author_id" validate:"required"  example:"1"`
	CreatedAt time.Time `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601" example:"2023-11-20T11:15:37Z"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601" example:"2023-11-20T11:15:37Z"`
}

type SecretResp struct {
	ID        int         `json:"id" db:"id" jsonapi:"primary,secrets" example:"11"`
	Title     string      `json:"title" jsonapi:"attr,title" validate:"required" example:"адрес базы данных"`
	Key       string      `json:"key" jsonapi:"attr,key" validate:"required" example:"dbHost"`
	Data      string      `json:"data" jsonapi:"attr,data" validate:"required" example:"localhost:5432"`
	Type      string      `json:"type" jsonapi:"attr,type" validate:"oneof= auth ssh env" example:"env"`
	AuthorId  int         `json:"author_id" db:"author_id" jsonapi:"attr,author_id" validate:"required"  example:"1"`
	Users     interface{} `json:"users" jsonapi:"attr,users"`
	CreatedAt time.Time   `json:"created_at" db:"created_at" jsonapi:"attr,created_at,iso8601" example:"2023-11-20T11:15:37Z"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at" jsonapi:"attr,updated_at,iso8601" example:"2023-11-20T11:15:37Z"`
}

type SecretReq struct {
	Title    string `json:"title" jsonapi:"attr,title" validate:"required" example:"адрес базы данных"`
	Key      string `json:"key" jsonapi:"attr,key" validate:"required" example:"dbHost"`
	Data     string `json:"data" jsonapi:"attr,data" validate:"required" example:"localhost:5432"`
	Type     string `json:"type" jsonapi:"attr,type" validate:"oneof= auth ssh env" example:"env"`
	AuthorId int    `json:"author_id" db:"author_id" jsonapi:"attr,author_id" validate:"required"  example:"1"`
}

type UsersSecret struct {
	UserIds  []int `json:"user_ids" db:"user_id" example:"13"`
	SecretId int   `json:"secret_id" db:"secret_id" example:"41"`
}

type SecretSwaggerData struct {
	ID         string `json:"id" example:"11"`
	Type       string `json:"type" example:"secrets"`
	Attributes Secret `json:"attributes"`
}

type SecretSwaggerJson struct {
	Data []SecretSwaggerData `json:"data"`
}
type SecretSwaggerJsonObj struct {
	Data SecretSwaggerData `json:"data"`
}

type SecretSwaggerReq struct {
	Title    string `json:"title" jsonapi:"attr,title" validate:"required" example:"адрес базы данных"`
	Key      string `json:"key" jsonapi:"attr,key" validate:"required" example:"dbHost"`
	Data     string `json:"data" jsonapi:"attr,data" validate:"required" example:"localhost:5432"`
	Type     string `json:"type" example:"env"`
	AuthorId int    `json:"author_id" db:"author_id" jsonapi:"attr,author_id" validate:"required"  example:"1"`
}
