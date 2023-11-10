package config

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Config struct {
	DB     *sqlx.DB
	Router *mux.Router
}

func StartApp() {
	app := Config{}

	//if err := app.ConnectToDB(); err != nil {
	//	log.Fatal(err)
	//}
	app.LoadRoutes()
	handler := handleCORS(app.Router)

	log.Fatal(http.ListenAndServe(":8888", handler))
}
