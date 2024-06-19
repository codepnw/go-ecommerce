package server

import (
	"log"

	"github.com/codepnw/go-ecommerce/pkg/database"
	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Start()
	GetServer() *server
}

type server struct {
	db  database.Service
	app *fiber.App
}

func NewServer(db database.Service) Server {
	return &server{
		db:  db,
		app: fiber.New(),
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {
	log.Print("server is starting on :8080")
	s.app.Listen(":8080")
}
