package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
)

func (app *Config) ConnectToDB() (err error) {
	dsn := fmt.Sprint("host=localhost port=32781 user=bekarys password=mynewpassword dbname=secrets sslmode=disable")

	app.DB, err = sqlx.Open("postgres", dsn)

	if err != nil {
		return fmt.Errorf("error connecting database: %v", err)
	}

	if err = app.DB.Ping(); err != nil {
		return fmt.Errorf("error ping database: %v", err)
	}

	slog.Info("postgres connection established")

	return nil
}
