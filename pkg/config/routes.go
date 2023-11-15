package config

import (
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/users"
	"github.com/bekarys11/evrika-secrets/pkg/auth"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
)

func (app *Config) LoadRoutes() {
	userRepo := &users.Repo{DB: app.DB, LDAP: app.LDAP}
	authRepo := &auth.Repo{DB: app.DB}

	app.Router = mux.NewRouter()
	app.Post("/api/v1/login", app.HandleRequest(authRepo.Login))
	app.Get("/api/v1/users", app.HandleRequest(userRepo.All))
	app.Post("/api/v1/users", app.HandleRequest(userRepo.Create))

	slog.Info("app running on PORT:" + os.Getenv("APP_PORT"))
}

type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

func (app *Config) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("GET")
}

func (app *Config) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("POST")
}

func (app *Config) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("PUT")
}

func (app *Config) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	app.Router.HandleFunc(path, f).Methods("DELETE")
}

func (app *Config) HandleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}

func (app *Config) HandleGuardedRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isValid, err := auth.IsValidToken(r.Header.Get("Authorization"))

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Error with token validation: %v", err.Error())))
		}

		if isValid {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
