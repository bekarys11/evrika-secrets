package config

import (
	"fmt"
	_ "github.com/bekarys11/evrika-secrets/docs"
	"github.com/bekarys11/evrika-secrets/internal/roles"
	"github.com/bekarys11/evrika-secrets/internal/secrets"
	"github.com/bekarys11/evrika-secrets/internal/users"
	"github.com/bekarys11/evrika-secrets/pkg/auth"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"
)

// @title           Evrika Secrets API
// @version         1.0
// @description     Platform for managing application secrets and keys.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      10.10.1.59:44044
// @BasePath  /api/v1
// @securityDefinitions.apiKey  ApiKeyAuth
// @in header
// @name Authorization
func (app *Config) LoadRoutes() {
	validate := validator.New(validator.WithRequiredStructEnabled())

	userRepo := &users.Repo{DB: app.DB, LDAP: app.LDAP, Validation: validate}
	authRepo := &auth.Repo{DB: app.DB}
	secretRepo := &secrets.Repo{DB: app.DB}
	roleRepo := &roles.Repo{DB: app.DB}

	app.Router = mux.NewRouter()
	app.Post("/api/v1/login", app.HandleRequest(authRepo.Login))

	app.Get("/api/v1/users", app.HandleGuardedRequest(userRepo.All))
	app.Post("/api/v1/users", app.HandleGuardedRequest(userRepo.Create))

	app.Get("/api/v1/profile", app.HandleGuardedRequest(userRepo.GetProfile))

	app.Get("/api/v1/secrets/{user_id}", app.HandleGuardedRequest(secretRepo.All))
	app.Post("/api/v1/secrets", app.HandleGuardedRequest(secretRepo.Create))
	app.Post("/api/v1/secrets/share", app.HandleGuardedRequest(secretRepo.ShareSecret))

	app.Get("/api/v1/roles", app.HandleGuardedRequest(roleRepo.All))

	app.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL("http://10.10.1.59:44044/swagger/doc.json"))).Methods("GET")

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
			resp.ErrorJSON(w, fmt.Errorf("token error: %v", err), http.StatusUnauthorized)
			return
		}

		if isValid {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
