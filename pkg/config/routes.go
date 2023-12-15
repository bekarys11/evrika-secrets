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
	"github.com/go-ldap/ldap"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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
func loadRoutes(db *sqlx.DB, ldapConn *ldap.Conn) (router *mux.Router) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	e, err := casbin.NewEnforcer("auth_model.conf", "policy.csv")
	if err != nil {
		log.Fatalf("Failed to create new enforcer: %v", err)
	}

	userRepository := users.NewRepository(db, ldapConn, validate)
	userService := users.NewUserService(userRepository)
	userServer := users.NewHttpServer(userService)

	roleRepository := roles.NewRepository(db)
	roleService := roles.NewRoleService(roleRepository)
	roleServer := roles.NewHttpServer(roleService)

	authRepo := &auth.Repo{DB: db}
	secretRepo := &secrets.Repo{DB: db, QBuilder: psql}
	//roleRepo := &roles.Repo{DB: db}

	router = mux.NewRouter()
	router.HandleFunc("/api/v1/login", authRepo.Login)

	api := router.PathPrefix("/api/v1").Subrouter()
	api.Use(Authenticator())
	api.Use(Authorizer(e))

	api.HandleFunc("/users", userServer.GetUsers).Methods("GET")
	api.HandleFunc("/users", userServer.CreateUser).Methods("POST")
	api.HandleFunc("/profile", userServer.GetProfile).Methods("GET")

	api.HandleFunc("/secrets", secretRepo.All).Methods("GET")
	api.HandleFunc("/secrets/{secret_id}", secretRepo.One).Methods("GET")
	api.HandleFunc("/secrets/{secret_id}", secretRepo.Update).Methods("PUT")
	api.HandleFunc("/secrets/{secret_id}", secretRepo.Delete).Methods("DELETE")

	api.HandleFunc("/secrets", secretRepo.Create).Methods("POST")
	api.HandleFunc("/secrets/share", secretRepo.ShareSecret).Methods("POST")

	api.HandleFunc("/roles", roleServer.GetRoles).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("%s%s/swagger/doc.json", os.Getenv("APP_URL"), os.Getenv("SWAGGER_PORT"))))).Methods("GET")
	slog.Info("app running on PORT:" + os.Getenv("APP_PORT"))
	return router
}
