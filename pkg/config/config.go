package config

import (
	"fmt"
	seed "github.com/bekarys11/evrika-secrets/cmd/seeder"
	"github.com/go-ldap/ldap"
	"github.com/gorilla/handlers"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Server *http.Server
	DB     *sqlx.DB
	LDAP   *ldap.Conn
}

func New() *Config {
	APP_FULL_URL := os.Getenv("APP_FULL_URL")

	db, err := connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	if err := seed.PopulateRoles(db); err != nil {
		log.Fatal(err)
	}

	if err := seed.PopulateUsers(db); err != nil {
		log.Fatal(err)
	}

	ldapConn, err := connectToLDAP()
	if err != nil {
		log.Fatal(err)
	}

	logger := newLogger()
	router := loadRoutes(db, ldapConn, logger)
	handler := handleCORS(router)

	handler = handlers.LoggingHandler(os.Stdout, router)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s", APP_FULL_URL),
		Handler:      handler,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	app := &Config{
		Server: server,
		DB:     db,
		LDAP:   ldapConn,
	}

	return app
}
