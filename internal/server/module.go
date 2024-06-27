package server

import (
	"github.com/codepnw/go-ecommerce/internal/middleware"
	"github.com/codepnw/go-ecommerce/internal/monitor"
	"github.com/codepnw/go-ecommerce/internal/users/usersHandlers"
	"github.com/codepnw/go-ecommerce/internal/users/usersRepositories"
	"github.com/codepnw/go-ecommerce/internal/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
	m middleware.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, m middleware.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
		m: m,
	}
}

func InitMiddleware(s *server) middleware.IMiddlewareHandler {
	repo := middleware.MiddlewareRepository(s.db.Get())
	usecase := middleware.MiddlewareUsecase(repo)
	return middleware.MiddlewareHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitor.MonitorHandler(m.s.cfg)

	m.r.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repo := usersRepositories.UsersRepository(m.s.db.Get())
	usecase := usersUsecases.UsersUsecase(m.s.cfg, repo)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")
	
	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
}