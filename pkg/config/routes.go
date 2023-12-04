package config

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/bekarys11/evrika-secrets/docs"
	"github.com/bekarys11/evrika-secrets/internal/roles"
	"github.com/bekarys11/evrika-secrets/internal/secrets"
	"github.com/bekarys11/evrika-secrets/internal/users"
	"github.com/bekarys11/evrika-secrets/pkg/auth"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"log/slog"
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
	e, err := casbin.NewEnforcer("auth_model.conf", "policy.csv")
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	if err != nil {
		log.Fatalf("Failed to create new enforcer: %v", err)
	}

	userRepo := &users.Repo{DB: app.DB, LDAP: app.LDAP, Validation: validate}
	authRepo := &auth.Repo{DB: app.DB}
	secretRepo := &secrets.Repo{DB: app.DB, QBuilder: psql}
	roleRepo := &roles.Repo{DB: app.DB}

	app.Router = mux.NewRouter()
	app.Router.HandleFunc("/api/v1/login", authRepo.Login)

	api := app.Router.PathPrefix("/api/v1").Subrouter()
	api.Use(Authenticator())
	api.Use(Authorizer(e))

	api.HandleFunc("/users", userRepo.All).Methods("GET")

	api.HandleFunc("/users", userRepo.All).Methods("GET")
	api.HandleFunc("/users", userRepo.Create).Methods("POST")

	api.HandleFunc("/profile", userRepo.GetProfile).Methods("GET")

	api.HandleFunc("/secrets", secretRepo.All).Methods("GET")
	api.HandleFunc("/secrets/{secret_id}", secretRepo.One).Methods("GET")

	api.HandleFunc("/secrets", secretRepo.Create).Methods("POST")
	api.HandleFunc("/secrets/share", secretRepo.ShareSecret).Methods("POST")

	api.HandleFunc("/roles", roleRepo.All).Methods("GET")

	app.Router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s%s/swagger/doc.json", os.Getenv("APP_URL"), os.Getenv("APP_PORT"))))).Methods("GET")

	slog.Info("app running on PORT:" + os.Getenv("APP_PORT"))
}
