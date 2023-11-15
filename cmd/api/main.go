package main

import (
	"github.com/bekarys11/evrika-secrets/pkg/config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config.StartApp()
}
