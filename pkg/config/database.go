package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

func connectToDB() (db *sqlx.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	db, err = sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("error connecting database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error ping database: %v", err)
	}

	slog.Info("postgres connection established")

	return db, nil
}
