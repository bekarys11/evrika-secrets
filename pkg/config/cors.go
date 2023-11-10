package config

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func handleCORS(router *mux.Router) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	return c.Handler(router)
}
