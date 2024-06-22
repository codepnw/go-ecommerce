package server

import (
	"log"

	"github.com/codepnw/go-ecommerce/config"
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
	cfg config.Config
}

func NewServer(db database.Service, cfg config.Config) Server {
	return &server{
		db:  db,
		app: fiber.New(),
		cfg: cfg,
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {
	v1 := s.app.Group("v1")
	module := InitModule(v1, s)

	module.MonitorModule()

	log.Print("server is starting on :8080")
	s.app.Listen(":8080")
}
