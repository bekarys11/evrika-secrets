package config

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *Config) LoadRoutes() {
	app.Router = mux.NewRouter()

	app.Router.HandleFunc("/hello", ArticlesCategoryHandler).Methods("GET")

}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello there"))
}
