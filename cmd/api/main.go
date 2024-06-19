package main

import (
	"os"

	"github.com/codepnw/go-ecommerce/config"
	"github.com/codepnw/go-ecommerce/internal/server"
	"github.com/codepnw/go-ecommerce/pkg/database"
	_ "github.com/joho/godotenv/autoload"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := database.DBConnect(cfg)
	db.Health()
	defer db.Close()

	server.NewServer(db).Start()
}
