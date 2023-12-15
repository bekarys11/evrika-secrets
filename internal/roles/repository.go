package roles

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) GetRoles() ([]*Role, error) {
	var roles []*Role

	rows, err := repo.DB.Queryx("SELECT * FROM roles LIMIT 10")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var role Role
		if err := rows.StructScan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}
