package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/codepnw/go-ecommerce/config"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Service interface {
	Health() map[string]string
	Close() error
	Get() *sql.DB
}

type service struct {
	db *sql.DB
}

func DBConnect(cfg config.Config) Service {
	db, err := sql.Open(cfg.Db().Driver(), cfg.Db().Url())
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(cfg.Db().MaxOpenConns())

	return &service{
		db: db,
	}
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "Database connected!"
	return stats
}

func (s *service) Close() error {
	log.Println("Disconnected from database")
	return s.db.Close()
}

func (s *service) Get() *sql.DB {
	return s.db
}
