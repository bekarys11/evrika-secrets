package config

import (
	"fmt"
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

func New() *Config {
	PORT := os.Getenv("APP_PORT")

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

	router := loadRoutes(db, ldapConn)
	handler := handleCORS(router)

	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", PORT),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app := &Config{
		Server: server,
		DB:     db,
		Router: router,
		LDAP:   ldapConn,
	}

	return app
}
