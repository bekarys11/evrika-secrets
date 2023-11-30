package config

import (
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/users"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/casbin/casbin/v2"
	"log"
	"log/slog"
	"net/http"
)

func Authorizer(e *casbin.Enforcer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			claims, err := users.GetTokenClaims(r)
			if err != nil {
				resp.ErrorJSON(w, fmt.Errorf("get profile error: %v", err), http.StatusInternalServerError)
				return
			}

			role, ok := claims["role"]
			if !ok {
				slog.Info("role not found, role 'guest' is given")
				role = "guest"
			}
			// casbin rule enforcing
			log.Printf("role: %s, path: %s, method: %s", role, r.URL.Path, r.Method)
			isAuthorized, err := e.Enforce(role, r.URL.Path, r.Method)
			if err != nil {
				slog.Error("unable to enforce roles:", err)
				resp.ErrorJSON(w, fmt.Errorf("enforce: %v", err), 500)
				return
			}

			if isAuthorized {
				log.Println("user is authorized")
				next.ServeHTTP(w, r)
			} else {
				log.Println("user is not authorized")
				resp.ErrorJSON(w, fmt.Errorf("user is not authorized"), http.StatusUnauthorized)
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}
