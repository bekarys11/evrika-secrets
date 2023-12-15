package main

import (
	"context"
	"github.com/bekarys11/evrika-secrets/pkg/config"
	_ "github.com/joho/godotenv/autoload"
	"time"
)

func main() {
	server, db := config.StartApp()

	wait := config.GracefulShutdown(context.Background(), 2*time.Second, map[string]config.Operation{
		"database": func(ctx context.Context) error {
			return db.Close()
		},
		"http-server": func(ctx context.Context) error {
			return server.Shutdown(context.Background())
		},
	})

	<-wait

	server.ListenAndServe()
}
