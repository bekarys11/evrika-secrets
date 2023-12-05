package config

import (
	seed "github.com/bekarys11/evrika-secrets/cmd/seeder"
	"github.com/go-ldap/ldap"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

type Config struct {
	DB     *sqlx.DB
	Router *mux.Router
	LDAP   *ldap.Conn
}

func StartApp() {
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

	//if err := app.ConnectToLDAP(); err != nil {
	//	log.Fatal(err)
	//}

	app.LoadRoutes()
	handler := handleCORS(app.Router)

	defer app.DB.Close()
	defer app.LDAP.Close()

	log.Fatal(http.ListenAndServe(PORT, handler))
}
