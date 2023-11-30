package config

import (
	"fmt"
	"github.com/bekarys11/evrika-secrets/internal/users"
	resp "github.com/bekarys11/evrika-secrets/pkg/response"
	"github.com/casbin/casbin/v2"
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
				role = "guest"
			}
			// casbin rule enforcing
			res, err := e.Enforce(role, r.URL.Path, r.Method)

			if err != nil {
				resp.WriteJSON(w, 500, fmt.Errorf("enforce: %v", err))
				return
			}

			if res {
				next.ServeHTTP(w, r)
			} else {
				resp.WriteJSON(w, 500, "enforce: not allowed")
				return
			}
		}

		return http.HandlerFunc(fn)
	}
}
