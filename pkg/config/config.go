package config

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

type Config struct {
	DB     *sqlx.DB
	Router *mux.Router
}

func StartApp() {
	PORT := os.Getenv("APP_PORT")
	app := Config{}

	if err := app.ConnectToDB(); err != nil {
		log.Fatal(err)
	}
	defer app.DB.Close()

	app.LoadRoutes()
	handler := handleCORS(app.Router)

	log.Fatal(http.ListenAndServe(PORT, handler))
}
