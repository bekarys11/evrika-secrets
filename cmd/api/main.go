package main

import (
	"context"
	"github.com/bekarys11/evrika-secrets/pkg/config"
	_ "github.com/joho/godotenv/autoload"
	"time"
)

func main() {
	cfg := config.New()

	go cfg.Server.ListenAndServe()

	wait := config.GracefulShutdown(context.Background(), 2*time.Second, map[string]config.Operation{
		"database": func(ctx context.Context) error {
			return cfg.DB.Close()
		},
		"http-server": func(ctx context.Context) error {
			return cfg.Server.Shutdown(context.Background())
		},
	})
	<-wait
}
