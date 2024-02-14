package roles

import (
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Repository struct {
	DB     *sqlx.DB
	Logger *slog.Logger
}

func NewRepository(db *sqlx.DB, logger *slog.Logger) *Repository {
	return &Repository{
		DB:     db,
		Logger: logger,
	}
}

func (repo *Repository) GetRoles() ([]*Role, error) {
	var roles []*Role

	rows, err := repo.DB.Queryx("SELECT * FROM roles LIMIT 10")
	if err != nil {
		repo.Logger.Error("role query failed", err)
		return nil, err
	}

	for rows.Next() {
		var role Role
		if err := rows.StructScan(&role); err != nil {
			repo.Logger.Error("role scan failed", err)
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}
