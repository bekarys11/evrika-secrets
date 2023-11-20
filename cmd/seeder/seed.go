package seed

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func PopulateRoles(db *sqlx.DB) error {
	_, err := db.Exec(`
	INSERT INTO roles (name, alias) 
	VALUES ('Админ', 'admin'), ('Пользователь', 'user')
	ON CONFLICT DO NOTHING;
`)
	if err != nil {
		return fmt.Errorf("error inserting roles into database: %v", err)
	}
	return nil
}

func PopulateUsers(db *sqlx.DB) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

	db.Exec(`
		INSERT INTO users (name, email, password, is_active, role_id) 
		VALUES ('bekarys', 'bekarys.t@evrika.com', $1, true, 1)
		ON CONFLICT DO NOTHING`, string(hashed))
	if err != nil {
		return fmt.Errorf("error creating users: %v", err.Error())
	}

	return nil
}
