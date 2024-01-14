package config

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connectToDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	sqlDB, err := sql.Open("postgres", dsn)
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})

	sqlxDB := sqlx.NewDb(sqlDB, "postgres")
	if err != nil {
		return nil, fmt.Errorf("error connecting database: %v", err)
	}

	if err = sqlxDB.Ping(); err != nil {
		return nil, fmt.Errorf("error ping database: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/migrations",
		"postgres", driver)

	if err != nil {
		log.Fatal("Error creating migration instance: %v", err)
	}

	//if err := m.Down(); err != nil {
	//	log.Fatalf("Faile to run cleanup migrations: %v", err)
	//}

	if err := m.Up(); err != nil {
		log.Printf("Failed to migrate: %v", err)
	}

	slog.Info("postgres db connection established")

	return sqlxDB, nil
}
