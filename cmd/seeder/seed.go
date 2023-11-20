package seed

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func PopulateRoles(db *sqlx.DB) error {
	if _, err := db.Exec("TRUNCATE roles RESTART IDENTITY CASCADE ;"); err == nil {
		return fmt.Errorf("error truncating roles table: %v", err)
	}

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
	hashed, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("USER_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error generating password for user: %v", err.Error())
	}

	if _, err := db.Exec("TRUNCATE users RESTART IDENTITY CASCADE ;"); err == nil {
		return fmt.Errorf("error truncating roles table: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO users (name, email, password, is_active, role_id) 
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT DO NOTHING`, "bekarys", "bekarys.t@evrika.com", string(hashed), true, 1)
	if err != nil {
		return fmt.Errorf("error creating user: %v", err.Error())
	}

	return nil
}
