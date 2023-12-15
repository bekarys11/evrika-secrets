package config

import (
	seed "github.com/bekarys11/evrika-secrets/cmd/seeder"
	"github.com/go-ldap/ldap"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Server *http.Server
	DB     *sqlx.DB
	Router *mux.Router
	LDAP   *ldap.Conn
}

func StartApp() (*http.Server, *sqlx.DB) {
	PORT := os.Getenv("APP_PORT")
	app := Config{}

	if err := app.ConnectToDB(); err != nil {
		log.Fatal(err)
	}

	if err := seed.PopulateRoles(app.DB); err != nil {
		log.Fatal(err)
	}

	if err := seed.PopulateUsers(app.DB); err != nil {
		log.Fatal(err)
	}

	if err := app.ConnectToLDAP(); err != nil {
		log.Fatal(err)
	}

	app.LoadRoutes()
	handler := handleCORS(app.Router)

	app.Server = &http.Server{
		Addr:           PORT,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	defer app.DB.Close()
	defer app.LDAP.Close()

	return app.Server, app.DB
}
