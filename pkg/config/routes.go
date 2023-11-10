package config

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
)

func (app *Config) LoadRoutes() {
	app.Router = mux.NewRouter()
	app.Router.HandleFunc("/hello", ArticlesCategoryHandler).Methods("GET")

	slog.Info("app running on PORT:" + os.Getenv("APP_PORT"))
}

func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello there"))
}
